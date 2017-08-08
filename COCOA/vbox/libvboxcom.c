//
//  libvboxcom.c
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/12/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#include <assert.h>
#include <unistd.h>
#include <Block.h>

#include "common.h"
#include "vbox.h"
#include "session.h"
#include "machine.h"
#include "bios_settings.h"
#include "medium.h"
#include "medium_format.h"
#include "progress.h"
#include "storage_controller.h"
#include "network.h"
#include "host.h"

#include "libvboxcom.h"

#pragma mark - TYPES

#define ERROR_MESSAGE_BUF_SIZE 256

typedef struct ivbox_session {
    char                  error_msg[ERROR_MESSAGE_BUF_SIZE];
    IVirtualBox*          vbox;                             // virtualbox
    IVirtualBoxClient*    client;                           // vbox client
    ISession*             vsession;                         // vbox session
    IMachine*             machine;                          // vbox machine
    char*                 machine_id;                       // machine id
    char*                 setting_file_path;                // setting file
} ivbox_session;

#pragma mark - MACROS

// these two are convert between internal <-> external types
#define toiVBoxSession(ptr) ((ivbox_session*)ptr)
#define toVBoxGlue(ptr)     ((VBoxGlue)ptr)

#pragma mark - DECLARATION

#pragma mark build machine base

static HRESULT
vbox_machine_build(IVirtualBox* virtualbox, IMachine* vbox_machine, int cpu_count, int memory_size, char* error_message);

static HRESULT
vbox_machine_add_1st_nat_network(IMachine* vbox_machine, ISession* vbox_session, char* error_message);

static HRESULT
vbox_machine_add_2nd_bridged_network(IMachine* vbox_machine, ISession* vbox_session, const char* host_interface, char* error_message);

static HRESULT
vbox_machine_add_shared_folder(IMachine* vbox_machine, ISession* vbox_session, const char* shared_name, const char *host_folder, char* error_message);

static HRESULT
vbox_machine_add_storage_controller(IMachine* vbox_machine, ISession* vbox_session, const char* storage_controller_name, char* error_message);

static HRESULT
vbox_machine_add_boot_image(IVirtualBox* virtualbox, IMachine* vbox_machine, ISession* vbox_session, const char* storage_controller_name, const char *boot_image_path, char *error_message);

static HRESULT
vbox_machine_add_hard_disk(IVirtualBox* virtualbox, IMachine* vbox_machine, ISession* vbox_session, const char* storage_controller_name, const char *hdd_medium_path, int disk_size, char *error_message);


#pragma mark utils
/**
 * Print detailed error information if available.
 * @param   pszErrorMsg     string containing the code location specific error message
 * @param   rc              COM/XPCOM result code
 */
static inline void
print_error_info(char *buffer, const char *message, HRESULT rc) {
    memset(buffer, 0, ERROR_MESSAGE_BUF_SIZE);
    sprintf(buffer, "%s (CODE : 0x%08X)", message, rc);
}

#pragma mark - DEFINITION

#pragma mark init & close
VBGlueResult
NewVBoxGlue(VBoxGlue* glue) {

    // make sure the pointer passed IS null.
    assert(glue != NULL && *glue == NULL);
    
    HRESULT result;
    ivbox_session* session = (ivbox_session*)calloc(1, sizeof(ivbox_session));
    // assign to return value
    *glue = toVBoxGlue(session);
    
    result = VBoxCGlueInit();
    if (FAILED(result)) {
        // it is more reasonable to print lib's error message
        strcpy(session->error_msg, g_szVBoxErrMsg);
        return VBGlue_Fail;
    }

    result = VboxClientInitialize(&(session->client));
    if ( FAILED(result) || session->client == NULL ) {
        print_error_info(session->error_msg, "[VBox] failed to init vbox client", result);
        return VBGlue_Fail;
    }

    result = VboxGetVirtualBox(session->client, &(session->vbox));
    if ( FAILED(result) || session->vbox == NULL ) {
        print_error_info(session->error_msg, "[VBox] failed to get VirtualBox reference", result);
        return VBGlue_Fail;
    }

    result = VboxGetSession(session->client, &(session->vsession));
    if ( FAILED(result) || session->vsession == NULL ) {
        print_error_info(session->error_msg, "[VBox] Failed to get Session reference", result);
        return VBGlue_Fail;
    }

    return VBGlue_Ok;
}

VBGlueResult
CloseVBoxGlue(VBoxGlue glue) {

    // make sure the pointer passed is not null.
    assert(glue != NULL);
    
    ivbox_session* session = toiVBoxSession(glue);
    HRESULT result;
    
    if ( session->vsession != NULL ) {
        result = VboxISessionRelease(session->vsession);
        if (FAILED(result)) {
            print_error_info(session->error_msg, "[VBox] Failed to release ISession reference", result);
            return VBGlue_Fail;
        }
    }
    if ( session->vbox != NULL ) {
        result = VboxIVirtualBoxRelease(session->vbox);
        if (FAILED(result)) {
            print_error_info(session->error_msg, "[VBox] failed to release vbox reference", result);
            return VBGlue_Fail;
        }
    }
    if ( session->client != NULL ) {
        result = VboxClientRelease(session->client);
        if (FAILED(result)) {
            print_error_info(session->error_msg, "[VBox] failed to release vbox client", result);
            return VBGlue_Fail;
        }
    }
    
    VboxClientUninitialize();
    VBoxCGlueTerm();
    free(session);
    
    return VBGlue_Ok;
}


