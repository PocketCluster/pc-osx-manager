#include "../include/common.h"
#include "../include/machine.h"

HRESULT VboxGetMachineName(IMachine* cmachine, char** cname) {
    BSTR wname = NULL;
    HRESULT result = IMachine_GetName(cmachine, &wname);
    if (FAILED(result)) {
        return result;
    }

    g_pVBoxFuncs->pfnUtf16ToUtf8(wname, cname);
    g_pVBoxFuncs->pfnComUnallocString(wname);
    return result;
}

HRESULT VboxGetMachineOSTypeId(IMachine* cmachine, char** cosTypeId) {
    BSTR wosTypeId = NULL;
    HRESULT result = IMachine_GetOSTypeId(cmachine, &wosTypeId);
    if (FAILED(result)) {
        return result;
    }

    g_pVBoxFuncs->pfnUtf16ToUtf8(wosTypeId, cosTypeId);
    g_pVBoxFuncs->pfnComUnallocString(wosTypeId);
    return result;

}

HRESULT VboxGetMachineSettingsFilePath(IMachine* cmachine, char** cpath) {
    BSTR wpath = NULL;
    HRESULT result = IMachine_GetSettingsFilePath(cmachine, &wpath);
    if (FAILED(result)) {
        return result;
    }

    g_pVBoxFuncs->pfnUtf16ToUtf8(wpath, cpath);
    g_pVBoxFuncs->pfnComUnallocString(wpath);
    return result;
}

HRESULT VboxGetMachineMemorySize(IMachine* cmachine, PRUint32* cram) {
    return IMachine_GetMemorySize(cmachine, cram);
}

HRESULT VboxSetMachineMemorySize(IMachine* cmachine, PRUint32 cram) {
    return IMachine_SetMemorySize(cmachine, cram);
}

HRESULT VboxGetMachineVRAMSize(IMachine* cmachine, PRUint32* cvram) {
    return IMachine_GetVRAMSize(cmachine, cvram);
}

HRESULT VboxSetMachineVRAMSize(IMachine* cmachine, PRUint32 cvram) {
    return IMachine_SetVRAMSize(cmachine, cvram);
}

HRESULT VboxGetMachinePointingHIDType(IMachine* cmachine, PRUint32* ctype) {
    return IMachine_GetPointingHIDType(cmachine, ctype);
}

HRESULT VboxSetMachinePointingHIDType(IMachine* cmachine, PRUint32 ctype) {
    return IMachine_SetPointingHIDType(cmachine, ctype);
}

HRESULT VboxGetMachineSettingsModified(IMachine* cmachine, PRBool* cmodified) {
    return IMachine_GetSettingsModified(cmachine, cmodified);
}

HRESULT VboxMachineSaveSettings(IMachine* cmachine) {
    return IMachine_SaveSettings(cmachine);
}

HRESULT VboxMachineUnregister(IMachine* cmachine, PRUint32 cleanupMode, IMedium*** cmedia, ULONG* mediaCount) {
    SAFEARRAY *safeArray = g_pVBoxFuncs->pfnSafeArrayOutParamAlloc();
    HRESULT result = IMachine_Unregister(cmachine, cleanupMode, ComSafeArrayAsOutIfaceParam(safeArray, IMedium *));
    if (!FAILED(result)) {
        result = g_pVBoxFuncs->pfnSafeArrayCopyOutIfaceParamHelper((IUnknown ***)cmedia, mediaCount, safeArray);
    }
    g_pVBoxFuncs->pfnSafeArrayDestroy(safeArray);
    return result;
}

HRESULT VboxMachineDeleteConfig(IMachine* cmachine, PRUint32 mediaCount, IMedium** cmedia, IProgress** cprogress) {
    SAFEARRAY *pSafeArray = g_pVBoxFuncs->pfnSafeArrayCreateVector(VT_UNKNOWN, 0, mediaCount);
    g_pVBoxFuncs->pfnSafeArrayCopyInParamHelper(pSafeArray, cmedia, sizeof(IMedium*) * mediaCount);
    HRESULT result = IMachine_DeleteConfig(cmachine, ComSafeArrayAsInParam(pSafeArray), cprogress);
    g_pVBoxFuncs->pfnSafeArrayDestroy(pSafeArray);
    return result;
}

HRESULT VboxMachineAttachDevice(IMachine* cmachine, const char* cname, PRInt32 cport, PRInt32 cdevice, PRUint32 ctype, IMedium* cmedium) {
    BSTR wname;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cname, &wname);
    if (FAILED(result)) {
        return result;
    }
    
    result = IMachine_AttachDevice(cmachine, wname, cport, cdevice, ctype, cmedium);
    g_pVBoxFuncs->pfnUtf16Free(wname);

    return result;
}

HRESULT VboxMachineUnmountMedium(IMachine* cmachine, const char* cname, PRInt32 cport, PRInt32 cdevice, PRBool cforce) {
    BSTR wname;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cname, &wname);
    if (FAILED(result)) {
        return result;
    }

    result = IMachine_UnmountMedium(cmachine, wname, cport, cdevice, cforce);
    g_pVBoxFuncs->pfnUtf16Free(wname);

    return result;
}

HRESULT VboxMachineGetMedium(IMachine* cmachine, char* cname, PRInt32 cport, PRInt32 cdevice, IMedium** cmedium) {
    BSTR wname;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cname, &wname);
    if (FAILED(result)) {
        return result;
    }

    result = IMachine_GetMedium(cmachine, wname, cport, cdevice, cmedium);
    g_pVBoxFuncs->pfnUtf16Free(wname);

    return result;
}

