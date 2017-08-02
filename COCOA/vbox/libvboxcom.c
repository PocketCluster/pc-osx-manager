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

typedef struct iVBoxSession {
    char                  error_msg[ERROR_MESSAGE_BUF_SIZE];
    IVirtualBox*          vbox;                             // virtualbox
    IVirtualBoxClient*    client;                           // vbox client
    ISession*             vsession;                         // vbox session
    IMachine*             machine;                          // vbox machine
    char*                 machine_id;                       // machine id
    char*                 setting_file_path;                // setting file
} iVBoxSession;

#pragma mark - MACROS

// these two are convert between internal <-> external types
#define toiVBoxSession(ptr) ((iVBoxSession*)ptr)
#define toVBoxGlue(ptr)     ((VBoxGlue)ptr)


#define CLIENT_DPTR(ptr)  ((IVirtualBoxClient**)ptr)     // Client double pointer
#define CLIENT_DREF(ptr)  ((IVirtualBoxClient*)(*(ptr))) // Client deref

#define SESSION_DPTR(ptr) ((ISession**)ptr)              // Session double pointer
#define SESSION_DREF(ptr) ((ISession*)(*(ptr)))          // Session deref

#define VBOX_DPTR(ptr)    ((IVirtualBox**)ptr)           // VirtualBox Double Pointer
#define VBOX_DREF(ptr)    ((IVirtualBox*)(*(ptr)))       // Virtualbox deref

#define MACHINE_DPTR(ptr) ((IMachine**)ptr)              // Machine Double Pointer
#define MACHINE_DREF(ptr) ((IMachine*)(*(ptr)))          // Machine deref

#pragma mark - DECLARATION

#pragma mark build machine base
typedef void** VOID_DPTR;

VBGlueResult vbox_machine_build(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, int cpu_count, int memory_size, char* error_message);

VBGlueResult vbox_machine_add_bridged_network(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* host_interface, char* error_message);

VBGlueResult vbox_machine_add_shared_folder(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* shared_name, const char *host_folder, char* error_message);


VBGlueResult vbox_machine_add_storage_controller(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, char* error_message);

VBGlueResult vbox_machine_add_boot_image(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, const char *boot_image_path, char *error_message);

VBGlueResult vbox_machine_add_hard_disk(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, const char *hdd_medium_path, int disk_size, void(^build_progress)(int progress, int done), char *error_message);


#pragma mark delete & release machine
VBGlueResult vbox_machine_destroy(VOID_DPTR vbox_machine, char* base_folder, const char* storage_controller_name, int remove_dvd, void(^build_progress)(int progress, int done), char *error_message);


#pragma mark utils
/**
 * Print detailed error information if available.
 * @param   pszErrorMsg     string containing the code location specific error message
 * @param   rc              COM/XPCOM result code
 */
static inline void
print_error_info(char *message_buffer, const char *pszErrorMsg, HRESULT rc) {
    memset(message_buffer, 0, ERROR_MESSAGE_BUF_SIZE);
    sprintf(message_buffer, "\n--- %s (rc=%#010x) ---\n", pszErrorMsg, rc);
}

#pragma mark - DEFINITION

#pragma mark init & close
VBGlueResult
NewVBoxGlue(VBoxGlue* glue) {

    // make sure the pointer passed IS null.
    assert(glue != NULL && *glue == NULL);
    
    HRESULT result;
    iVBoxSession* session = (iVBoxSession*)calloc(1, sizeof(iVBoxSession));
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
    
    iVBoxSession* session = toiVBoxSession(glue);
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

    iVBoxSession* session = toiVBoxSession(glue);
    PRBool changed = PR_FALSE;
    *isMachineChanged = (bool)changed;
    IMachine *mutable_machine;
    
    //firstly lock the machine
    HRESULT result = VboxLockMachine(session->machine, session->vsession, LockType_Write);
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] Failed to lock machine for adding storage controller", result);
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
            print_error_info(session->error_msg, "[VBox] Failed to release locked machine for attaching adapter", result);
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


