#include "VBoxCAPIGlue.h"
#include "common.h"

HRESULT VboxGetGuestOSTypes(IVirtualBox* cbox, IGuestOSType*** ctypes, ULONG* typeCount) {
    SAFEARRAY *safeArray = g_pVBoxFuncs->pfnSafeArrayOutParamAlloc();
    HRESULT result = IVirtualBox_GetGuestOSTypes(cbox, ComSafeArrayAsOutIfaceParam(safeArray, IGuestOSType *));
    if (!FAILED(result)) {
        result = g_pVBoxFuncs->pfnSafeArrayCopyOutIfaceParamHelper((IUnknown ***)ctypes, typeCount, safeArray);
    }
    g_pVBoxFuncs->pfnSafeArrayDestroy(safeArray);
    return result;
}

HRESULT VboxGetGuestOSTypeId(IGuestOSType* ctype, char** cid) {
    BSTR wid = NULL;
    HRESULT result = IGuestOSType_GetId(ctype, &wid);
    if (FAILED(result)) {
        return result;
    }

    g_pVBoxFuncs->pfnUtf16ToUtf8(wid, cid);
    g_pVBoxFuncs->pfnComUnallocString(wid);
    return result;
}

HRESULT VboxIGuestOSTypeRelease(IGuestOSType* ctype) {
    return IGuestOSType_Release(ctype);
}
