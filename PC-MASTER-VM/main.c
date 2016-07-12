//
//  main.c
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#if 0

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <signal.h>
#include <unistd.h>
#include <sys/poll.h>
#include "VBoxCAPIGlue.h"

/** Set by Ctrl+C handler. */
static volatile int g_fStop = 0;

static const char*
GetStateName(MachineState_T machineState) {
    
    switch (machineState) {
        case MachineState_Null:
            return "<null>";
            
        case MachineState_PoweredOff:
            return "PoweredOff";
            
        case MachineState_Saved:
            return "Saved";
            
        case MachineState_Teleported:
            return "Teleported";
            
        case MachineState_Aborted:
            return "Aborted";
            
        case MachineState_Running:
            return "Running";
            
        case MachineState_Paused:
            return "Paused";
            
        case MachineState_Stuck:
            return "Stuck";
            
        case MachineState_Teleporting:
            return "Teleporting";
            
        case MachineState_LiveSnapshotting:
            return "LiveSnapshotting";
            
        case MachineState_Starting:
            return "Starting";
            
        case MachineState_Stopping:
            return "Stopping";
            
        case MachineState_Saving:
            return "Saving";
            
        case MachineState_Restoring:
            return "Restoring";
            
        case MachineState_TeleportingPausedVM:
            return "TeleportingPausedVM";
            
        case MachineState_TeleportingIn:
            return "TeleportingIn";
            
        case MachineState_FaultTolerantSyncing:
            return "FaultTolerantSyncing";
            
        case MachineState_DeletingSnapshotOnline:
            return "DeletingSnapshotOnline";
            
        case MachineState_DeletingSnapshotPaused:
            return "DeletingSnapshotPaused";
            
        case MachineState_RestoringSnapshot:
            return "RestoringSnapshot";
            
        case MachineState_DeletingSnapshot:
            return "DeletingSnapshot";
            
        case MachineState_SettingUp:
            return "SettingUp";
            
        default:
            return "no idea";
    }
}

/**
 * Ctrl+C handler, terminate event listener.
 *
 * Remember that most function calls are not allowed in this context (including
 * printf!), so make sure that this does as little as possible.
 *
 * @param  iInfo    Platform dependent detail info (ignored).
 */
static BOOL VBOX_WINAPI
ctrlCHandler(DWORD iInfo) {
    (void)iInfo;
    g_fStop = 1;
    return TRUE;
}

/**
 * Sample event processing function, dumping some event information.
 * Shared between active and passive event demo, to highlight that this part
 * is identical between the two.
 */
static HRESULT
EventListenerDemoProcessEvent(IEvent *event) {

    VBoxEventType_T evType;
    HRESULT rc;
    
    if (!event) {
        printf("event null\n");
        return S_OK;
    }
    
    evType = VBoxEventType_Invalid;
    rc = IEvent_get_Type(event, &evType);
    if (FAILED(rc)) {
        printf("cannot get event type, rc=%#x\n", rc);
        return S_OK;
    }
    
    switch (evType) {
        case VBoxEventType_OnMousePointerShapeChanged:
            printf("OnMousePointerShapeChanged\n");
            break;
            
        case VBoxEventType_OnMouseCapabilityChanged:
            printf("OnMouseCapabilityChanged\n");
            break;
            
        case VBoxEventType_OnKeyboardLedsChanged:
            printf("OnMouseCapabilityChanged\n");
            break;
            
        case VBoxEventType_OnStateChanged: {
            IStateChangedEvent *ev = NULL;
            enum MachineState state;
            rc = IEvent_QueryInterface(event, &IID_IStateChangedEvent, (void **)&ev);
            if (FAILED(rc)) {
                printf("cannot get StateChangedEvent interface, rc=%#x\n", rc);
                return S_OK;
            }
            
            if (!ev) {
                printf("StateChangedEvent reference null\n");
                return S_OK;
            }

            rc = IStateChangedEvent_get_State(ev, &state);
            if (FAILED(rc)) {
                printf("warning: cannot get state, rc=%#x\n", rc);
            }
            IStateChangedEvent_Release(ev);
            printf("OnStateChanged: %s\n", GetStateName(state));
            
            fflush(stdout);
            if (state == MachineState_PoweredOff ||
                state == MachineState_Saved      ||
                state == MachineState_Teleported ||
                state == MachineState_Aborted) {
                g_fStop = 1;
            }
            break;
        }
            
        case VBoxEventType_OnAdditionsStateChanged:
            printf("OnAdditionsStateChanged\n");
            break;
            
        case VBoxEventType_OnNetworkAdapterChanged:
            printf("OnNetworkAdapterChanged\n");
            break;
            
        case VBoxEventType_OnSerialPortChanged:
            printf("OnSerialPortChanged\n");
            break;
            
        case VBoxEventType_OnParallelPortChanged:
            printf("OnParallelPortChanged\n");
            break;
            
        case VBoxEventType_OnStorageControllerChanged:
            printf("OnStorageControllerChanged\n");
            break;
            
        case VBoxEventType_OnMediumChanged:
            printf("OnMediumChanged\n");
            break;
            
        case VBoxEventType_OnVRDEServerChanged:
            printf("OnVRDEServerChanged\n");
            break;
            
        case VBoxEventType_OnUSBControllerChanged:
            printf("OnUSBControllerChanged\n");
            break;
            
        case VBoxEventType_OnUSBDeviceStateChanged:
            printf("OnUSBDeviceStateChanged\n");
            break;
            
        case VBoxEventType_OnSharedFolderChanged:
            printf("OnSharedFolderChanged\n");
            break;
            
        case VBoxEventType_OnRuntimeError:
            printf("OnRuntimeError\n");
            break;
            
        case VBoxEventType_OnCanShowWindow:
            printf("OnCanShowWindow\n");
            break;
        case VBoxEventType_OnShowWindow:
            printf("OnShowWindow\n");
            break;
            
        default:
            printf("unknown event: %d\n", evType);
    }
    
    return S_OK;
}

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
 * @param   pszExecutable   string with the executable name
 * @param   pszErrorMsg     string containing the code location specific error message
 * @param   rc              COM/XPCOM result code
 */