#pragma mark find, build & destroy machine
VBGlueResult
VBoxFindMachineByNameOrID(VBoxGlue glue, const char* machine_name) {
    
    // make sure the pointer passed is not null.
    assert(glue != NULL);

    iVBoxSession* session = toiVBoxSession(glue);

    if ( session->machine  == NULL ) {
        HRESULT result;
        result = VboxFindMachine(session->vbox, machine_name, &(session->machine));
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
        print_error_info(session->error_msg, "[VBox] machine instance already exists", S_OK);
        return VBGlue_Fail;
    }
    return VBGlue_Ok;
}

VBGlueResult
VBoxCreateMachineByName(VBoxGlue glue, const char* base_folder, const char* machine_name) {
    
    // make sure the pointer passed is not null.
    assert(glue != NULL);
    assert(machine_name != NULL || strlen(machine_name) != 0);
    
    iVBoxSession* session = toiVBoxSession(glue);
    HRESULT result;

    if ( session->machine  != NULL ) {
        print_error_info(session->error_msg, "[VBox] machine instance already exists", S_OK);
        return VBGlue_Fail;
    }
    
    // create machine file name
    result = VboxComposeMachineFilename(session->vbox, machine_name, "", (char *)base_folder, &(session->setting_file_path));
    if (FAILED(result)) {
        print_error_info(session->error_msg, "[VBox] Failed composing machine name", result);
        return VBGlue_Fail;
    }
    // create machine based on the
    result = VboxCreateMachine(session->vbox, session->setting_file_path, machine_name, "Linux26_64", "", &(session->machine));
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

    iVBoxSession* session = toiVBoxSession(glue);

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


#pragma mark build machine base
VBGlueResult
vbox_machine_build(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, int cpu_count, int memory_size, char* error_message) {
    HRESULT result;
    VBGlueResult ret = VBGlue_Ok;

    assert(VBOX_DREF(virtualbox)      != NULL);
    assert(MACHINE_DREF(vbox_machine) != NULL);

    // Setup BIOS
    {
        // get BIOS settings
        IBIOSSettings *bios;
        result = VboxGetMachineBIOSSettings(MACHINE_DREF(vbox_machine), &bios);
        if (FAILED(result) || bios == NULL) {
            print_error_info(error_message, "[VBox] Failed to acquire bios settings", result);
            return VBGlue_Fail;
        }
        // enable I/O APIC
        result = IBIOSSettings_SetIOAPICEnabled(bios, (PRBool)1);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to enable IO/APIC", result);
                return VBGlue_Fail;
            }
        }
        // set ACPI enabled
        result = IBIOSSettings_SetACPIEnabled(bios, (PRBool)1);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to enable ACPI", result);
                return VBGlue_Fail;
            }
        }
        // release bios settings
        result = VboxIBiosSettingsRelease(bios);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to release BIOS", result);
                return VBGlue_Fail;
            }
        }
    }
    
    // Motherboard Settings
    {
        // set memory
        result = VboxSetMachineMemorySize(MACHINE_DREF(vbox_machine), memory_size);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBGlue_Error] Failed to set memory size", result);
            return VBGlue_Fail;
        }
        
        // set up Boot Order
        result = IMachine_SetBootOrder(MACHINE_DREF(vbox_machine), 1, DeviceType_DVD);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBGlue_Error] Failed to fix boot order", result);
            return VBGlue_Fail;
        }
        result = IMachine_SetBootOrder(MACHINE_DREF(vbox_machine), 2, DeviceType_HardDisk);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBGlue_Error] Failed to fix boot order", result);
            return VBGlue_Fail;
        }
        result = IMachine_SetBootOrder(MACHINE_DREF(vbox_machine), 3, DeviceType_Null);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to fix boot order", result);
                return VBGlue_Fail;
            }
        }
        result = IMachine_SetBootOrder(MACHINE_DREF(vbox_machine), 4, DeviceType_Null);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to fix boot order", result);
                return VBGlue_Fail;
            }
        }
        
        // set Chipset type
        result = IMachine_SetChipsetType(MACHINE_DREF(vbox_machine), ChipsetType_ICH9);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to setting chipset type", result);
                return VBGlue_Fail;
            }
        }
        // set RTC timer
        result = IMachine_SetRTCUseUTC(MACHINE_DREF(vbox_machine), (PRBool)1);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to setting Hardware UTC timer", result);
                return VBGlue_Fail;
            }
        }
    }
    
    // Processor Setting
    {
        // set CPU Count
        result = IMachine_SetCPUCount(MACHINE_DREF(vbox_machine), cpu_count);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBGlue_Error] Failed to setting CPU count", result);
            return VBGlue_Fail;
        }
        // set Execution Cap
        result = IMachine_SetCPUExecutionCap(MACHINE_DREF(vbox_machine), 100);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to setting CPU execution cap", result);
                return VBGlue_Fail;
            }
        }
        // PAE enabled
        PRBool enabled = (PRBool)1;
        result = IMachine_GetCPUProperty(MACHINE_DREF(vbox_machine), CPUPropertyType_PAE, &enabled);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to setting PAE/NX enabling", result);
                return VBGlue_Fail;
            }
        }
    }
    
    // Acceleration
    {
        // Paravirtualization setting
        result = IMachine_SetParavirtProvider(MACHINE_DREF(vbox_machine), ParavirtProvider_Default);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to setting Pravirtualization", result);
                return VBGlue_Fail;
            }
        }
        // Nested Paging
        result = IMachine_SetHWVirtExProperty(MACHINE_DREF(vbox_machine), HWVirtExPropertyType_Enabled, (PRBool)1);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to setting HWVirtExPropertyType", result);
                return VBGlue_Fail;
            }
        }
    }
    
    // Display
    {
        // set VRAM
        result = VboxSetMachineVRAMSize(MACHINE_DREF(vbox_machine), 12);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to VRAM size", result);
                return VBGlue_Fail;
            }
        }
    }
    
    // SAVE SETTINGS & REGISTER MACHINE BEFORE ATTACH A MEDIUM
    {
        // save settings
        result = VboxMachineSaveSettings(MACHINE_DREF(vbox_machine));
        if (FAILED(result)) {
            print_error_info(error_message, "[VBGlue_Error] Failed to save machine before attaching a medium", result);
            return VBGlue_Fail;
        }
        // Register machine
        result = VboxRegisterMachine(VBOX_DREF(virtualbox), MACHINE_DREF(vbox_machine));
        if (FAILED(result)) {
            print_error_info(error_message, "[VBGlue_Error] Failed to register machine", result);
            return VBGlue_Fail;
        }
    }
    return ret;
}

