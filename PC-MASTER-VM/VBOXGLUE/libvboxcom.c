//
//  libvboxcom.c
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/12/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

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
#include <assert.h>
#include <unistd.h>

#include "libvboxcom.h"

#pragma mark - MACRO

#define CLIENT_DPTR(ptr)  ((IVirtualBoxClient**)ptr)     // Client double pointer
#define CLIENT_DREF(ptr)  ((IVirtualBoxClient*)(*(ptr))) // Client deref

#define SESSION_DPTR(ptr) ((ISession**)ptr)              // Session double pointer
#define SESSION_DREF(ptr) ((ISession*)(*(ptr)))          // Session deref

#define VBOX_DPTR(ptr)    ((IVirtualBox**)ptr)           // VirtualBox Double Pointer
#define VBOX_DREF(ptr)    ((IVirtualBox*)(*(ptr)))       // Virtualbox deref

#define MACHINE_DPTR(ptr) ((IMachine**)ptr)              // Machine Double Pointer
#define MACHINE_DREF(ptr) ((IMachine*)(*(ptr)))          // Machine deref

#pragma mark - UTIL

/**
 * Print detailed error information if available.
 * @param   pszErrorMsg     string containing the code location specific error message
 * @param   rc              COM/XPCOM result code
 */
static void
print_error_info(char *message_buffer, const char *pszErrorMsg, HRESULT rc)
{
    sprintf(message_buffer, "\n--- %s (rc=%#010x) ---\n", pszErrorMsg, (unsigned)S_OK);
#if 0
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
#endif
}

# if 0
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

#pragma mark - APP & API VERSION

unsigned int
vbox_app_version() {
    return VboxGetAppVersion();
}

#pragma mark - GET MACHINE ID
VBRESULT
vbox_machine_getid(VOID_DPTR vbox_machine, char** machine_id, char *error_message) {
    HRESULT result = VboxMachineGetID(MACHINE_DREF(vbox_machine), machine_id);
    if (FAILED(result)) {
        print_error_info(error_message, "[INFO] Failed to get Machine ID", result);
        return FATAL;
    }
    return GOOD;
}


#pragma mark - MACHINE STATUS
VBRESULT
vbox_machine_setting_path(VOID_DPTR vbox_machine, char** base_folder, char* error_message) {
    HRESULT result = VboxGetMachineSettingsFilePath(MACHINE_DREF(vbox_machine), base_folder);
    if (FAILED(result)) {
        print_error_info(error_message, "[INFO] Fail to get setting path", result);
        return INFO;
    }
    return GOOD;
}