#pragma mark host network interface
VBGlueResult
VBoxSearchHostNetworkInterfaceByName(VBoxGlue glue, const char* queryName, char** nameFound) {

    assert(glue != NULL);
    assert(queryName != NULL && strlen(queryName) != 0);
    
    ivbox_session* session = toiVBoxSession(glue);
    IHost* host = NULL;
    HRESULT result;
    
    result = VboxGetHostFromVirtualbox(session->vbox, &host);
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] failed to get Host reference", result);
        return VBGlue_Fail;
    }
    
    result = VboxHostSearchNetworkInterfacesFromList(host, queryName, nameFound);
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] failed to find host networkInterface", result);
        return VBGlue_Fail;
    }
    
    return result;
}

#pragma mark machine status
/*
 typedef enum MachineState {
   MachineState_Null = 0,
   MachineState_PoweredOff = 1,
   MachineState_Saved = 2,
   MachineState_Teleported = 3,
   MachineState_Aborted = 4,
   MachineState_Running = 5,
   MachineState_Paused = 6,
   MachineState_Stuck = 7,
   MachineState_Teleporting = 8,
   MachineState_LiveSnapshotting = 9,
   MachineState_Starting = 10,
   MachineState_Stopping = 11,
   MachineState_Saving = 12,
   MachineState_Restoring = 13,
   MachineState_TeleportingPausedVM = 14,
   MachineState_TeleportingIn = 15,
   MachineState_FaultTolerantSyncing = 16,
   MachineState_DeletingSnapshotOnline = 17,
   MachineState_DeletingSnapshotPaused = 18,
   MachineState_OnlineSnapshotting = 19,
   MachineState_RestoringSnapshot = 20,
   MachineState_DeletingSnapshot = 21,
   MachineState_SettingUp = 22,
   MachineState_Snapshotting = 23,
   MachineState_FirstOnline = 5,
   MachineState_LastOnline = 19,
   MachineState_FirstTransient = 8,
   MachineState_LastTransient = 23
 } MachineState;
 #define MachineState_T PRUint32
 */

VBGlueMachineState
VBoxMachineGetCurrentState(VBoxGlue glue) {
    // make sure the pointer passed is not null.
    assert(glue != NULL);

    ivbox_session* session = toiVBoxSession(glue);
    MachineState mState;
    HRESULT result;

    result = VboxMachineGetState(session->machine, &mState);
    if (FAILED(result)) {
        return VBGlueMachine_Illegal;
    }

    switch (mState) {
        case MachineState_PoweredOff:
        case MachineState_Aborted:
        case MachineState_Running:
        case MachineState_Paused:
        case MachineState_Stuck:
        case MachineState_Starting:
        case MachineState_Stopping: {
            return (VBGlueMachineState)mState;
        }
        default:
            return VBGlueMachine_Illegal;
    }

    return VBGlueMachine_Illegal;
}

VBGlueResult
VBoxMachineIsSettingChanged(VBoxGlue glue, bool* isMachineChanged) {
    
    // make sure the pointer passed is not null.
    assert(glue != NULL);

    ivbox_session* session = toiVBoxSession(glue);
    PRBool changed = PR_FALSE;
    *isMachineChanged = (bool)changed;
    IMachine *mutable_machine;
    
    //firstly lock the machine
    HRESULT result = VboxLockMachine(session->machine, session->vsession, LockType_Write);
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] Failed to lock machine for checking machine setting", result);
        return VBGlue_Fail;
    }
    // get mutable machine
    result = VboxGetSessionMachine(session->vsession, &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(session->error_msg, "[VBox] Failed to get a mutable copy of a machine", result);
        return VBGlue_Fail;
    }
    // check if settings modified
    result = VboxGetMachineSettingsModified(mutable_machine, &changed);
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] Fail to get setting modification", result);
        return VBGlue_Fail;
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            print_error_info(session->error_msg, "[VBox] Failed to release locked machine", result);
            return VBGlue_Fail;
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(session->vsession);
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] Failed to unlock machine for attaching adapter", result);
        return VBGlue_Fail;
    }

    *isMachineChanged = (bool)(changed == PR_TRUE);
    return VBGlue_Ok;
}