VBGlueResult
vbox_machine_add_bridged_network(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* host_interface, char* error_message) {

    INetworkAdapter *adapter = NULL;
    VBGlueResult ret = VBGlue_Ok;
    
    //firstly lock the machine
    HRESULT result = VboxLockMachine(MACHINE_DREF(vbox_machine), SESSION_DREF(vbox_session), LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBGlue_Error] Failed to lock machine for networking", result);
        return VBGlue_Fail;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(SESSION_DREF(vbox_session), &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[VBGlue_Error] Failed to get a mutable copy of a machine for networking", result);
        ret = VBGlue_Fail;
    }
    // get network adapter
    result = VboxMachineGetNetworkAdapter(mutable_machine, 0, &adapter);
    if (FAILED(result) || adapter == NULL) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to acquire adapter from slot 0", result);
            ret = VBGlue_Fail;
        }
    }
    // enable network adapter
    result = VboxNetworkAdapterSetEnabled(adapter, TRUE);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to enable network adapter", result);
            ret = VBGlue_Fail;
        }
    }
    // set bridged network type
    result = VboxNetworkAdapterSetAttachmentType(adapter, NetworkAttachmentType_Bridged);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to set network attachement type", result);
            ret = VBGlue_Fail;
        }
    }

    // set host network adapter this bridge should connect to
    result = VboxNetworkAdapterSetBridgedHostInterface(adapter, host_interface);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBGlue_Error] Failed to connect to host network interface", result);
        ret = VBGlue_Fail;
    }    
    // set adapter type (AMD PCnet-FAST III, VBox Default)
    result = VboxNetworkAdapterSetAdapterType(adapter, NetworkAdapterType_Am79C973);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to set network adapter type", result);
            ret = VBGlue_Fail;
        }
    }
    // promiscuous mode policy
    result = VboxNetworkAdapterSetPromiscModePolicy(adapter, NetworkAdapterPromiscModePolicy_Deny);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to set promiscuous mode", result);
            return VBGlue_Fail;
        }
    }
    // set cable connected
    result = VboxNetworkAdapterSetCableConnected(adapter, TRUE);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to set cable connected", result);
            return VBGlue_Fail;
        }
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to save machine after attaching hard disk medium", result);
            ret = VBGlue_Fail;
        }
    }
    // release the first adapter
    result = VboxNetworkAdapterRelease(adapter);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to release adapter", result);
            return VBGlue_Fail;
        }
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to release locked machine for attaching adapter", result);
                return VBGlue_Fail;
            }
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(SESSION_DREF(vbox_session));
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to unlock machine for attaching adapter", result);
            return VBGlue_Fail;
        }
    }
    
    return ret;
}