int
vbox_machine_is_setting_changed(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, char *error_message) {

    PRBool changed = (PRBool)0;
    IMachine *mutable_machine;
    
    //firstly lock the machine
    HRESULT result = VboxLockMachine(MACHINE_DREF(vbox_machine), SESSION_DREF(vbox_session), LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[INFO] Failed to lock machine for adding storage controller", result);
    }
    // get mutable machine
    result = VboxGetSessionMachine(SESSION_DREF(vbox_session), &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[INFO] Failed to get a mutable copy of a machine", result);
    }
    // check if settings modified
    result = VboxGetMachineSettingsModified(mutable_machine, &changed);
    if (FAILED(result)) {
        print_error_info(error_message, "[INFO] Fail to get setting modification", result);
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            print_error_info(error_message, "[INFO] Failed to release locked machine for attaching adapter", result);
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(SESSION_DREF(vbox_session));
    if (FAILED(result)) {
        print_error_info(error_message, "[INFO] Failed to unlock machine for attaching adapter", result);
    }
    return changed;
}


#pragma mark - INIT & CLOSE
VBRESULT
vbox_init(char* error_message) {
    HRESULT result;
    result = VboxInit();
    if (FAILED(result)) {
        print_error_info(error_message, "[FATAL] Failed to initialize VBOX", result);
        return FATAL;
    }
    return GOOD;
}

void
vbox_term() {
    VboxTerm();
}


#pragma mark - SESSION INIT
VBRESULT
vbox_session_init(VOID_DPTR vbox_client, VOID_DPTR vbox_session, VOID_DPTR virtualbox, char* error_message) {
    HRESULT result;
    VboxClientInitialize(CLIENT_DPTR(vbox_client));
    if (!vbox_client || CLIENT_DREF(vbox_client) == NULL) {
        print_error_info(error_message, "[FATAL] Failed to get VirtualBoxClient reference", result);
        return FATAL;
    }
    result = VboxGetVirtualBox(CLIENT_DREF(vbox_client), VBOX_DPTR(virtualbox));
    if (FAILED(result) || VBOX_DREF(virtualbox) == NULL) {
        print_error_info(error_message, "[FATAL] Failed to get VirtualBox reference", result);
        return FATAL;
    }
    result = VboxGetSession(CLIENT_DREF(vbox_client), SESSION_DPTR(vbox_session));
    if (FAILED(result) || SESSION_DREF(vbox_session) == NULL) {
        print_error_info(error_message, "[FATAL] Failed to get Session reference", result);
        return FATAL;
    }
    return GOOD;
}

VBRESULT
vbox_session_close(VOID_DPTR vbox_client, VOID_DPTR vbox_session, VOID_DPTR virtualbox, char* error_message) {
    VBRESULT ret = GOOD;
    HRESULT result;
    if (SESSION_DREF(vbox_session) != NULL) {
        result = VboxISessionRelease(SESSION_DREF(vbox_session));
        if (FAILED(result)) {
            print_error_info(error_message, "[INFO] Failed to release ISession reference", result);
            ret = INFO;
        }
    }
    if (VBOX_DREF(virtualbox) != NULL) {
        result = VboxIVirtualBoxRelease(VBOX_DREF(virtualbox));
        if (FAILED(result)) {
            print_error_info(error_message, "[INFO] Failed to release IVirtualBox reference", result);
            ret = INFO;
        }
    }
    if (CLIENT_DREF(vbox_client) != NULL) {
        result = VboxClientRelease(CLIENT_DREF(vbox_client));
        if (FAILED(result)) {
            print_error_info(error_message, "[INFO] Failed to release IVirtualBoxClient reference", result);
            ret = INFO;
        }
    }
    VboxClientUninitialize();
    return ret;
}


#pragma mark - FIND, BUILD & DESTROY MACHINE
VBRESULT
vbox_machine_find(VOID_DPTR virtualbox, VOID_DPTR vbox_machine , const char* machine_name, char* error_message) {

    VBRESULT result = VboxFindMachine(VBOX_DREF(virtualbox), machine_name, MACHINE_DPTR(vbox_machine));
    if (FAILED(result)) {
        print_error_info(error_message, "[INFO] Failed to find machine", result);
        return INFO;
    }
    return GOOD;
}

VBRESULT
vbox_machine_create(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, const char* machine_name, char** base_folder, char* error_message) {
    HRESULT result;

    assert(machine_name != NULL || strlen(machine_name) != 0);
    assert(VBOX_DREF(virtualbox) != NULL);

    // create machine file name
    result = VboxComposeMachineFilename(VBOX_DREF(virtualbox), machine_name, "", "", base_folder);
    if (FAILED(result)) {
        print_error_info(error_message, "[FATAL] Failed composing machine name", result);
        return FATAL;
    }
    // create machine based on the
    result = VboxCreateMachine(VBOX_DREF(virtualbox), *base_folder, machine_name, "Linux26_64", "", MACHINE_DPTR(vbox_machine));
    if (FAILED(result) || MACHINE_DREF(vbox_machine) == NULL) {
        print_error_info(error_message, "[FATAL] Failed to create machine", result);
        return FATAL;
    }
    return GOOD;
}

VBRESULT
vbox_machine_release(VOID_DPTR vbox_machine, char* base_folder, char* error_message) {
    HRESULT result = GOOD;
    // release machine
    if (MACHINE_DREF(vbox_machine) != NULL) {
        HRESULT result = VboxIMachineRelease(MACHINE_DREF(vbox_machine));
        if (FAILED(result)) {
            print_error_info(error_message, "[INFO] Failed to close machine referenece", result);
        }
    }
    // release base folder
    if (base_folder != NULL) {
        VboxUtf8Free(base_folder);
    }
    
    return result;
}


#pragma mark - BUILD MACHINE BASE
VBRESULT
vbox_machine_build(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, int cpu_count, int memory_size, char* error_message) {
    HRESULT result;
    VBRESULT ret = GOOD;
    
    // Setup BIOS
    {
        // get BIOS settings
        IBIOSSettings *bios;
        result = VboxGetMachineBIOSSettings(MACHINE_DREF(vbox_machine), &bios);
        if (FAILED(result) || bios == NULL) {
            print_error_info(error_message, "[INFO] Failed to acquire bios settings", result);
            ret = INFO;
        }
        // enable I/O APIC
        result = IBIOSSettings_SetIOAPICEnabled(bios, (PRBool)1);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to enable IO/APIC", result);
                ret = INFO;
            }
        }
        // set ACPI enabled
        result = IBIOSSettings_SetACPIEnabled(bios, (PRBool)1);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to enable ACPI", result);
                ret = INFO;
            }
        }
        // release bios settings
        result = VboxIBiosSettingsRelease(bios);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to release BIOS", result);
                ret = INFO;
            }
        }
    }
    
    // Motherboard Settings
    {
        // set memory
        result = VboxSetMachineMemorySize(MACHINE_DREF(vbox_machine), memory_size);
        if (FAILED(result)) {
            print_error_info(error_message, "[FATAL] Failed to set memory size", result);
            return FATAL;
        }
        
        // set up Boot Order
        result = IMachine_SetBootOrder(MACHINE_DREF(vbox_machine), 1, DeviceType_DVD);
        if (FAILED(result)) {
            print_error_info(error_message, "[FATAL] Failed to fix boot order", result);
            return FATAL;
        }
        result = IMachine_SetBootOrder(MACHINE_DREF(vbox_machine), 2, DeviceType_HardDisk);
        if (FAILED(result)) {
            print_error_info(error_message, "[FATAL] Failed to fix boot order", result);
            return FATAL;
        }
        result = IMachine_SetBootOrder(MACHINE_DREF(vbox_machine), 3, DeviceType_Null);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to fix boot order", result);
                ret = INFO;
            }
        }
        result = IMachine_SetBootOrder(MACHINE_DREF(vbox_machine), 4, DeviceType_Null);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to fix boot order", result);
                ret = INFO;
            }
        }
        
        // set Chipset type
        result = IMachine_SetChipsetType(MACHINE_DREF(vbox_machine), ChipsetType_ICH9);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to setting chipset type", result);
                ret = INFO;
            }
        }
        // set RTC timer
        result = IMachine_SetRTCUseUTC(MACHINE_DREF(vbox_machine), (PRBool)1);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to setting Hardware UTC timer", result);
                ret = INFO;
            }
        }
    }
    
    // Processor Setting
    {
        // set CPU Count
        result = IMachine_SetCPUCount(MACHINE_DREF(vbox_machine), cpu_count);
        if (FAILED(result)) {
            print_error_info(error_message, "[FATAL] Failed to setting CPU count", result);
            return FATAL;
        }
        // set Execution Cap
        result = IMachine_SetCPUExecutionCap(MACHINE_DREF(vbox_machine), 100);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to setting CPU execution cap", result);
                ret = INFO;
            }
        }
        // PAE enabled
        PRBool enabled = (PRBool)1;
        result = IMachine_GetCPUProperty(MACHINE_DREF(vbox_machine), CPUPropertyType_PAE, &enabled);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to setting PAE/NX enabling", result);
                ret = INFO;
            }
        }
    }
    
    // Acceleration
    {
        // Paravirtualization setting
        result = IMachine_SetParavirtProvider(MACHINE_DREF(vbox_machine), ParavirtProvider_Default);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to setting Pravirtualization", result);
                ret = INFO;
            }
        }
        // Nested Paging
        result = IMachine_SetHWVirtExProperty(MACHINE_DREF(vbox_machine), HWVirtExPropertyType_Enabled, (PRBool)1);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to setting HWVirtExPropertyType", result);
                ret = INFO;
            }
        }
    }
    
    // Display
    {
        // set VRAM
        result = VboxSetMachineVRAMSize(MACHINE_DREF(vbox_machine), 12);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to VRAM size", result);
                ret = INFO;
            }
        }
    }
    
    // SAVE SETTINGS & REGISTER MACHINE BEFORE ATTACH A MEDIUM
    {
        // save settings
        result = VboxMachineSaveSettings(MACHINE_DREF(vbox_machine));
        if (FAILED(result)) {
            print_error_info(error_message, "[FATAL] Failed to save machine before attaching a medium", result);
            return FATAL;
        }
        // Register machine
        result = VboxRegisterMachine(VBOX_DREF(virtualbox), MACHINE_DREF(vbox_machine));
        if (FAILED(result)) {
            print_error_info(error_message, "[FATAL] Failed to register machine", result);
            return FATAL;
        }
    }
    return ret;
}