#pragma mark find, create, build, & release, destroy machine
VBGlueResult
VBoxMachineFindByNameOrID(VBoxGlue glue, const char* machineName) {
    
    // make sure the pointer passed is not null.
    assert(glue != NULL);

    ivbox_session* session = toiVBoxSession(glue);

    if ( session->machine  == NULL ) {
        HRESULT result;
        result = VboxFindMachine(session->vbox, machineName, &(session->machine));
        if (FAILED(result)) {
            print_error_info(session->error_msg, "[VBox] Failed to find machine", result);
            return VBGlue_Fail;
        }
        result = VboxGetMachineSettingsFilePath(session->machine, &(session->setting_file_path));
        if (FAILED(result)) {
            print_error_info(session->error_msg, "[VBox] Fail to get setting file path", result);
            return VBGlue_Fail;
        }
        result = VboxMachineGetID(session->machine, &(session->machine_id));
        if (FAILED(result)) {
            print_error_info(session->error_msg, "[VBox] Failed to get Machine ID", result);
            return VBGlue_Fail;
        }

    } else {
        print_error_info(session->error_msg, "[VBox] machine instance already exists", NS_OK);
        return VBGlue_Fail;
    }
    return VBGlue_Ok;
}

VBGlueResult
VBoxMachineCreateByName(VBoxGlue glue, const char* baseFolder, const char* machineName) {
    
    // make sure the pointer passed is not null.
    assert(glue != NULL);
    assert(machineName != NULL || strlen(machineName) != 0);
    
    ivbox_session* session = toiVBoxSession(glue);
    HRESULT result;

    if ( session->machine  != NULL ) {
        print_error_info(session->error_msg, "[VBox] machine instance already exists", NS_OK);
        return VBGlue_Fail;
    }
    
    // create machine file name
    result = VboxComposeMachineFilename(session->vbox, machineName, "", (char *)baseFolder, &(session->setting_file_path));
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] Failed composing machine name", result);
        return VBGlue_Fail;
    }
    // create machine based on the
    result = VboxCreateMachine(session->vbox, session->setting_file_path, machineName, "Linux26_64", "", &(session->machine));
    if (FAILED(result) || session->machine == NULL) {
        print_error_info(session->error_msg, "[VBox] Failed to create machine", result);
        return VBGlue_Fail;
    }
    result = VboxMachineGetID(session->machine, &(session->machine_id));
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] Failed to get Machine ID", result);
        return VBGlue_Fail;
    }

    return VBGlue_Ok;
}