VBGlueResult
vbox_machine_add_shared_folder(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* shared_name, const char *host_folder, char* error_message) {

    VBGlueResult ret = VBGlue_Ok;
    HRESULT result;
    
    //firstly lock the machine
    result = VboxLockMachine(MACHINE_DREF(vbox_machine), SESSION_DREF(vbox_session), LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBGlue_Error] Failed to lock machine for shared folder", result);
        return VBGlue_Fail;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(SESSION_DREF(vbox_session), &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[VBGlue_Error] Failed to get a mutable copy of a machine for shared folder", result);
        ret = VBGlue_Fail;
    }
    // create shared folder
    result = VboxMachineCreateSharedFolder(mutable_machine, shared_name, host_folder, TRUE, TRUE);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to add shared folder", result);
            ret = VBGlue_Fail;
        }
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to save machine for adding shared folder", result);
            ret = VBGlue_Fail;
        }
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to release machine for adding shared folder", result);
                ret = VBGlue_Fail;
            }
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(SESSION_DREF(vbox_session));
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to unlock machine for adding shared folder", result);
            return VBGlue_Fail;
        }
    }
    return ret;
}

VBGlueResult
vbox_machine_add_storage_controller(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, char* error_message) {

    VBGlueResult ret = VBGlue_Ok;
    HRESULT result;
    IStorageController *storage_controller;
    
    //firstly lock the machine
    result = VboxLockMachine(MACHINE_DREF(vbox_machine), SESSION_DREF(vbox_session), LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBGlue_Error] Failed to lock machine for adding storage controller", result);
        return VBGlue_Fail;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(SESSION_DREF(vbox_session), &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[VBGlue_Error] Failed to get a mutable copy of a machine", result);
        ret = VBGlue_Fail;
    }
    // add storage controller
    result = VboxMachineAddStorageController(mutable_machine, storage_controller_name, StorageBus_SATA, &storage_controller);
    if (FAILED(result) || storage_controller_name == NULL) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to add storage controller", result);
            ret = VBGlue_Fail;
        }
    }
    // storage controller type
    result = VboxSetStorageControllerType(storage_controller, StorageControllerType_IntelAhci);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to set storage controller type", result);
            return VBGlue_Fail;
        }
    }
    // storage controller set # of ports
    result = IStorageController_SetPortCount(storage_controller, 10);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to increase port count", result);
            ret = VBGlue_Fail;
        }
    }
    // Enable host IO cache for imaging
    PRBool use_host_iocache = (PRBool)1;
    result = IStorageController_SetUseHostIOCache(storage_controller, use_host_iocache);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to enable host IO cache", result);
            return VBGlue_Fail;
        }
    }
    // release storage controller
    result = VboxIStorageControllerRelease(storage_controller);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to release storage controller for attaching storage controller", result);
            ret = VBGlue_Fail;
        }
        
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to save machine after attaching storage controller", result);
            ret = VBGlue_Fail;
        }
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to release mutable machine for attaching storage controller", result);
                return VBGlue_Fail;
            }
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(SESSION_DREF(vbox_session));
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to unlock machine for attaching storage controller", result);
            return VBGlue_Fail;
        }
    }

    return ret;
}