VBRESULT
vbox_machine_add_bridged_network(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* host_interface, char* error_message) {

    INetworkAdapter *adapter = NULL;
    VBRESULT ret = GOOD;
    
    //firstly lock the machine
    HRESULT result = VboxLockMachine(MACHINE_DREF(vbox_machine), SESSION_DREF(vbox_session), LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[FATAL] Failed to lock machine for networking", result);
        return FATAL;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(SESSION_DREF(vbox_session), &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[FATAL] Failed to get a mutable copy of a machine for networking", result);
        ret = FATAL;
    }
    // get network adapter
    result = VboxMachineGetNetworkAdapter(mutable_machine, 0, &adapter);
    if (FAILED(result) || adapter == NULL) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to acquire adapter from slot 0", result);
            ret = FATAL;
        }
    }
    // enable network adapter
    result = VboxNetworkAdapterSetEnabled(adapter, TRUE);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to enable network adapter", result);
            ret = FATAL;
        }
    }
    // set bridged network type
    result = VboxNetworkAdapterSetAttachmentType(adapter, NetworkAttachmentType_Bridged);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to set network attachement type", result);
            ret = FATAL;
        }
    }

    // set host network adapter this bridge should connect to
    result = VboxNetworkAdapterSetBridgedHostInterface(adapter, host_interface);
    if (FAILED(result)) {
        print_error_info(error_message, "[FATAL] Failed to connect to host network interface", result);
        ret = FATAL;
    }    
    // set adapter type (AMD PCnet-FAST III, VBox Default)
    result = VboxNetworkAdapterSetAdapterType(adapter, NetworkAdapterType_Am79C973);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to set network adapter type", result);
            ret = FATAL;
        }
    }
    // promiscuous mode policy
    result = VboxNetworkAdapterSetPromiscModePolicy(adapter, NetworkAdapterPromiscModePolicy_Deny);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to set promiscuous mode", result);
            ret = INFO;
        }
    }
    // set cable connected
    result = VboxNetworkAdapterSetCableConnected(adapter, TRUE);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to set cable connected", result);
            ret = INFO;
        }
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to save machine after attaching hard disk medium", result);
            ret = FATAL;
        }
    }
    // release the first adapter
    result = VboxNetworkAdapterRelease(adapter);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to release adapter", result);
            ret = INFO;
        }
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to release locked machine for attaching adapter", result);
                ret = INFO;
            }
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(SESSION_DREF(vbox_session));
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to unlock machine for attaching adapter", result);
            ret = INFO;
        }
    }
    
    return ret;
}