#pragma mark build & destroy machine
HRESULT
vbox_machine_build(IVirtualBox* virtualbox, IMachine* vbox_machine, int cpu_count, int memory_size, char* error_message) {

    assert(virtualbox  != NULL);
    assert(vbox_machine != NULL);

    HRESULT result;

    // Setup BIOS
    {
        // get BIOS settings
        IBIOSSettings *bios;
        result = VboxGetMachineBIOSSettings(vbox_machine, &bios);
        if ( FAILED(result) || bios == NULL ) {
            print_error_info(error_message, "[VBox] Failed to acquire bios settings", result);
            return result;
        }
        // enable I/O APIC
        result = IBIOSSettings_SetIOAPICEnabled(bios, PR_TRUE);
        if ( FAILED(result) ) {
            print_error_info(error_message, "[VBox] Failed to enable IO/APIC", result);
            return result;
        }
        // set ACPI enabled
        result = IBIOSSettings_SetACPIEnabled(bios, PR_TRUE);
        if ( FAILED(result) ) {
            print_error_info(error_message, "[VBox] Failed to enable ACPI", result);
            return result;
        }
        result = IBIOSSettings_SetLogoFadeIn(bios, PR_FALSE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to disable bios fade-in", result);
            return result;
        }
        result = IBIOSSettings_SetLogoFadeOut(bios, PR_FALSE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to disable bios fade-out", result);
            return result;
        }
        result = IBIOSSettings_SetLogoDisplayTime(bios, 0);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to disable bios display time", result);
            return result;
        }
        result = IBIOSSettings_SetBootMenuMode(bios, BIOSBootMenuMode_Disabled);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to disable bios menu", result);
            return result;
        }
        // release bios settings
        result = VboxIBiosSettingsRelease(bios);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to release BIOS", result);
            return result;
        }
    }
    
    // Motherboard Settings
    {
        // set memory
        result = VboxSetMachineMemorySize(vbox_machine, memory_size);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to set memory size", result);
            return result;
        }
        
        // set up Boot Order
        result = IMachine_SetBootOrder(vbox_machine, 1, DeviceType_DVD);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to fix boot order", result);
            return result;
        }
        result = IMachine_SetBootOrder(vbox_machine, 2, DeviceType_HardDisk);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to fix boot order", result);
            return result;
        }
        result = IMachine_SetBootOrder(vbox_machine, 3, DeviceType_Null);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to fix boot order", result);
            return result;
        }
        result = IMachine_SetBootOrder(vbox_machine, 4, DeviceType_Null);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to fix boot order", result);
            return result;
        }
        
        result = IMachine_SetFirmwareType(vbox_machine, FirmwareType_BIOS);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to set firmware time", result);
            return result;
        }
        
        // enable high precision event timer
        result = IMachine_SetHPETEnabled(vbox_machine, PR_TRUE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to enable HPET", result);
            return result;
        }
        // set Chipset type
        result = IMachine_SetChipsetType(vbox_machine, ChipsetType_ICH9);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to setting chipset type", result);
            return result;
        }
        // set RTC timer
        result = IMachine_SetRTCUseUTC(vbox_machine, PR_TRUE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to setting Hardware UTC timer", result);
            return result;
        }
    }
    
    // Processor Setting
    {
        // set CPU Count
        result = IMachine_SetCPUCount(vbox_machine, cpu_count);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBGlue_Error] Failed to setting CPU count", result);
            return result;
        }
        // set Execution Cap
        result = IMachine_SetCPUExecutionCap(vbox_machine, 100);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to setting CPU execution cap", result);
            return result;
        }
        // PAE enabled
        result = IMachine_SetCPUProperty(vbox_machine, CPUPropertyType_PAE, PR_TRUE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to enable PAE/NX", result);
            return result;
        }
        // long mode enabled for 64bit os
        result = IMachine_SetCPUProperty(vbox_machine, CPUPropertyType_LongMode, PR_TRUE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to enable longmode", result);
            return result;
        }
        // use hardware virtualization (VT-x/AMD-V) if available
        result = IMachine_SetHWVirtExProperty(vbox_machine, HWVirtExPropertyType_Enabled, PR_TRUE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to enable (VT-x/AMD-V)", result);
            return result;
        }
        // use VT-x VPID if available
        result = IMachine_SetHWVirtExProperty(vbox_machine, HWVirtExPropertyType_VPID, PR_TRUE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to enable VT-x VPID", result);
            return result;
        }
        // use Nested Paging if available
        result = IMachine_SetHWVirtExProperty(vbox_machine, HWVirtExPropertyType_NestedPaging, PR_TRUE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to enable Nested Paging", result);
            return result;
        }
        // use large page allocation if available. Requires nested paging and a 64-bit host.
        result = IMachine_SetHWVirtExProperty(vbox_machine, HWVirtExPropertyType_LargePages, PR_TRUE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to enable large page allocation", result);
            return result;
        }
    }
    
    // Acceleration
    {
        // Paravirtualization setting
        result = IMachine_SetParavirtProvider(vbox_machine, ParavirtProvider_Default);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to setting Pravirtualization", result);
            return result;
        }
        // Nested Paging
        result = IMachine_SetHWVirtExProperty(vbox_machine, HWVirtExPropertyType_Enabled, (PRBool)1);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to setting HWVirtExPropertyType", result);
            return result;
        }
    }
    
    // Display
    {
        // set VRAM
        result = VboxSetMachineVRAMSize(vbox_machine, 16);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to VRAM size", result);
            return result;
        }
    }
    
    // audio
    {
        IAudioAdapter *audio = NULL;
        result = IMachine_get_AudioAdapter(vbox_machine, &audio);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] failed get audio control", result);
            return result;
        }
        result = IAudioAdapter_SetEnabled(audio, PR_FALSE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] failed get audio control", result);
            return result;
        }
        result = IAudioAdapter_Release(audio);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] failed release audio control", result);
            return result;
        }
        audio = NULL;
    }

    // misc settings
    {
        // disable teleportation (this cause false "setting changed" flag)
        result = IMachine_SetTeleporterEnabled(vbox_machine, PR_FALSE);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to disable teleportation", result);
            return result;
        }
    }
    
    // SAVE SETTINGS & REGISTER MACHINE BEFORE ATTACH A MEDIUM
    {
        // save settings
        result = VboxMachineSaveSettings(vbox_machine);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to save machine before attaching a medium", result);
            return result;
        }
        // Register machine
        result = VboxRegisterMachine(virtualbox, vbox_machine);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to register machine", result);
            return result;
        }
    }
    return NS_OK;
}