static void
PrintErrorInfo(const char *pszExecutable, const char *pszErrorMsg, HRESULT rc)
{
    IErrorInfo *ex;
    HRESULT rc2 = S_OK;
    
    fprintf(stderr, "%s: %s (rc=%#010x)\n", pszExecutable, pszErrorMsg, (unsigned)rc);
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
                fprintf(stderr, "Extended error info (IVirtualBoxErrorInfo):\n");
                
                IVirtualBoxErrorInfo_get_ResultCode(ei, &resultCode);
                fprintf(stderr, "  resultCode=%#010x\n", (unsigned)resultCode);
                
                IVirtualBoxErrorInfo_get_Component(ei, &componentUtf16);
                g_pVBoxFuncs->pfnUtf16ToUtf8(componentUtf16, &component);
                g_pVBoxFuncs->pfnComUnallocString(componentUtf16);
                fprintf(stderr, "  component=%s\n", component);
                g_pVBoxFuncs->pfnUtf8Free(component);
                
                IVirtualBoxErrorInfo_get_Text(ei, &textUtf16);
                g_pVBoxFuncs->pfnUtf16ToUtf8(textUtf16, &text);
                g_pVBoxFuncs->pfnComUnallocString(textUtf16);
                fprintf(stderr, "  text=%s\n", text);
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

/**
 * Start a VM.
 *
 * @param   argv0       executable name
 * @param   virtualBox  ptr to IVirtualBox object
 * @param   session     ptr to ISession object
 * @param   id          identifies the machine to start
 */
static void
startVM(const char *argv0, IVirtualBox *virtualBox, ISession *session, BSTR id)
{
    HRESULT rc;
    IMachine  *machine    = NULL;
    IProgress *progress   = NULL;
    BSTR env              = NULL;
    BSTR sessionType;
    SAFEARRAY *groupsSA = g_pVBoxFuncs->pfnSafeArrayOutParamAlloc();
    
    rc = IVirtualBox_FindMachine(virtualBox, id, &machine);
    if (FAILED(rc) || !machine) {
        PrintErrorInfo(argv0, "Error: Couldn't get the Machine reference", rc);
        return;
    }
    
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
            printf("Groups[%d]: %s\n", i, group);
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
#ifdef USE_ACTIVE_EVENT_LISTENER
            registerActiveEventListener(virtualBox, session, id);
#else /* !USE_ACTIVE_EVENT_LISTENER */
            registerPassiveEventListener(virtualBox, session, id);
#endif /* !USE_ACTIVE_EVENT_LISTENER */
        }
        IProgress_Release(progress);
    } else {
        PrintErrorInfo(argv0, "Error: LaunchVMProcess failed", rc);
    }

    /* It's important to always release resources. */
    IMachine_Release(machine);
}

/**
 * List the registered VMs.
 *
 * @param   argv0       executable name
 * @param   virtualBox  ptr to IVirtualBox object
 * @param   session     ptr to ISession object
 */