VBGlueResult
vbox_machine_add_boot_image(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, const char *boot_image_path, char *error_message) {
    
    IMedium *boot_image;
    VBGlueResult ret = VBGlue_Ok;
    HRESULT result;

    // open medium
    PRBool use_old_uuid = (PRBool)0;
    result = VboxOpenMedium(VBOX_DREF(virtualbox), boot_image_path, DeviceType_DVD, AccessMode_ReadOnly, use_old_uuid, &boot_image);
    if (FAILED(result) || boot_image == NULL) {
        print_error_info(error_message, "[VBGlue_Error] Failed to open boot image for attaching boot image", result);
        return VBGlue_Fail;
    }
    //firstly lock the machine
    result = VboxLockMachine(MACHINE_DREF(vbox_machine), SESSION_DREF(vbox_session), LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[VBGlue_Error] Failed to lock machine for attaching boot image", result);
        return VBGlue_Fail;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(SESSION_DREF(vbox_session), &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[VBGlue_Error] Failed to get a mutable copy of a machine", result);
        ret = VBGlue_Fail;
    }
    // attach a medium
    result = VboxMachineAttachDevice(mutable_machine, storage_controller_name, 0, 0, DeviceType_DVD, boot_image);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to attach boot image", result);
            ret = VBGlue_Fail;
        }
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBGlue_Error] Failed to save machine after attaching boot image", result);
            ret = VBGlue_Fail;
        }
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to release machine after attaching boot image", result);
                return VBGlue_Fail;
            }
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(SESSION_DREF(vbox_session));
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to unlock machine for attaching boot image", result);
            return VBGlue_Fail;
        }
    }
    // release medium
    result = VboxIMediumRelease(boot_image);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to release boot image", result);
            return VBGlue_Fail;
        }
    }
    return ret;
}