HRESULT
vbox_machine_add_1st_nat_network(IMachine* vbox_machine, ISession* vbox_session, char* error_message) {
    static const char* NAT_RULE_PCSSH = "PC_SSH";
    static const char* NAT_HOST_IP    = "127.0.0.1";
    static const char* NAT_GUEST_IP   = "";
    
    INetworkAdapter *adapter = NULL;
    
    assert(vbox_machine != NULL);
    assert(vbox_session != NULL);
    
    //firstly lock the machine
    HRESULT result = VboxLockMachine(vbox_machine, vbox_session, LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to lock machine for networking", result);
        return result;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(vbox_session, &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[VBox] Failed to get a mutable copy of a machine for networking", result);
        return result;
    }
    // get network adapter
    result = VboxMachineGetNetworkAdapter(mutable_machine, 0, &adapter);
    if (FAILED(result) || adapter == NULL) {
        print_error_info(error_message, "[VBox] Failed to acquire adapter from slot 0", result);
        return result;
    }
    // enable network adapter
    result = VboxNetworkAdapterSetEnabled(adapter, TRUE);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to enable network adapter", result);
        return result;
    }
    // set nat network type
    result = VboxNetworkAdapterSetAttachmentType(adapter, NetworkAttachmentType_NAT);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to set network attachement type", result);
        return result;
    }
    // set adapter type (Virtio-Net)
    result = VboxNetworkAdapterSetAdapterType(adapter, NetworkAdapterType_Virtio);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to set network adapter type", result);
        return result;
    }
    // set cable connected
    result = VboxNetworkAdapterSetCableConnected(adapter, TRUE);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to set cable connected", result);
        return result;
    }
    
    // add port forwarding rule
    result = VboxNetworkAddPortForwardingRule(adapter, NAT_RULE_PCSSH, NATProtocol_TCP, NAT_HOST_IP, (unsigned short)3022, NAT_GUEST_IP, (unsigned short)3022);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to port forward for ssh", result);
        return result;
    }
    
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to save machine after attaching hard disk medium", result);
        return result;
    }
    // release the first adapter
    result = VboxNetworkAdapterRelease(adapter);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to release adapter", result);
        return result;
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to release locked machine for attaching adapter", result);
            return result;
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(vbox_session);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to unlock machine for attaching adapter", result);
        return result;
    }
    
    return NS_OK;
}

HRESULT
vbox_machine_add_2nd_bridged_network(IMachine* vbox_machine, ISession* vbox_session, const char* host_interface, char* error_message) {

    INetworkAdapter *adapter = NULL;
    
    assert(vbox_machine != NULL);
    assert(vbox_session != NULL);
    assert(host_interface != NULL && strlen(host_interface) != 0);

    //firstly lock the machine
    HRESULT result = VboxLockMachine(vbox_machine, vbox_session, LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to lock machine for networking", result);
        return result;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(vbox_session, &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[VBox] Failed to get a mutable copy of a machine for networking", result);
        return result;
    }
    // get network adapter
    result = VboxMachineGetNetworkAdapter(mutable_machine, 1, &adapter);
    if (FAILED(result) || adapter == NULL) {
        print_error_info(error_message, "[VBox] Failed to acquire adapter from slot 0", result);
        return result;
    }
    // enable network adapter
    result = VboxNetworkAdapterSetEnabled(adapter, TRUE);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to enable network adapter", result);
        return result;
    }
    // set bridged network type
    result = VboxNetworkAdapterSetAttachmentType(adapter, NetworkAttachmentType_Bridged);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to set network attachement type", result);
        return result;
    }
    // set host network adapter this bridge should connect to
    result = VboxNetworkAdapterSetBridgedHostInterface(adapter, host_interface);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to connect to host network interface", result);
        return result;
    }    
    // set adapter type (Virtio-Net)
    result = VboxNetworkAdapterSetAdapterType(adapter, NetworkAdapterType_Virtio);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to set network adapter type", result);
        return result;
    }
    // promiscuous mode policy
    result = VboxNetworkAdapterSetPromiscModePolicy(adapter, NetworkAdapterPromiscModePolicy_Deny);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to set promiscuous mode", result);
        return result;
    }
    // set cable connected
    result = VboxNetworkAdapterSetCableConnected(adapter, TRUE);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to set cable connected", result);
        return result;
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to save machine after attaching hard disk medium", result);
        return result;
    }
    // release the first adapter
    result = VboxNetworkAdapterRelease(adapter);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to release adapter", result);
        return result;
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to release locked machine for attaching adapter", result);
            return result;
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(vbox_session);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to unlock machine for attaching adapter", result);
        return result;
    }
    
    return NS_OK;
}

HRESULT
vbox_machine_add_shared_folder(IMachine* vbox_machine, ISession* vbox_session, const char* shared_name, const char *host_folder, char* error_message) {
    
    HRESULT result;
    
    assert(vbox_machine != NULL);
    assert(vbox_session != NULL);
    assert(shared_name != NULL && strlen(shared_name) != 0);
    
    //firstly lock the machine
    result = VboxLockMachine(vbox_machine, vbox_session, LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to lock machine for shared folder", result);
        return result;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(vbox_session, &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[VBox] Failed to get a mutable copy of a machine for shared folder", result);
        return result;
    }
    // create shared folder
    result = VboxMachineCreateSharedFolder(mutable_machine, shared_name, host_folder, TRUE, TRUE);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to add shared folder", result);
        return result;
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to save machine for adding shared folder", result);
        return result;
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to release machine for adding shared folder", result);
            return result;
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(vbox_session);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to unlock machine for adding shared folder", result);
        return result;
    }
    return NS_OK;
}