static void
listVMs(const char *argv0, IVirtualBox *virtualBox, ISession *session) {
    
    HRESULT rc;
    SAFEARRAY *machinesSA = g_pVBoxFuncs->pfnSafeArrayOutParamAlloc();
    IMachine **machines = NULL;
    ULONG machineCnt = 0;
    ULONG i;
    unsigned start_id;
    
    /*
     * Get the list of all registered VMs.
     */
    rc = IVirtualBox_get_Machines(virtualBox, ComSafeArrayAsOutIfaceParam(machinesSA, IMachine *));
    if (FAILED(rc)) {
        PrintErrorInfo(argv0, "could not get list of machines", rc);
        return;
    }
    
    /*
     * Extract interface pointers from machinesSA, and update the reference
     * counter of each object, as destroying machinesSA would call Release.
     */
    g_pVBoxFuncs->pfnSafeArrayCopyOutIfaceParamHelper((IUnknown ***)&machines, &machineCnt, machinesSA);
    g_pVBoxFuncs->pfnSafeArrayDestroy(machinesSA);
    
    if (!machineCnt) {
        g_pVBoxFuncs->pfnArrayOutFree(machines);
        printf("\tNo VMs\n");
        return;
    }
    
    printf("VM List:\n\n");
    
    /*
     * Iterate through the collection.
     */
    
    for (i = 0; i < machineCnt; ++i) {
        IMachine *machine      = machines[i];
        BOOL      isAccessible = FALSE;
        
        printf("\tMachine #%u\n", (unsigned)i);
        
        if (!machine) {
            printf("\t(skipped, NULL)\n");
            continue;
        }
        
        IMachine_get_Accessible(machine, &isAccessible);
        
        if (isAccessible) {
            BSTR machineNameUtf16;
            char *machineName;
            
            IMachine_get_Name(machine, &machineNameUtf16);
            g_pVBoxFuncs->pfnUtf16ToUtf8(machineNameUtf16,&machineName);
            g_pVBoxFuncs->pfnComUnallocString(machineNameUtf16);
            printf("\tName:        %s\n", machineName);
            g_pVBoxFuncs->pfnUtf8Free(machineName);
        } else {
            printf("\tName:        <inaccessible>\n");
        }
        
        {
            BSTR uuidUtf16;
            char      *uuidUtf8;
            
            IMachine_get_Id(machine, &uuidUtf16);
            g_pVBoxFuncs->pfnUtf16ToUtf8(uuidUtf16, &uuidUtf8);
            g_pVBoxFuncs->pfnComUnallocString(uuidUtf16);
            printf("\tUUID:        %s\n", uuidUtf8);
            g_pVBoxFuncs->pfnUtf8Free(uuidUtf8);
        }
        
        if (isAccessible)
        {
            {
                BSTR      configFileUtf16;
                char      *configFileUtf8;
                
                IMachine_get_SettingsFilePath(machine, &configFileUtf16);
                g_pVBoxFuncs->pfnUtf16ToUtf8(configFileUtf16, &configFileUtf8);
                g_pVBoxFuncs->pfnComUnallocString(configFileUtf16);
                printf("\tConfig file: %s\n", configFileUtf8);
                g_pVBoxFuncs->pfnUtf8Free(configFileUtf8);
            }
            
            {
                ULONG memorySize;
                
                IMachine_get_MemorySize(machine, &memorySize);
                printf("\tMemory size: %uMB\n", memorySize);
            }
            
            {
                BSTR typeId;
                BSTR osNameUtf16;
                char *osName;
                IGuestOSType *osType = NULL;
                
                IMachine_get_OSTypeId(machine, &typeId);
                IVirtualBox_GetGuestOSType(virtualBox, typeId, &osType);
                g_pVBoxFuncs->pfnComUnallocString(typeId);
                IGuestOSType_get_Description(osType, &osNameUtf16);
                g_pVBoxFuncs->pfnUtf16ToUtf8(osNameUtf16,&osName);
                g_pVBoxFuncs->pfnComUnallocString(osNameUtf16);
                printf("\tGuest OS:    %s\n\n", osName);
                g_pVBoxFuncs->pfnUtf8Free(osName);
                
                IGuestOSType_Release(osType);
            }
        }
    }
    
    /*
     * Let the user chose a machine to start.
     */
    
    printf("Type Machine# to start (0 - %u) or 'quit' to do nothing: ", (unsigned)(machineCnt - 1));
    
    fflush(stdout);
    
    if (scanf("%u", &start_id) == 1 && start_id < machineCnt) {
        IMachine *machine = machines[start_id];
        
        if (machine) {
            
            BSTR uuidUtf16 = NULL;
            IMachine_get_Id(machine, &uuidUtf16);
            startVM(argv0, virtualBox, session, uuidUtf16);
            g_pVBoxFuncs->pfnComUnallocString(uuidUtf16);
        }
    }
    
    /*
     * Don't forget to release the objects in the array.
     */
    
    for (i = 0; i < machineCnt; ++i) {
        IMachine *machine = machines[i];
        
        if (machine) {
            IMachine_Release(machine);
        }
    }
    g_pVBoxFuncs->pfnArrayOutFree(machines);
}