VBRESULT
vbox_machine_add_shared_folder(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* shared_name, const char *host_folder, char* error_message) {

    VBRESULT ret = GOOD;
    HRESULT result;
    
    //firstly lock the machine
    result = VboxLockMachine(MACHINE_DREF(vbox_machine), SESSION_DREF(vbox_session), LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[FATAL] Failed to lock machine for shared folder", result);
        return FATAL;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(SESSION_DREF(vbox_session), &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[FATAL] Failed to get a mutable copy of a machine for shared folder", result);
        ret = FATAL;
    }
    // create shared folder
    result = VboxMachineCreateSharedFolder(mutable_machine, shared_name, host_folder, TRUE, TRUE);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to add shared folder", result);
            ret = FATAL;
        }
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to save machine for adding shared folder", result);
            ret = FATAL;
        }
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to release machine for adding shared folder", result);
                ret = FATAL;
            }
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(SESSION_DREF(vbox_session));
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to unlock machine for adding shared folder", result);
            ret = INFO;
        }
    }
    return ret;
}

VBRESULT
vbox_machine_add_storage_controller(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, char* error_message) {

    VBRESULT ret = GOOD;
    HRESULT result;
    IStorageController *storage_controller;
    
    //firstly lock the machine
    result = VboxLockMachine(MACHINE_DREF(vbox_machine), SESSION_DREF(vbox_session), LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[FATAL] Failed to lock machine for adding storage controller", result);
        return FATAL;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(SESSION_DREF(vbox_session), &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[FATAL] Failed to get a mutable copy of a machine", result);
        ret = FATAL;
    }
    // add storage controller
    result = VboxMachineAddStorageController(mutable_machine, storage_controller_name, StorageBus_SATA, &storage_controller);
    if (FAILED(result) || storage_controller_name == NULL) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to add storage controller", result);
            ret = FATAL;
        }
    }
    // storage controller type
    result = VboxSetStorageControllerType(storage_controller, StorageControllerType_IntelAhci);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to set storage controller type", result);
            ret = INFO;
        }
    }
    // storage controller set # of ports
    result = IStorageController_SetPortCount(storage_controller, 10);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to increase port count", result);
            ret = FATAL;
        }
    }
    // Enable host IO cache for imaging
    PRBool use_host_iocache = (PRBool)1;
    result = IStorageController_SetUseHostIOCache(storage_controller, use_host_iocache);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to enable host IO cache", result);
            ret = INFO;
        }
    }
    // release storage controller
    result = VboxIStorageControllerRelease(storage_controller);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to release storage controller for attaching storage controller", result);
            ret = FATAL;
        }
        
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to save machine after attaching storage controller", result);
            ret = FATAL;
        }
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to release mutable machine for attaching storage controller", result);
                ret = INFO;
            }
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(SESSION_DREF(vbox_session));
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to unlock machine for attaching storage controller", result);
            ret = INFO;
        }
    }

    return ret;
}