HRESULT
vbox_machine_add_storage_controller(IMachine* vbox_machine, ISession* vbox_session, const char* storage_controller_name, char* error_message) {

    HRESULT result;
    IStorageController *storage_controller;
    
    assert(vbox_machine != NULL);
    assert(vbox_session != NULL);
    assert(storage_controller_name != NULL && strlen(storage_controller_name) != 0);
    
    //firstly lock the machine
    result = VboxLockMachine(vbox_machine, vbox_session, LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to lock machine for adding storage controller", result);
        return result;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(vbox_session, &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[VBox] Failed to get a mutable copy of a machine", result);
        return result;
    }
    // add storage controller
    result = VboxMachineAddStorageController(mutable_machine, storage_controller_name, StorageBus_SATA, &storage_controller);
    if (FAILED(result) || storage_controller_name == NULL) {
        print_error_info(error_message, "[VBox] Failed to add storage controller", result);
        return result;
    }
    // storage controller type
    result = VboxSetStorageControllerType(storage_controller, StorageControllerType_IntelAhci);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to set storage controller type", result);
        return result;
    }
    // storage controller set # of ports
    result = IStorageController_SetPortCount(storage_controller, 2);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to increase port count", result);
        return result;
    }
    // Enable host IO cache for imaging
    PRBool use_host_iocache = (PRBool)1;
    result = IStorageController_SetUseHostIOCache(storage_controller, use_host_iocache);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to enable host IO cache", result);
        return result;
    }
    // release storage controller
    result = VboxIStorageControllerRelease(storage_controller);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to release storage controller for attaching storage controller", result);
        return result;
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to save machine after attaching storage controller", result);
        return result;
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to release mutable machine for attaching storage controller", result);
            return result;
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(vbox_session);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to unlock machine for attaching storage controller", result);
        return result;
    }

    return NS_OK;
}

HRESULT
vbox_machine_add_boot_image(IVirtualBox* virtualbox, IMachine* vbox_machine, ISession* vbox_session, const char* storage_controller_name, const char *boot_image_path, char *error_message) {
    
    IMedium *boot_image;
    HRESULT result;

    assert(virtualbox   != NULL);
    assert(vbox_machine != NULL);
    assert(vbox_session != NULL);
    assert(storage_controller_name != NULL && strlen(storage_controller_name) != 0);
    assert(boot_image_path != NULL && strlen(boot_image_path) != 0);
    
    // open medium
    PRBool use_old_uuid = (PRBool)0;
    result = VboxOpenMedium(virtualbox, boot_image_path, DeviceType_DVD, AccessMode_ReadOnly, use_old_uuid, &boot_image);
    if (FAILED(result) || boot_image == NULL) {
        print_error_info(error_message, "[VBox] Failed to open boot image for attaching boot image", result);
        return result;
    }
    //firstly lock the machine
    result = VboxLockMachine(vbox_machine, vbox_session, LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to lock machine for attaching boot image", result);
        return result;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(vbox_session, &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[VBox] Failed to get a mutable copy of a machine", result);
        return result;
    }
    // attach a medium
    result = VboxMachineAttachDevice(mutable_machine, storage_controller_name, 0, 0, DeviceType_DVD, boot_image);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to attach boot image", result);
        return result;
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to save machine after attaching boot image", result);
        return result;
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to release machine after attaching boot image", result);
            return result;
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(vbox_session);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to unlock machine for attaching boot image", result);
        return result;
    }
    // release medium
    result = VboxIMediumRelease(boot_image);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to release boot image", result);
        return result;
    }
    return NS_OK;
}