int
main(int argc, const char * argv[]) {

    IVirtualBoxClient *vboxclient = NULL;
    IVirtualBox *vbox            = NULL;
    ISession   *session          = NULL;
    ULONG       revision         = 0;
    BSTR        versionUtf16     = NULL;
    BSTR        homefolderUtf16  = NULL;
    HRESULT    rc;     /* Result code of various function (method) calls. */
    
    printf("Starting main()\n");
    
    if (VBoxCGlueInit())
    {
        fprintf(stderr, "%s: FATAL: VBoxCGlueInit failed: %s\n",
                argv[0], g_szVBoxErrMsg);
        return EXIT_FAILURE;
    }
    
    {
        unsigned ver = g_pVBoxFuncs->pfnGetVersion();
        printf("VirtualBox version: %u.%u.%u\n", ver / 1000000, ver / 1000 % 1000, ver % 1000);
        ver = g_pVBoxFuncs->pfnGetAPIVersion();
        printf("VirtualBox API version: %u.%u\n", ver / 1000, ver % 1000);
    }
    
    g_pVBoxFuncs->pfnClientInitialize(NULL, &vboxclient);
    if (!vboxclient)
    {
        fprintf(stderr, "%s: FATAL: could not get VirtualBoxClient reference\n", argv[0]);
        return EXIT_FAILURE;
    }
    
    printf("----------------------------------------------------\n");
    
    rc = IVirtualBoxClient_get_VirtualBox(vboxclient, &vbox);
    if (FAILED(rc) || !vbox)
    {
        PrintErrorInfo(argv[0], "FATAL: could not get VirtualBox reference", rc);
        return EXIT_FAILURE;
    }
    rc = IVirtualBoxClient_get_Session(vboxclient, &session);
    if (FAILED(rc) || !session)
    {
        PrintErrorInfo(argv[0], "FATAL: could not get Session reference", rc);
        return EXIT_FAILURE;
    }
    
    /*
     * Now ask for revision, version and home folder information of
     * this vbox. Were not using fancy macros here so it
     * remains easy to see how we access C++'s vtable.
     */
    
    /* 1. Revision */
    
    rc = IVirtualBox_get_Revision(vbox, &revision);
    if (SUCCEEDED(rc))
        printf("\tRevision: %u\n", revision);
    else
        PrintErrorInfo(argv[0], "GetRevision() failed", rc);
    
    /* 2. Version */
    
    rc = IVirtualBox_get_Version(vbox, &versionUtf16);
    if (SUCCEEDED(rc))
    {
        char *version = NULL;
        g_pVBoxFuncs->pfnUtf16ToUtf8(versionUtf16, &version);
        printf("\tVersion: %s\n", version);
        g_pVBoxFuncs->pfnUtf8Free(version);
        g_pVBoxFuncs->pfnComUnallocString(versionUtf16);
    }
    else
        PrintErrorInfo(argv[0], "GetVersion() failed", rc);
    
    /* 3. Home Folder */
    
    rc = IVirtualBox_get_HomeFolder(vbox, &homefolderUtf16);
    if (SUCCEEDED(rc))
    {
        char *homefolder = NULL;
        g_pVBoxFuncs->pfnUtf16ToUtf8(homefolderUtf16, &homefolder);
        printf("\tHomeFolder: %s\n", homefolder);
        g_pVBoxFuncs->pfnUtf8Free(homefolder);
        g_pVBoxFuncs->pfnComUnallocString(homefolderUtf16);
    }
    else
        PrintErrorInfo(argv[0], "GetHomeFolder() failed", rc);
    
    listVMs(argv[0], vbox, session);
    ISession_UnlockMachine(session);
    
    printf("----------------------------------------------------\n");
    
    /*
     * Do as mom told us: always clean up after yourself.
     */
    if (session)
    {
        ISession_Release(session);
        session = NULL;
    }
    if (vbox)
    {
        IVirtualBox_Release(vbox);
        vbox = NULL;
    }
    if (vboxclient)
    {
        IVirtualBoxClient_Release(vboxclient);
        vboxclient = NULL;
    }
    
    g_pVBoxFuncs->pfnClientUninitialize();
    VBoxCGlueTerm();
    printf("Finished main()\n");
    
    return 0;
}
#endif


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

/**
 * Print detailed error information if available.
 * @param   pszErrorMsg     string containing the code location specific error message
 * @param   rc              COM/XPCOM result code
 */
