#include "../include/common.h"
#include "../include/medium.h"

HRESULT VboxCreateHardDisk(IVirtualBox* cbox, char* cformat, char* clocation, DeviceType deviceType, AccessMode accessMode, IMedium** cmedium) {
    BSTR wformat;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cformat, &wformat);
    if (FAILED(result)) {
        return result;
    }

    BSTR wlocation;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(clocation, &wlocation);
    if (FAILED(result)) {
        g_pVBoxFuncs->pfnUtf16Free(wformat);
        return result;
    }

    //result = IVirtualBox_CreateHardDisk(cbox, wformat, wlocation, cmedium);
    result = IVirtualBox_CreateMedium(cbox, wformat, wlocation, accessMode, deviceType, cmedium);
    g_pVBoxFuncs->pfnUtf16Free(wlocation);
    g_pVBoxFuncs->pfnUtf16Free(wformat);

    return result;
}

HRESULT VboxOpenMedium(IVirtualBox* cbox, char* clocation, DeviceType cdeviceType, AccessMode caccessType, PRBool cforceNewUuid, IMedium** cmedium) {
    BSTR wlocation;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(clocation, &wlocation);
    if (FAILED(result)) {
        return result;
    }

    result = IVirtualBox_OpenMedium(cbox, wlocation, cdeviceType, caccessType, cforceNewUuid, cmedium);
    g_pVBoxFuncs->pfnUtf16Free(wlocation);

    return result;
}

HRESULT VboxMediumCreateBaseStorage(IMedium* cmedium, PRInt64 size, PRUint32 variantCount, PRUint32* cvariant, IProgress** cprogress) {

    SAFEARRAY *pSafeArray = g_pVBoxFuncs->pfnSafeArrayCreateVector(VT_UI4, 0, variantCount);
    g_pVBoxFuncs->pfnSafeArrayCopyInParamHelper(pSafeArray, cvariant, sizeof(PRUint32) * variantCount);
    HRESULT result = IMedium_CreateBaseStorage(cmedium, size, ComSafeArrayAsInParam(pSafeArray), cprogress);
    g_pVBoxFuncs->pfnSafeArrayDestroy(pSafeArray);
    return result;
}

HRESULT VboxMediumDeleteStorage(IMedium* cmedium, IProgress** cprogress) {
    return IMedium_DeleteStorage(cmedium, cprogress);
}

HRESULT VboxMediumClose(IMedium* cmedium) {
    return IMedium_Close(cmedium);
}

HRESULT VboxGetMediumLocation(IMedium* cmedium, char** clocation) {
    BSTR wlocation = NULL;
    HRESULT result = IMedium_GetLocation(cmedium, &wlocation);
    if (FAILED(result)) {
        return result;
    }

    g_pVBoxFuncs->pfnUtf16ToUtf8(wlocation, clocation);
    g_pVBoxFuncs->pfnComUnallocString(wlocation);
    return result;
}

HRESULT VboxGetMediumState(IMedium* cmedium, PRUint32* cstate) {
    return IMedium_GetState(cmedium, cstate);
}

HRESULT VboxGetMediumSize(IMedium* cmedium, PRInt64* csize) {
    return IMedium_GetSize(cmedium, csize);
}

HRESULT VboxIMediumRelease(IMedium* cmedium) {
    HRESULT result = IMedium_Release(cmedium);
    cmedium = NULL;
    return result;
}