HRESULT VboxMachineCreateSharedFolder(IMachine *cmachine, const char* cName, const char* cHostPath, BOOL cWritable, BOOL cAutomount) {
    BSTR wName;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cName, &wName);
    if (FAILED(result)) {
        return result;
    }
    
    BSTR wHostPath;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cHostPath, &wHostPath);
    if (FAILED(result)) {
        return result;
    }
    
    PRBool writable  = (PRBool)cWritable;
    PRBool automount = (PRBool)cAutomount;
    
    result = IMachine_CreateSharedFolder(cmachine, wName, wHostPath, writable, automount);
    g_pVBoxFuncs->pfnUtf16Free(wName);
    g_pVBoxFuncs->pfnUtf16Free(wHostPath);
    return result;
}

HRESULT VboxMachineGetNetworkAdapter(IMachine *cmachine, PRUint32 slotNumber, INetworkAdapter **cAdapter) {
    return IMachine_GetNetworkAdapter(cmachine, slotNumber, cAdapter);
}

HRESULT VboxMachineLaunchVMProcess(IMachine* cmachine, ISession* csession, char* cuiType, char* cenvironment, IProgress** cprogress) {
    BSTR wuiType;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cuiType, &wuiType);
    if (FAILED(result)) {
        return result;
    }

    BSTR wenvironment;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cenvironment, &wenvironment);
    if (FAILED(result)) {
        g_pVBoxFuncs->pfnUtf16Free(wuiType);
        return result;
    }

    result = IMachine_LaunchVMProcess(cmachine, csession, wuiType, wenvironment, cprogress);
    g_pVBoxFuncs->pfnUtf16Free(wenvironment);
    g_pVBoxFuncs->pfnUtf16Free(wuiType);

    return result;
}

HRESULT VboxIMachineRelease(IMachine* cmachine) {
    HRESULT result = IMachine_Release(cmachine);
    cmachine = NULL;
    return result;
}


HRESULT VboxCreateMachine(IVirtualBox* cbox, char* cSettingsFile, const char* cname, char* cosTypeId, char* cflags, IMachine** cmachine) {
    BSTR wSettingsFile;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cSettingsFile, &wSettingsFile);
    if (FAILED(result)) {
        return result;
    }
    
    BSTR wname;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cname, &wname);
    if (FAILED(result)) {
        return result;
    }

    BSTR wosTypeId;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cosTypeId, &wosTypeId);
    if (FAILED(result)) {
        g_pVBoxFuncs->pfnUtf16Free(wname);
        return result;
    }

    BSTR wflags = NULL;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cflags, &wflags);
    if (FAILED(result)) {
        g_pVBoxFuncs->pfnUtf16Free(wosTypeId);
        g_pVBoxFuncs->pfnUtf16Free(wname);
    }

    SAFEARRAY *pSafeArray = g_pVBoxFuncs->pfnSafeArrayCreateVector(VT_BSTR, 0, 0);
    result = IVirtualBox_CreateMachine(cbox, wSettingsFile, wname, ComSafeArrayAsInParam(pSafeArray), wosTypeId, wflags, cmachine);
    g_pVBoxFuncs->pfnSafeArrayDestroy(pSafeArray);
    g_pVBoxFuncs->pfnUtf16Free(wflags);
    g_pVBoxFuncs->pfnUtf16Free(wosTypeId);
    g_pVBoxFuncs->pfnUtf16Free(wname);
    g_pVBoxFuncs->pfnUtf16Free(wSettingsFile);

    return result;
}

HRESULT VboxMachineGetID(IMachine* cmachine, char** cMachineId) {
    BSTR uuidUtf16;
    HRESULT result = IMachine_get_Id(cmachine, &uuidUtf16);
    g_pVBoxFuncs->pfnUtf16ToUtf8(uuidUtf16, cMachineId);
    g_pVBoxFuncs->pfnComUnallocString(uuidUtf16);
    return result;
}

HRESULT VboxFindMachine(IVirtualBox* cbox, const char* cnameOrId, IMachine** cmachine) {
    BSTR wnameOrId;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cnameOrId, &wnameOrId);
    if (FAILED(result)) {
        return result;
    }

    result = IVirtualBox_FindMachine(cbox, wnameOrId, cmachine);
    g_pVBoxFuncs->pfnUtf16Free(wnameOrId);

    return result;
}

HRESULT VboxGetMachines(IVirtualBox* cbox, IMachine*** cmachines, ULONG* machineCount) {
    SAFEARRAY *safeArray = g_pVBoxFuncs->pfnSafeArrayOutParamAlloc();
    HRESULT result = IVirtualBox_GetMachines(cbox, ComSafeArrayAsOutIfaceParam(safeArray, IMachine *));
    if (!FAILED(result)) {
        result = g_pVBoxFuncs->pfnSafeArrayCopyOutIfaceParamHelper((IUnknown ***)cmachines, machineCount, safeArray);
    }
    g_pVBoxFuncs->pfnSafeArrayDestroy(safeArray);
    return result;
}

HRESULT VboxRegisterMachine(IVirtualBox* cbox, IMachine* cmachine) {
    return IVirtualBox_RegisterMachine(cbox, cmachine);
}

