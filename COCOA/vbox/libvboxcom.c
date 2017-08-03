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


#pragma mark machine status
VBGlueResult
VBoxIsMachineSettingChanged(VBoxGlue glue, bool* isMachineChanged) {
    
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


#pragma mark find, create, & build machine
VBGlueResult
VBoxFindMachineByNameOrID(VBoxGlue glue, const char* machineName) {
    
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
VBoxCreateMachineByName(VBoxGlue glue, const char* baseFolder, const char* machineName) {
    
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

VBGlueResult
VBoxReleaseMachine(VBoxGlue glue) {
    
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
                VboxGetProgressPercent(progress, &progress_percent);
                usleep(500000);
            } while (progress_percent < 100);
        }
        
        // get completion code
        PRInt32 code;
        result = VboxGetProgressResultCode(progress, &code);
        if (FAILED(result)|| code != 0) {
            print_error_info(error_message, "[VBox] Failed to actuqire storage creation result code", result);
            return result;
        }
        // release progress
        VboxIProgressRelease(progress);
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
VBoxMakeBuildOption(int cpu, int mem, const char* host, const char* spath, const char* sname, const char* boot, const char* hdd) {
    VBoxBuildOption* option = (VBoxBuildOption*)calloc(1, sizeof(VBoxBuildOption));
    option->CpuCount      = cpu;
    option->MemSize       = mem;
    option->HostInterface = host;
    option->SharedDirPath = spath;
    option->SharedDirName = sname;
    option->BootImagePath = boot;
    option->HddImagePath  = hdd;
    return option;
}

VBGlueResult
VBoxBuildMachine(VBoxGlue glue, VBoxBuildOption* option) {
    
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
VBoxDestoryMachine(VBoxGlue glue) {

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
        VboxGetProgressPercent(progress, &progress_percent);
        usleep(500000);
    } while (progress_percent < 100);
    VboxIProgressRelease(progress);
    
    // free media array
    VboxArrayOutFree(media);

    // release machine
    return VBoxReleaseMachine(glue);
}


#pragma mark start & stop machine

HRESULT
VboxMachineStart(IVirtualBox *virtualBox, ISession *session, IMachine *cmachine, const char* session_type) {
    HRESULT rc;
    IMachine  *machine    = NULL;
    IProgress *progress   = NULL;
    BSTR env              = NULL;
    BSTR sessionType;
    SAFEARRAY *groupsSA = g_pVBoxFuncs->pfnSafeArrayOutParamAlloc();
    
    rc = IMachine_get_Groups(machine, ComSafeArrayAsOutTypeParam(groupsSA, BSTR));
    if (SUCCEEDED(rc)) {
        BSTR *groups = NULL;
        ULONG cbGroups = 0;
        ULONG i, cGroups;
        g_pVBoxFuncs->pfnSafeArrayCopyOutParamHelper((void **)&groups, &cbGroups, VT_BSTR, groupsSA);
        g_pVBoxFuncs->pfnSafeArrayDestroy(groupsSA);
        cGroups = cbGroups / sizeof(groups[0]);
        for (i = 0; i < cGroups; ++i) {
            /* Note that the use of %S might be tempting, but it is not
             * available on all platforms, and even where it is usable it
             * may depend on correct compiler options to make wchar_t a
             * 16 bit number. So better play safe and use UTF-8. */
            char *group;
            g_pVBoxFuncs->pfnUtf16ToUtf8(groups[i], &group);
            g_pVBoxFuncs->pfnUtf8Free(group);
        }
        for (i = 0; i < cGroups; ++i) {
            g_pVBoxFuncs->pfnComUnallocString(groups[i]);
        }
        g_pVBoxFuncs->pfnArrayOutFree(groups);
    }
    
    g_pVBoxFuncs->pfnUtf8ToUtf16("gui", &sessionType);
    rc = IMachine_LaunchVMProcess(machine, session, sessionType, env, &progress);
    g_pVBoxFuncs->pfnUtf16Free(sessionType);
    if (SUCCEEDED(rc)) {
        
        BOOL completed;
        LONG resultCode;
        
        printf("Waiting for the remote session to open...\n");
        IProgress_WaitForCompletion(progress, -1);
        
        rc = IProgress_get_Completed(progress, &completed);
        if (FAILED(rc)) {
            fprintf(stderr, "Error: GetCompleted status failed\n");
        }
        
        IProgress_get_ResultCode(progress, &resultCode);
        if (FAILED(resultCode)) {
            IVirtualBoxErrorInfo *errorInfo;
            BSTR textUtf16;
            char *text;
            
            IProgress_get_ErrorInfo(progress, &errorInfo);
            IVirtualBoxErrorInfo_get_Text(errorInfo, &textUtf16);
            g_pVBoxFuncs->pfnUtf16ToUtf8(textUtf16, &text);
            printf("Error: %s\n", text);
            
            g_pVBoxFuncs->pfnComUnallocString(textUtf16);
            g_pVBoxFuncs->pfnUtf8Free(text);
            IVirtualBoxErrorInfo_Release(errorInfo);
        } else {
            fprintf(stderr, "VM process has been successfully started\n");
            
            /* Kick off the event listener demo part, which is quite separate.
             * Ignore it if you need a more basic sample. */
            //registerPassiveEventListener(virtualBox, session, id);
        }
        IProgress_Release(progress);
    }
    
    /* It's important to always release resources. */
    //IMachine_Release(machine);
    return rc;
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
VboxGetSettingFilePath(VBoxGlue glue) {
    return toiVBoxSession(glue)->setting_file_path;
}

const char*
VboxGetMachineID(VBoxGlue glue) {
    return toiVBoxSession(glue)->machine_id;
}


# if 0
/**
 * Register passive event listener for the selected VM.
 *
 * @param   virtualBox ptr to IVirtualBox object
 * @param   session    ptr to ISession object
 * @param   id         identifies the machine to start
 */
static void
registerPassiveEventListener(IVirtualBox *virtualBox, ISession *session, BSTR machineId) {
    
    IConsole *console = NULL;
    HRESULT rc;
    
    rc = ISession_get_Console(session, &console);
    if ((SUCCEEDED(rc)) && console) {
        
        IEventSource *es = NULL;
        rc = IConsole_get_EventSource(console, &es);
        
        if (SUCCEEDED(rc) && es) {
            static const ULONG interestingEvents[] = {
                VBoxEventType_OnMousePointerShapeChanged,
                VBoxEventType_OnMouseCapabilityChanged,
                VBoxEventType_OnKeyboardLedsChanged,
                VBoxEventType_OnStateChanged,
                VBoxEventType_OnAdditionsStateChanged,
                VBoxEventType_OnNetworkAdapterChanged,
                VBoxEventType_OnSerialPortChanged,
                VBoxEventType_OnParallelPortChanged,
                VBoxEventType_OnStorageControllerChanged,
                VBoxEventType_OnMediumChanged,
                VBoxEventType_OnVRDEServerChanged,
                VBoxEventType_OnUSBControllerChanged,
                VBoxEventType_OnUSBDeviceStateChanged,
                VBoxEventType_OnSharedFolderChanged,
                VBoxEventType_OnRuntimeError,
                VBoxEventType_OnCanShowWindow,
                VBoxEventType_OnShowWindow
            };
            SAFEARRAY *interestingEventsSA = NULL;
            IEventListener *consoleListener = NULL;
            
            /* The VirtualBox API expects enum values as VT_I4, which in the
             * future can be hopefully relaxed. */
            interestingEventsSA = g_pVBoxFuncs->pfnSafeArrayCreateVector(VT_I4, 0, sizeof(interestingEvents) / sizeof(interestingEvents[0]));
            g_pVBoxFuncs->pfnSafeArrayCopyInParamHelper(interestingEventsSA, &interestingEvents, sizeof(interestingEvents));
            
            rc = IEventSource_CreateListener(es, &consoleListener);
            if (SUCCEEDED(rc) && consoleListener) {
                rc = IEventSource_RegisterListener(es,
                                                   consoleListener,
                                                   ComSafeArrayAsInParam(interestingEventsSA),
                                                   0 /* passive */);
                
                if (SUCCEEDED(rc)) {
                    
                    /* Just wait here for events, no easy way to do this better
                     * as there's not much to do after this completes. */
                    printf("Entering event loop, PowerOff the machine to exit or press Ctrl-C to terminate\n");
                    fflush(stdout);
                    signal(SIGINT, (void (*)(int))ctrlCHandler);
                    
                    while (!g_fStop) {
                        IEvent *ev = NULL;
                        rc = IEventSource_GetEvent(es, consoleListener, 250, &ev);
                        if (FAILED(rc)) {
                            printf("Failed getting event: %#x\n", rc);
                            g_fStop = 1;
                            continue;
                        }
                        /* handle timeouts, resulting in NULL events */
                        if (!ev) {
                            continue;
                        }
                        rc = EventListenerDemoProcessEvent(ev);
                        if (FAILED(rc)) {
                            printf("Failed processing event: %#x\n", rc);
                            g_fStop = 1;
                            /* finish processing the event */
                        }
                        rc = IEventSource_EventProcessed(es, consoleListener, ev);
                        if (FAILED(rc)) {
                            printf("Failed to mark event as processed: %#x\n", rc);
                            g_fStop = 1;
                            /* continue with event release */
                        }
                        if (ev) {
                            IEvent_Release(ev);
                            ev = NULL;
                        }
                    }
                    signal(SIGINT, SIG_DFL);
                } else {
                    printf("Failed to register event listener.\n");
                }
                IEventSource_UnregisterListener(es, (IEventListener *)consoleListener);
                IEventListener_Release(consoleListener);
            } else {
                printf("Failed to create an event listener instance.\n");
            }
            g_pVBoxFuncs->pfnSafeArrayDestroy(interestingEventsSA);
            IEventSource_Release(es);
        } else {
            printf("Failed to get the event source instance.\n");
        }
        IConsole_Release(console);
    }
}

/**
 * Print detailed error information if available.
 * @param   pszErrorMsg     string containing the code location specific error message
 * @param   rc              COM/XPCOM result code
 */
void
print_error_info(char *message_buffer, const char *pszErrorMsg, HRESULT rc)
{
    IErrorInfo *ex;
    HRESULT rc2 = S_OK;
    
    sprintf(message_buffer, "\n--- %s (rc=%#010x) ---\n", pszErrorMsg, (unsigned)rc);
    rc2 = g_pVBoxFuncs->pfnGetException(&ex);
    
    if (SUCCEEDED(rc2) && ex) {
        
        IVirtualBoxErrorInfo *ei;
        rc2 = IErrorInfo_QueryInterface(ex, &IID_IVirtualBoxErrorInfo, (void **)&ei);
        
        if (FAILED(rc2)) {
            ei = NULL;
        }
        
        if (ei) {
            /* got extended error info, maybe multiple infos */
            do {
                LONG resultCode = S_OK;
                BSTR componentUtf16 = NULL;
                char *component = NULL;
                BSTR textUtf16 = NULL;
                char *text = NULL;
                IVirtualBoxErrorInfo *ei_next = NULL;
                sprintf(message_buffer, "Extended error info (IVirtualBoxErrorInfo):\n");
                
                IVirtualBoxErrorInfo_get_ResultCode(ei, &resultCode);
                sprintf(message_buffer, "  resultCode=%#010x\n", (unsigned)resultCode);
                
                IVirtualBoxErrorInfo_get_Component(ei, &componentUtf16);
                g_pVBoxFuncs->pfnUtf16ToUtf8(componentUtf16, &component);
                g_pVBoxFuncs->pfnComUnallocString(componentUtf16);
                sprintf(message_buffer, "  component=%s\n", component);
                g_pVBoxFuncs->pfnUtf8Free(component);
                
                IVirtualBoxErrorInfo_get_Text(ei, &textUtf16);
                g_pVBoxFuncs->pfnUtf16ToUtf8(textUtf16, &text);
                g_pVBoxFuncs->pfnComUnallocString(textUtf16);
                sprintf(message_buffer, "  text=%s\n", text);
                g_pVBoxFuncs->pfnUtf8Free(text);
                
                rc2 = IVirtualBoxErrorInfo_get_Next(ei, &ei_next);
                if (FAILED(rc2))
                    ei_next = NULL;
                IVirtualBoxErrorInfo_Release(ei);
                ei = ei_next;
            } while (ei);
        }
        
        IErrorInfo_Release(ex);
        g_pVBoxFuncs->pfnClearException();
    }
}

// You must free the result if result is non-NULL.
static char*
str_replace(char* orig, const char* rep, const char* with) {
    char *result;         // the return string
    char *ins;            // the next insert point
    char *tmp;            // varies
    size_t len_rep;       // length of rep
    size_t len_with;      // length of with
    size_t len_front;     // distance between rep and end of last rep
    size_t count;         // number of replacements
    
    if (!orig) {
        return (char *)NULL;
    }
    if (!rep) {
        rep = "";
    }
    len_rep = strlen(rep);
    if (!with) {
        with = "";
    }
    len_with = strlen(with);
    
    ins = orig;
    for (count = 0; (tmp = strstr(ins, rep)); ++count) {
        ins = tmp + len_rep;
    }
    
    // first time through the loop, all the variable are set correctly
    // from here on,
    //    tmp points to the end of the result string
    //    ins points to the next occurrence of rep in orig
    //    orig points to the remainder of orig after "end of rep"
    tmp = result = malloc(strlen(orig) + (len_with - len_rep) * count + 1);
    
    if (!result) {
        return NULL;
    }
    while (count--) {
        ins = strstr(orig, rep);
        len_front = ins - orig;
        tmp = strncpy(tmp, orig, len_front) + len_front;
        tmp = strcpy(tmp, with) + len_with;
        orig += len_front + len_rep; // move to next "end of rep"
    }
    strcpy(tmp, orig);
    return result;
}
#endif