HRESULT
vbox_machine_add_hard_disk(IVirtualBox* virtualbox, IMachine* vbox_machine, ISession* vbox_session, const char* storage_controller_name, const char *hdd_medium_path, int disk_size, char *error_message) {

    // Create and Open hard drive
    HRESULT result;
    IMedium *hdd_medium;
    
    assert(virtualbox   != NULL);
    assert(vbox_machine != NULL);
    assert(vbox_session != NULL);
    assert(storage_controller_name != NULL && strlen(storage_controller_name) != 0);
    assert(hdd_medium_path != NULL && strlen(hdd_medium_path) != 0);

    {
// (2017/08/03) we'll use VBoxManage and raw input to create VMDK
#if CREATE_HDD_WITH_XPCOM
        // set medium image
        result = VboxCreateHardDisk(virtualbox, "VMDK", hdd_medium_path, DeviceType_HardDisk, AccessMode_ReadWrite, &hdd_medium);
        if (FAILED(result) || hdd_medium == NULL) {
            print_error_info(error_message, "[VBGlue_Error] Failed to create harddrive", result);
            return result;
        }
        // create medium
        //REF : https://www.virtualbox.org/sdkref/_virtual_box_8idl.html#adedcbf1a6e5e35fe7a0ca0c4b3447154
        PRUint32 cVariant[2] = {MediumVariant_Standard, MediumVariant_NoCreateDir};
        PRUint32 variantCount = sizeof(cVariant) / sizeof(cVariant[0]);
        IProgress *progress;
        result = VboxMediumCreateBaseStorage(hdd_medium, (1 << 24), variantCount, cVariant, &progress);
        if (FAILED(result)){
            print_error_info(error_message, "[VBox] Failed to create base storage", result);
            return result;
        }
        
        // it is recommended to wait short amount of time
        result = VboxProgressWaitForCompletion(progress, 500);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to complete creating base storage", result);
            return result;
        }
        else {
            PRUint32 progress_percent = 0;
            do {
                VboxProgressGetPercent(progress, &progress_percent);
                usleep(500000);
            } while (progress_percent < 100);
        }
        
        // get completion code
        PRInt32 code;
        result = VboxProgressGetResultCode(progress, &code);
        if (FAILED(result)|| code != 0) {
            print_error_info(error_message, "[VBox] Failed to actuqire storage creation result code", result);
            return result;
        }
        // release progress
        VboxProgressRelease(progress);
#endif
        // open medium
        PRBool use_old_uuid = (PRBool)0;
        result = VboxOpenMedium(virtualbox, hdd_medium_path, DeviceType_HardDisk, AccessMode_ReadWrite, use_old_uuid, &hdd_medium);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to open hard drive", result);
            return result;
        }
    }
    
    // Attach a medium to storage controller
    {
        //firstly lock the machine
        result = VboxLockMachine(vbox_machine, vbox_session, LockType_Write);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to lock machine for attaching hdd to storage controller", result);
            return result;
        }
        // get mutable machine
        IMachine *mutable_machine;
        result = VboxGetSessionMachine(vbox_session, &mutable_machine);
        if (FAILED(result) || mutable_machine == NULL) {
            print_error_info(error_message, "[VBox] Failed to get a mutable copy of a machine for attaching hard disk medium", result);
            return result;
        }
        // attach a medium
        result = VboxMachineAttachDevice(mutable_machine, storage_controller_name, 1, 0, DeviceType_HardDisk, hdd_medium);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to attach hard disk medium", result);
            return result;
        }
        // save setting
        result = VboxMachineSaveSettings(mutable_machine);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to save machine after attaching hard disk medium", result);
            return result;
        }
        // then we can safely release the mutable machine
        if (mutable_machine) {
            result = VboxIMachineRelease(mutable_machine);
            if (FAILED(result)) {
                print_error_info(error_message, "[VBox] Failed to release machine after attaching hard disk medium", result);
                return result;
            }
        }
        // then unlock machine
        result = VboxUnlockMachine(vbox_session);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBox] Failed to unlock machine after attaching hard disk medium", result);
            return result;
        }
    }
    
    // Close & release hard drive
    result = VboxIMediumRelease(hdd_medium);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBox] Failed to release a hard drive", result);
        return result;
    }
    
    return NS_OK;
}


#pragma mark -
VBoxBuildOption*
VBoxMakeBuildOption(int cpu, int mem, const char* host, const char* boot, const char* hdd, const char* spath, const char* sname) {
    VBoxBuildOption* option = (VBoxBuildOption*)calloc(1, sizeof(VBoxBuildOption));
    option->CpuCount      = cpu;
    option->MemSize       = mem;
    option->HostInterface = host;
    option->BootImagePath = boot;
    option->HddImagePath  = hdd;
    option->SharedDirPath = spath;
    option->SharedDirName = sname;
    return option;
}

VBGlueResult
VBoxMachineBuildWithOption(VBoxGlue glue, VBoxBuildOption* option) {
    
    static const char* STORAGE_CONTROLLER_NAME = "SATA";

    ivbox_session* session = toiVBoxSession(glue);
    HRESULT result = NS_OK;

    // sanity check
    assert(session           != NULL);
    assert(session->vbox     != NULL);
    assert(session->client   != NULL);
    assert(session->machine  != NULL);
    assert(session->vsession != NULL);

    assert(option != NULL);
    assert(0 < option->CpuCount);
    assert(2048 <= option->MemSize);
    assert(strlen(option->HostInterface) != 0);
    assert(strlen(option->SharedDirPath) != 0);
    assert(strlen(option->SharedDirName) != 0);
    assert(strlen(option->BootImagePath) != 0);
    assert(strlen(option->HddImagePath)  != 0);

    // build basic machine with bios & motherboard settings
    result = vbox_machine_build(session->vbox, session->machine, option->CpuCount, option->MemSize, session->error_msg);
    if (FAILED(result)) {
        return VBGlue_Fail;
    }

    // add nat network
    result = vbox_machine_add_1st_nat_network(session->machine, session->vsession, session->error_msg);
    if (FAILED(result)) {
        return VBGlue_Fail;
    }
    
    // add bridged network
    result = vbox_machine_add_2nd_bridged_network(session->machine, session->vsession, option->HostInterface, session->error_msg);
    if (FAILED(result)) {
        return VBGlue_Fail;
    }

    // add shared folder
    result = vbox_machine_add_shared_folder(session->machine, session->vsession, option->SharedDirName, option->SharedDirPath, session->error_msg);
    if (FAILED(result)) {
        return VBGlue_Fail;
    }

    // add storage controller
    result = vbox_machine_add_storage_controller(session->machine, session->vsession, STORAGE_CONTROLLER_NAME, session->error_msg);
    if (FAILED(result)) {
        return VBGlue_Fail;
    }

    // add boot image
    result = vbox_machine_add_boot_image(session->vbox, session->machine, session->vsession, STORAGE_CONTROLLER_NAME, option->BootImagePath, session->error_msg);
    if (FAILED(result)) {
        return VBGlue_Fail;
    }
    
    // add hard drive
    result = vbox_machine_add_hard_disk(session->vbox, session->machine, session->vsession, STORAGE_CONTROLLER_NAME, option->HddImagePath, 200000, session->error_msg);
    if (FAILED(result)) {
        return VBGlue_Fail;
    }

    return VBGlue_Ok;
}