static void
print_error_info(const char *pszErrorMsg, HRESULT rc)
{
    IErrorInfo *ex;
    HRESULT rc2 = S_OK;

    fprintf(stderr, "\n--- %s (rc=%#010x) ---\n", pszErrorMsg, (unsigned)rc);
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
                fprintf(stderr, "Extended error info (IVirtualBoxErrorInfo):\n");
                
                IVirtualBoxErrorInfo_get_ResultCode(ei, &resultCode);
                fprintf(stderr, "  resultCode=%#010x\n", (unsigned)resultCode);
                
                IVirtualBoxErrorInfo_get_Component(ei, &componentUtf16);
                g_pVBoxFuncs->pfnUtf16ToUtf8(componentUtf16, &component);
                g_pVBoxFuncs->pfnComUnallocString(componentUtf16);
                fprintf(stderr, "  component=%s\n", component);
                g_pVBoxFuncs->pfnUtf8Free(component);
                
                IVirtualBoxErrorInfo_get_Text(ei, &textUtf16);
                g_pVBoxFuncs->pfnUtf16ToUtf8(textUtf16, &text);
                g_pVBoxFuncs->pfnComUnallocString(textUtf16);
                fprintf(stderr, "  text=%s\n", text);
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
char *str_replace(char *orig, char *rep, char *with) {
    char *result;            // the return string
    char *ins;               // the next insert point
    char *tmp;               // varies
    size_t len_rep;   // length of rep
    size_t len_with;  // length of with
    size_t len_front; // distance between rep and end of last rep
    size_t count;     // number of replacements
    
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


#define MAKE_HDD

int
main(int argc, const char * argv[]) {
    /* Result code of various function (method) calls. */
    HRESULT           result;
    
    IVirtualBoxClient* vboxclient     = NULL;
    IVirtualBox*       vbox           = NULL;
    ISession*          session        = NULL;
    IMachine*          machine        = NULL;
    char*              base_folder    = NULL;
    char*              machine_name   = "pc-master";
    
    
    result = VboxInit();
    if (FAILED(result)) {
        return result;
    }
    VboxClientInitialize(&vboxclient);
    if (!vboxclient) {
        print_error_info("FATAL: could not get VirtualBoxClient reference\n", result);
        return EXIT_FAILURE;
    }

    result = VboxGetVirtualBox(vboxclient, &vbox);
    if (FAILED(result) || !vbox) {
        print_error_info("FATAL: could not get VirtualBox reference %d", result);
        return EXIT_FAILURE;
    }
    result = VboxGetSession(vboxclient, &session);
    if (FAILED(result) || !session) {
        print_error_info("FATAL: could not get Session reference", result);
        return EXIT_FAILURE;
    }
    
    // create machine file name
    result = VboxComposeMachineFilename(vbox, machine_name, "", "", &base_folder);
    if (FAILED(result)) {
        print_error_info("Failed composing machine name", result);
    }
    
    // create machine based on the
    result = VboxCreateMachine(vbox, base_folder, machine_name, "Linux26_64", "", &machine);
    if (FAILED(result) || machine == NULL) {
        print_error_info("Failed to create machine", result);
    }

#pragma mark - VBOX CLIENT READY
    printf("------------------------VBOX CLIENT READY----------------------------\n");
    printf("Base folder : %s\n", base_folder);
    
    // Setup BIOS
    {
        // get BIOS settings
        IBIOSSettings *bios;
        result = VboxGetMachineBIOSSettings(machine, &bios);
        if (FAILED(result) || bios == NULL) {
            print_error_info("Failed to acquire bios settings", result);
        }
        // enable I/O APIC
        result = IBIOSSettings_SetIOAPICEnabled(bios, (PRBool)1);
        if (FAILED(result)) {
            print_error_info("Failed to enable IO APIC", result);
        }
        // set ACPI enabled
        result = IBIOSSettings_SetACPIEnabled(bios, (PRBool)1);
        if (FAILED(result)) {
            print_error_info("Failed to enable IO ACPI", result);
        }
        // release bios settings
        VboxIBiosSettingsRelease(bios);
    }

    // Motherboard Settings
    {
        // set memory
        result = VboxSetMachineMemorySize(machine, 1024);
        if (FAILED(result)) {
            print_error_info("Failed to set memory size", result);
        }
        
        // set up Boot Order
        IMachine_SetBootOrder(machine, 1, DeviceType_DVD);
        IMachine_SetBootOrder(machine, 2, DeviceType_HardDisk);
        IMachine_SetBootOrder(machine, 3, DeviceType_Null);
        IMachine_SetBootOrder(machine, 4, DeviceType_Null);
        
        // set Chipset type
        result = IMachine_SetChipsetType(machine, ChipsetType_ICH9);
        if (FAILED(result)) {
            print_error_info("Failed to setting chipset type", result);
        }

        // set RTC timer
        result = IMachine_SetRTCUseUTC(machine, (PRBool)1);
        if (FAILED(result)) {
            print_error_info("Failed to setting Hardware UTC timer", result);
        }
    }
    
    // Processor Setting
    {
        // set CPU Count
        result = IMachine_SetCPUCount(machine, 4);
        if (FAILED(result)) {
            print_error_info("Failed to setting CPU count", result);
        }
        
        // set Execution Cap
        result = IMachine_SetCPUExecutionCap(machine, 100);
        if (FAILED(result)) {
            print_error_info("Failed to setting CPU execution cap", result);
        }
        
        // PAE enabled
        PRBool enabled = (PRBool)1;
        result = IMachine_GetCPUProperty(machine, CPUPropertyType_PAE, &enabled);
        if (FAILED(result)) {
            print_error_info("Failed to setting PAE/NX enabling", result);
        }
    }
    
    // Acceleration
    {
        // Paravirtualization setting
        result = IMachine_SetParavirtProvider(machine, ParavirtProvider_Default);
        if (FAILED(result)) {
            print_error_info("Failed to setting Pravirtualization", result);
        }
        
        // Nested Paging
        result = IMachine_SetHWVirtExProperty(machine, HWVirtExPropertyType_Enabled, (PRBool)1);
        if (FAILED(result)) {
            print_error_info("Failed to setting HWVirtExPropertyType", result);
        }
    }
    
    // Display
    {
        // set VRAM
        result = VboxSetMachineVRAMSize(machine, 12);
        if (FAILED(result)) {
            print_error_info("Failed to VRAM size", result);
        }
    }
    
    // STORAGE
    {
        // STORAGE CONTROLLER
        char* storage_controller_name;
        {
            IStorageController *storage_controller;
            storage_controller_name = "SATA";
            result = VboxMachineAddStorageController(machine, storage_controller_name, StorageBus_SATA, &storage_controller);
            if (FAILED(result) || storage_controller_name == NULL) {
                print_error_info("Failed to add storage controller", result);
            }
            // storage controller type
            result = VboxSetStorageControllerType(storage_controller, StorageControllerType_IntelAhci);
            if (FAILED(result)) {
                print_error_info("Failed to set storage controller type", result);
            }

            // storage controller set # of ports
            result = IStorageController_SetPortCount(storage_controller, 10);
            if (FAILED(result)) {
                print_error_info("Failed to increase port count", result);
            }
            // Enable host IO cache for imaging
            PRBool use_host_iocache = (PRBool)1;
            result = IStorageController_SetUseHostIOCache(storage_controller, use_host_iocache);
            if (FAILED(result)) {
                print_error_info("Failed to enable host IO cache", result);
            }

            // release storage controller
            VboxIStorageControllerRelease(storage_controller);
        }
        
        // SAVE SETTINGS & REGISTER MACHINE BEFORE ATTACH A MEDIUM
        {
            result = VboxMachineSaveSettings(machine);
            if (FAILED(result)) {
                print_error_info("Failed to save machine before attaching a medium", result);
            }
            // Register machine
            result = VboxRegisterMachine(vbox, machine);
            if (FAILED(result)) {
                print_error_info("Failed to register machine", result);
            }
        }

#pragma mark - ATTACH BOOT IMAGE
        // CREATE & ATTACHE BOOT IMAGE
        {
            // Create and Open Boot Image
            IMedium *boot_image;
            char *boot_image_path = "/Users/almightykim/Downloads/boot2docker.iso";
            
            printf("BOOT IMAGE PATH : %s\n", boot_image_path);
            {
                // open medium
                PRBool use_old_uuid = (PRBool)0;
                result = VboxOpenMedium(vbox, boot_image_path, DeviceType_DVD, AccessMode_ReadOnly, use_old_uuid, &boot_image);
                if (FAILED(result) || boot_image == NULL) {
                    print_error_info("Failed to open boot image", result);
                }
            }
            
            // Attach a medium to storage controller
            {
                //firstly lock the machine
                result = VboxLockMachine(machine, session, LockType_Write);
                if (FAILED(result)) {
                    print_error_info("Failed to lock machine", result);
                }
                
                // get mutable machine
                IMachine *mutable_machine;
                result = VboxGetSessionMachine(session, &mutable_machine);
                if (FAILED(result) || mutable_machine == NULL) {
                    print_error_info("Failed to get a mutable copy of a machine", result);
                }
                // attach a medium
                result = VboxMachineAttachDevice(mutable_machine, storage_controller_name, 0, 0, DeviceType_DVD, boot_image);
                if (FAILED(result)) {
                    print_error_info("Failed to attache boot image", result);
                }
                // save setting
                result = VboxMachineSaveSettings(mutable_machine);
                if (FAILED(result)) {
                    print_error_info("Failed to save machine after attaching boot image", result);
                }
                // then we can safely release the mutable machine
                if (mutable_machine) {
                    VboxIMachineRelease(mutable_machine);
                    mutable_machine = NULL;
                }
                // then unlock machine
                result = VboxUnlockMachine(session);
                if (FAILED(result)) {
                    print_error_info("Failed to unlock machine", result);
                }
            }

            // Close & release boot image
            {
                // release medium
                result = VboxIMediumRelease(boot_image);
                if (FAILED(result)) {
                    print_error_info("Failed to release a hard drive", result);
                }
/*
                if (boot_image_path != NULL) {
                    free(boot_image_path);
                    boot_image_path = NULL;
                }
*/
            }
        }

#pragma mark - CREATE AND ATTACH HDD
        // CREATE & ATTACHE A HDD MEDIUM
        {
            // Create and Open hard drive
            IMedium *hdd_medium;
            char *hdd_medium_path = NULL;
            {
                PRBool use_old_uuid = (PRBool)0;
                hdd_medium_path = str_replace(base_folder, "pc-master.vbox", "pc-master-hdd.vmdk");
                printf("HDD path : %s\n", hdd_medium_path);
                
                result = VboxCreateHardDisk(vbox, "VMDK", hdd_medium_path, DeviceType_HardDisk, AccessMode_ReadWrite, &hdd_medium);
                if (FAILED(result) || hdd_medium == NULL) {
                    print_error_info("Failed to create harddrive", result);
                }
                
//TODO : make sure we put enough time buffer + progress monitor for hard disk to be created
                
                PRUint32 cVariant[2] = {MediumVariant_Standard, MediumVariant_NoCreateDir};
                PRUint32 variantCount = sizeof(cVariant) / sizeof(cVariant[0]);
                
                IProgress *progress;
                result = VboxMediumCreateBaseStorage(hdd_medium, (1 << 24), variantCount, cVariant, &progress);
                if (FAILED(result)){
                    print_error_info("Failed to create base storage", result);
                }
                result = VboxProgressWaitForCompletion(progress, 100000);
                if (FAILED(result)) {
                    print_error_info("Failed to complete creating base storage", result);
                }
                
                PRInt32 code;
                result = VboxGetProgressResultCode(progress, &code);
                if (FAILED(result)|| code != 0) {
                    print_error_info("Actuqired bad storage creation result code", result);
                }
                
                // open medium
                result = VboxOpenMedium(vbox, hdd_medium_path, DeviceType_HardDisk, AccessMode_ReadWrite, use_old_uuid, &hdd_medium);
                if (FAILED(result)) {
                    print_error_info("Failed to open hard drive", result);
                }
            }
            
            // Attach a medium to storage controller
            {
                
                //firstly lock the machine
                result = VboxLockMachine(machine, session, LockType_Write);
                if (FAILED(result)) {
                    print_error_info("Failed to lock machine", result);
                }
                // get mutable machine
                IMachine *mutable_machine;
                result = VboxGetSessionMachine(session, &mutable_machine);
                if (FAILED(result) || mutable_machine == NULL) {
                    print_error_info("Failed to get a mutable copy of a machine", result);
                }
                // attach a medium
                result = VboxMachineAttachDevice(mutable_machine, storage_controller_name, 1, 0, DeviceType_HardDisk, hdd_medium);
                if (FAILED(result)) {
                    print_error_info("Failed to attach hard disk medium", result);
                }
                // save setting
                result = VboxMachineSaveSettings(mutable_machine);
                if (FAILED(result)) {
                    print_error_info("Failed to save machine after attaching hard disk medium", result);
                }
                // then we can safely release the mutable machine
                if (mutable_machine) {
                    VboxIMachineRelease(mutable_machine);
                    mutable_machine = NULL;
                }
                // then unlock machine
                result = VboxUnlockMachine(session);
                if (FAILED(result)) {
                    print_error_info("Failed to unlock machine", result);
                }
            }
            
            // Close & release hard drive
            {
/*
                // medium umount
                PRBool force_umount = (PRBool)1;
                result = VboxMachineUnmountMedium(machine, storage_controller_name, 0, 0, force_umount);
                if (FAILED(result)) {
                    print_error_info("Failed to unmount medium", result);
                }
                // close medium
                result = VboxMediumClose(hdd_medium);
                if (FAILED(result)) {
                    print_error_info("Failed to close medium", result);
                }
*/
                // release medium
                result = VboxIMediumRelease(hdd_medium);
                if (FAILED(result)) {
                    print_error_info("Failed to release a hard drive", result);
                }
                
                if (hdd_medium_path != NULL) {
                    free(hdd_medium_path);
                    hdd_medium_path = NULL;
                }
            }
        }
    }

#pragma mark - ADD NETWORK ADAPTER
    // BRIDGED FIRST NETWORK
    {
        INetworkAdapter *adapter = NULL;
        {
            //firstly lock the machine
            result = VboxLockMachine(machine, session, LockType_Write);
            if (FAILED(result)) {
                print_error_info("Failed to lock machine for networking", result);
            }
            // get mutable machine
            IMachine *mutable_machine;
            result = VboxGetSessionMachine(session, &mutable_machine);
            if (FAILED(result) || mutable_machine == NULL) {
                print_error_info("Failed to get a mutable copy of a machine for networking", result);
            }
            // get network adapter
            result = VboxMachineGetNetworkAdapter(mutable_machine, 0, &adapter);
            if (FAILED(result) || adapter == NULL) {
                print_error_info("Failed to acquire adapter from slot 0", result);
            }
            // enable network adapter
            result = VboxNetworkAdapterSetEnabled(adapter, TRUE);
            if (FAILED(result)) {
                print_error_info("Failed to enable network adapter", result);
            }
            // set bridged network type
            result = VboxNetworkAdapterSetAttachmentType(adapter, NetworkAttachmentType_Bridged);
            if (FAILED(result)) {
                print_error_info("Failed to set network attachement type", result);
            }
/*
            // set host network adapter this bridge should connect to
            result = VboxNetworkAdapterSetBridgedHostInterface(adapter, NULL);
            if (FAILED(result)) {
                print_error_info("Failed to connect to host network interface", result);
            }
*/
            // set adapter type (AMD PCnet-FAST III, VBox Default)
            result = VboxNetworkAdapterSetAdapterType(adapter, NetworkAdapterType_Am79C973);
            if (FAILED(result)) {
                print_error_info("Failed to set network adapter type", result);
            }
            // promiscuous mode policy
            result = VboxNetworkAdapterSetPromiscModePolicy(adapter, NetworkAdapterPromiscModePolicy_Deny);
            if (FAILED(result)) {
                print_error_info("Failed to set promiscuous mode", result);
            }
            // set cable connected
            result = VboxNetworkAdapterSetCableConnected(adapter, TRUE);
            if (FAILED(result)) {
                print_error_info("Failed to set cable connected", result);
            }
            // save setting
            result = VboxMachineSaveSettings(mutable_machine);
            if (FAILED(result)) {
                print_error_info("Failed to save machine after attaching hard disk medium", result);
            }
            // release the first adapter
            VboxNetworkAdapterRelease(adapter);
            // then we can safely release the mutable machine
            if (mutable_machine) {
                VboxIMachineRelease(mutable_machine);
            }
            // then unlock machine
            result = VboxUnlockMachine(session);
            if (FAILED(result)) {
                print_error_info("Failed to unlock machine", result);
            }
        }
    }
    
#pragma mark - ADD SHARED FOLDER
    // ADD SHARED FOLDER
    {
        //firstly lock the machine
        result = VboxLockMachine(machine, session, LockType_Write);
        if (FAILED(result)) {
            print_error_info("Failed to lock machine for shared folder", result);
        }
        // get mutable machine
        IMachine *mutable_machine;
        result = VboxGetSessionMachine(session, &mutable_machine);
        if (FAILED(result) || mutable_machine == NULL) {
            print_error_info("Failed to get a mutable copy of a machine for shared folder", result);
        }
        // create shared folder
        result = VboxMachineCreateSharedFolder(mutable_machine, "/pocket", "/Users/almightykim/Workspace", TRUE, TRUE);
        if (FAILED(result)) {
            print_error_info("Failed to add shared folder /pocket", result);
        }
        // save setting
        result = VboxMachineSaveSettings(mutable_machine);
        if (FAILED(result)) {
            print_error_info("Failed to save machine after attaching hard disk medium", result);
        }
        // then we can safely release the mutable machine
        if (mutable_machine) {
            VboxIMachineRelease(mutable_machine);
        }
        // then unlock machine
        result = VboxUnlockMachine(session);
        if (FAILED(result)) {
            print_error_info("Failed to unlock machine", result);
        }
        
    }
    printf("------------------------VBOX CLIENT TERMINATED ----------------------------\n");
#pragma mark - VBOX CLIENT TERMINATED
    if (machine) {
        VboxIMachineRelease(machine);
    }
    if (base_folder) {
        VboxUtf8Free(base_folder);
    }
    if (session) {
        VboxISessionRelease(session);
    }
    if (vbox) {
        VboxIVirtualBoxRelease(vbox);
    }
    if (vboxclient) {
        VboxClientRelease(vboxclient);
    }
    
    VboxClientUninitialize();
    VboxTerm();
    
    printf("Finished main()\n");
}