VBRESULT
vbox_machine_add_boot_image(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, const char *boot_image_path, char *error_message) {
    
    IMedium *boot_image;
    VBRESULT ret = GOOD;
    HRESULT result;

    // open medium
    PRBool use_old_uuid = (PRBool)0;
    result = VboxOpenMedium(VBOX_DREF(virtualbox), boot_image_path, DeviceType_DVD, AccessMode_ReadOnly, use_old_uuid, &boot_image);
    if (FAILED(result) || boot_image == NULL) {
        print_error_info(error_message, "[FATAL] Failed to open boot image for attaching boot image", result);
        return FATAL;
    }
    //firstly lock the machine
    result = VboxLockMachine(MACHINE_DREF(vbox_machine), SESSION_DREF(vbox_session), LockType_Write);
    if (FAILED(result)) {
        print_error_info(error_message, "[FATAL] Failed to lock machine for attaching boot image", result);
        return FATAL;
    }
    // get mutable machine
    IMachine *mutable_machine;
    result = VboxGetSessionMachine(SESSION_DREF(vbox_session), &mutable_machine);
    if (FAILED(result) || mutable_machine == NULL) {
        print_error_info(error_message, "[FATAL] Failed to get a mutable copy of a machine", result);
        ret = FATAL;
    }
    // attach a medium
    result = VboxMachineAttachDevice(mutable_machine, storage_controller_name, 0, 0, DeviceType_DVD, boot_image);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to attach boot image", result);
            ret = FATAL;
        }
    }
    // save setting
    result = VboxMachineSaveSettings(mutable_machine);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[FATAL] Failed to save machine after attaching boot image", result);
            ret = FATAL;
        }
    }
    // then we can safely release the mutable machine
    if (mutable_machine) {
        result = VboxIMachineRelease(mutable_machine);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to release machine after attaching boot image", result);
                ret = INFO;
            }
        }
    }
    // then unlock machine
    result = VboxUnlockMachine(SESSION_DREF(vbox_session));
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to unlock machine for attaching boot image", result);
            ret = INFO;
        }
    }
    // release medium
    result = VboxIMediumRelease(boot_image);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to release boot image", result);
            ret = INFO;
        }
    }
    return ret;
}