VBGlueResult
VBoxMachineRelease(VBoxGlue glue) {

    // make sure the pointer passed is not null.
    assert(glue != NULL);

    ivbox_session* session = toiVBoxSession(glue);

    // release machine
    if ( session->setting_file_path != NULL ) {
        VboxUtf8Free(session->setting_file_path);
        session->setting_file_path = NULL;
    }
    if ( session->machine_id != NULL ) {
        VboxUtf8Free(session->machine_id);
        session->machine_id = NULL;
    }
    if (session->machine != NULL) {
        HRESULT result = VboxIMachineRelease(session->machine);
        if (FAILED(result)) {
            print_error_info(session->error_msg, "[VBox] Failed to close machine referenece", result);
            return VBGlue_Fail;
        } else {
            session->machine = NULL;
        }
    }

    return VBGlue_Ok;
}

VBGlueResult
VBoxMachineDestory(VBoxGlue glue) {

    // make sure the pointer passed is not null.
    assert(glue != NULL);
    
    ivbox_session* session = toiVBoxSession(glue);
    HRESULT result;
    ULONG media_count;
    IProgress *progress;
    IMedium** media;
    
    // unregister
    result = VboxMachineUnregister(session->machine, CleanupMode_DetachAllReturnHardDisksOnly, &media, &media_count);
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] Failed to unregister media", result);
        return VBGlue_Fail;
    }
    // delete medium
    result = VboxMachineDeleteConfig(session->machine, media_count, media, &progress);
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] Failed to delete medium", result);
        return VBGlue_Fail;
    }
    // delete progress
    VboxProgressWaitForCompletion(progress, 500);
    PRUint32 progress_percent = 0;
    do {
        VboxProgressGetPercent(progress, &progress_percent);
        usleep(500000);
    } while (progress_percent < 100);
    VboxProgressRelease(progress);
    
    // free media array
    VboxArrayOutFree(media);

    // release machine
    return VBoxMachineRelease(glue);
}


#pragma mark start & stop machine
VBGlueResult
VBoxMachineHeadlessStart(VBoxGlue glue) {

    static char* OPT_HEADLESS    = "headless";
    static char* OPT_ENVIRONMENT = "";

    ivbox_session* session = toiVBoxSession(glue);
    HRESULT result;
    IProgress *progress;
    PRInt32 resultCode;

    // make sure the pointer passed is not null.
    assert(glue != NULL);

    result = VboxMachineLaunchVMProcess(session->machine, session->vsession, OPT_HEADLESS, OPT_ENVIRONMENT, &progress);
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] Failed to launch virtual machine", result);
        return VBGlue_Fail;
    }

    // wait for startup
    VboxProgressWaitForCompletion(progress, 500);
    PRUint32 progress_percent = 0;
    do {
        VboxProgressGetPercent(progress, &progress_percent);
        usleep(500000);
    } while (progress_percent < 100);

    // get the progress result
    result = VboxProgressGetResultCode(progress, &resultCode);
    if (FAILED(result)) {
        char *error_message = NULL;
        HRESULT ret = VboxProgressGetResultInfo(progress, &error_message);
        if (SUCCEEDED(ret)) {
            strcpy(session->error_msg, error_message);
            free(error_message);
        }
    }

    VboxProgressRelease(progress);
    return (FAILED(result) ? VBGlue_Fail : VBGlue_Ok);
}

VBGlueResult
VBoxMachineStop(VBoxGlue glue) {
    return VBGlue_Ok;
}

#pragma mark - UTILS

VBGlueResult
VBoxTestErrorMessage(VBoxGlue glue) {
    ivbox_session* session = toiVBoxSession(glue);
    print_error_info(session->error_msg, "[VBox] VBoxGlue Error Message Test", NS_OK);
    return VBGlue_Fail;
}

const char*
VBoxGetErrorMessage(VBoxGlue glue) {
    return toiVBoxSession(glue)->error_msg;
}

const char*
VBoxGetSettingFilePath(VBoxGlue glue) {
    return toiVBoxSession(glue)->setting_file_path;
}

const char*
VBoxGetMachineID(VBoxGlue glue) {
    return toiVBoxSession(glue)->machine_id;
}