VBGlueResult
vbox_machine_add_hard_disk(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, const char *hdd_medium_path, int disk_size, void(^build_progress)(int progress, int done), char *error_message) {

    // Create and Open hard drive
    HRESULT result;
    VBGlueResult ret = VBGlue_Ok;
    IMedium *hdd_medium;
    {
        // set medium image
        result = VboxCreateHardDisk(VBOX_DREF(virtualbox), "VMDK", hdd_medium_path, DeviceType_HardDisk, AccessMode_ReadWrite, &hdd_medium);
        if (FAILED(result) || hdd_medium == NULL) {
            print_error_info(error_message, "[VBGlue_Error] Failed to create harddrive", result);
            return VBGlue_Fail;
        }
        // create medium
        //REF : https://www.virtualbox.org/sdkref/_virtual_box_8idl.html#adedcbf1a6e5e35fe7a0ca0c4b3447154
        PRUint32 cVariant[2] = {MediumVariant_Standard, MediumVariant_NoCreateDir};
        PRUint32 variantCount = sizeof(cVariant) / sizeof(cVariant[0]);
        IProgress *progress;
        result = VboxMediumCreateBaseStorage(hdd_medium, (1 << 24), variantCount, cVariant, &progress);
        if (FAILED(result)){
            print_error_info(error_message, "[VBGlue_Error] Failed to create base storage", result);
            return VBGlue_Fail;
        }
        
        // it is recommended to wait short amount of time
        result = VboxProgressWaitForCompletion(progress, 3);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBGlue_Error] Failed to complete creating base storage", result);
            ret = VBGlue_Fail;
        }
        else {
            PRUint32 progress_percent = 0;
            do {
                VboxGetProgressPercent(progress, &progress_percent);
                if (build_progress != NULL) {
                    build_progress((int)progress_percent, 0);
                }
                usleep(500000);
            } while (progress_percent < 100);
        }
        
        // get completion code
        PRInt32 code;
        result = VboxGetProgressResultCode(progress, &code);
        if (FAILED(result)|| code != 0) {
            print_error_info(error_message, "[VBGlue_Error] Failed to actuqire storage creation result code", result);
            return VBGlue_Fail;
        }
        // release progress
        VboxIProgressRelease(progress);
        
        // open medium
        PRBool use_old_uuid = (PRBool)0;
        result = VboxOpenMedium(VBOX_DREF(virtualbox), hdd_medium_path, DeviceType_HardDisk, AccessMode_ReadWrite, use_old_uuid, &hdd_medium);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBGlue_Error] Failed to open hard drive", result);
            return VBGlue_Fail;
        }
    }
    
    // Attach a medium to storage controller
    {
        //firstly lock the machine
        result = VboxLockMachine(MACHINE_DREF(vbox_machine), SESSION_DREF(vbox_session), LockType_Write);
        if (FAILED(result)) {
            print_error_info(error_message, "[VBGlue_Error] Failed to lock machine for attaching hdd to storage controller", result);
            return VBGlue_Fail;
        }
        // get mutable machine
        IMachine *mutable_machine;
        result = VboxGetSessionMachine(SESSION_DREF(vbox_session), &mutable_machine);
        if (FAILED(result) || mutable_machine == NULL) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBGlue_Error] Failed to get a mutable copy of a machine for attaching hard disk medium", result);
                ret = VBGlue_Fail;
            }
        }
        // attach a medium
        result = VboxMachineAttachDevice(mutable_machine, storage_controller_name, 1, 0, DeviceType_HardDisk, hdd_medium);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBGlue_Error] Failed to attach hard disk medium", result);
                ret = VBGlue_Fail;
            }
        }
        // save setting
        result = VboxMachineSaveSettings(mutable_machine);
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBGlue_Error] Failed to save machine after attaching hard disk medium", result);
                ret = VBGlue_Fail;
            }
        }
        // then we can safely release the mutable machine
        if (mutable_machine) {
            result = VboxIMachineRelease(mutable_machine);
            if (FAILED(result)) {
                if (ret == VBGlue_Ok) {
                    print_error_info(error_message, "[VBox] Failed to release machine after attaching hard disk medium", result);
                    return VBGlue_Fail;
                }
            }
        }
        // then unlock machine
        result = VboxUnlockMachine(SESSION_DREF(vbox_session));
        if (FAILED(result)) {
            if (ret == VBGlue_Ok) {
                print_error_info(error_message, "[VBox] Failed to unlock machine after attaching hard disk medium", result);
                return VBGlue_Fail;
            }
        }
    }
    
    // Close & release hard drive
    result = VboxIMediumRelease(hdd_medium);
    if (FAILED(result)) {
        if (ret == VBGlue_Ok) {
            print_error_info(error_message, "[VBox] Failed to release a hard drive", result);
            return VBGlue_Fail;
        }
    }
    
    // report the end of result
    if (build_progress != NULL) {
        build_progress(100, 1);
    }
    return ret;
}

#pragma mark destroy machine
VBGlueResult
VBoxDestoryMachine(VBoxGlue glue) {

    // make sure the pointer passed is not null.
    assert(glue != NULL);
    
    iVBoxSession* session = toiVBoxSession(glue);
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

/**
 * Register passive event listener for the selected VM.
 *
 * @param   virtualBox ptr to IVirtualBox object
 * @param   session    ptr to ISession object
 * @param   id         identifies the machine to start
 */
static void
registerPassiveEventListener(IVirtualBox *virtualBox, ISession *session, BSTR machineId) {
#if 0
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
#endif
}

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
    iVBoxSession* session = toiVBoxSession(glue);
    print_error_info(session->error_msg, "[VBox] VBoxGlue Error Message Test", (unsigned)S_OK);
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