VBRESULT
vbox_machine_add_hard_disk(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, const char *hdd_medium_path, int disk_size, void(^build_progress)(int progress, int done), char *error_message) {

    // Create and Open hard drive
    HRESULT result;
    VBRESULT ret = GOOD;
    IMedium *hdd_medium;
    {
        // set medium image
        result = VboxCreateHardDisk(VBOX_DREF(virtualbox), "VMDK", hdd_medium_path, DeviceType_HardDisk, AccessMode_ReadWrite, &hdd_medium);
        if (FAILED(result) || hdd_medium == NULL) {
            print_error_info(error_message, "[FATAL] Failed to create harddrive", result);
            return FATAL;
        }
        // create medium
        //REF : https://www.virtualbox.org/sdkref/_virtual_box_8idl.html#adedcbf1a6e5e35fe7a0ca0c4b3447154
        PRUint32 cVariant[2] = {MediumVariant_Standard, MediumVariant_NoCreateDir};
        PRUint32 variantCount = sizeof(cVariant) / sizeof(cVariant[0]);
        IProgress *progress;
        result = VboxMediumCreateBaseStorage(hdd_medium, (1 << 24), variantCount, cVariant, &progress);
        if (FAILED(result)){
            print_error_info(error_message, "[FATAL] Failed to create base storage", result);
            return FATAL;
        }
        
        // it is recommended to wait short amount of time
        result = VboxProgressWaitForCompletion(progress, 3);
        if (FAILED(result)) {
            print_error_info(error_message, "[FATAL] Failed to complete creating base storage", result);
            ret = FATAL;
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
            print_error_info(error_message, "[FATAL] Failed to actuqire storage creation result code", result);
            return FATAL;
        }
        // release progress
        VboxIProgressRelease(progress);
        
        // open medium
        PRBool use_old_uuid = (PRBool)0;
        result = VboxOpenMedium(VBOX_DREF(virtualbox), hdd_medium_path, DeviceType_HardDisk, AccessMode_ReadWrite, use_old_uuid, &hdd_medium);
        if (FAILED(result)) {
            print_error_info(error_message, "[FATAL] Failed to open hard drive", result);
            return FATAL;
        }
    }
    
    // Attach a medium to storage controller
    {
        //firstly lock the machine
        result = VboxLockMachine(MACHINE_DREF(vbox_machine), SESSION_DREF(vbox_session), LockType_Write);
        if (FAILED(result)) {
            print_error_info(error_message, "[FATAL] Failed to lock machine for attaching hdd to storage controller", result);
            return FATAL;
        }
        // get mutable machine
        IMachine *mutable_machine;
        result = VboxGetSessionMachine(SESSION_DREF(vbox_session), &mutable_machine);
        if (FAILED(result) || mutable_machine == NULL) {
            if (ret == GOOD) {
                print_error_info(error_message, "[FATAL] Failed to get a mutable copy of a machine for attaching hard disk medium", result);
                ret = FATAL;
            }
        }
        // attach a medium
        result = VboxMachineAttachDevice(mutable_machine, storage_controller_name, 1, 0, DeviceType_HardDisk, hdd_medium);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[FATAL] Failed to attach hard disk medium", result);
                ret = FATAL;
            }
        }
        // save setting
        result = VboxMachineSaveSettings(mutable_machine);
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[FATAL] Failed to save machine after attaching hard disk medium", result);
                ret = FATAL;
            }
        }
        // then we can safely release the mutable machine
        if (mutable_machine) {
            result = VboxIMachineRelease(mutable_machine);
            if (FAILED(result)) {
                if (ret == GOOD) {
                    print_error_info(error_message, "[INFO] Failed to release machine after attaching hard disk medium", result);
                    ret = INFO;
                }
            }
        }
        // then unlock machine
        result = VboxUnlockMachine(SESSION_DREF(vbox_session));
        if (FAILED(result)) {
            if (ret == GOOD) {
                print_error_info(error_message, "[INFO] Failed to unlock machine after attaching hard disk medium", result);
                ret = INFO;
            }
        }
    }
    
    // Close & release hard drive
    result = VboxIMediumRelease(hdd_medium);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to release a hard drive", result);
            ret = INFO;
        }
    }
    
    // report the end of result
    if (build_progress != NULL) {
        build_progress(100, 1);
    }
    return ret;
}

#pragma mark - DESTROY MACHINE
VBRESULT
vbox_machine_destroy(VOID_DPTR vbox_machine, char* base_folder, const char* storage_controller_name, int remove_dvd, void(^build_progress)(int progress, int done), char *error_message) {
    
    HRESULT result;
    ULONG media_count;
    VBRESULT ret = GOOD;
    IProgress *progress;
    IMedium** media;
    
    // unregister
    result = VboxMachineUnregister(MACHINE_DREF(vbox_machine), (remove_dvd == 0 ? CleanupMode_DetachAllReturnHardDisksOnly:CleanupMode_Full) , &media, &media_count);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to unregister media", result);
            ret = INFO;
        }
    }
    // delete medium
    result = VboxMachineDeleteConfig(MACHINE_DREF(vbox_machine), media_count, media, &progress);
    if (FAILED(result)) {
        if (ret == GOOD) {
            print_error_info(error_message, "[INFO] Failed to delete medium", result);
            ret = INFO;
        }
    }
    // delete progress
    VboxProgressWaitForCompletion(progress, 3);
    PRUint32 progress_percent = 0;
    do {
        VboxGetProgressPercent(progress, &progress_percent);
        if (build_progress != NULL) {
            build_progress((int)progress_percent, 0);
        }
        usleep(500000);
    } while (progress_percent < 100);
    VboxIProgressRelease(progress);
    
    // free media array
    VboxArrayOutFree(media);

    // release machine
    result = VboxIMachineRelease(MACHINE_DREF(vbox_machine));
    if (FAILED(result)) {
        print_error_info(error_message, "[INFO] Failed to close machine referenece", result);
    }
    // release base folder
    if (base_folder != NULL) {
        VboxUtf8Free(base_folder);
    }
    // report the end of result
    if (build_progress != NULL) {
        build_progress(100, 1);
    }
    return ret;
}


#pragma mark - START & STOP MACHINE


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
