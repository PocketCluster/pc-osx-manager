#include "../include/common.h"
#include "../include/medium_format.h"

HRESULT VboxGetMediumFormats(ISystemProperties* cprops, IMediumFormat*** cformats, ULONG* formatCount) {
    SAFEARRAY *safeArray = g_pVBoxFuncs->pfnSafeArrayOutParamAlloc();
    HRESULT result = ISystemProperties_GetMediumFormats(cprops, ComSafeArrayAsOutIfaceParam(safeArray, IMediumFormat *));
    if (!FAILED(result)) {
        result = g_pVBoxFuncs->pfnSafeArrayCopyOutIfaceParamHelper((IUnknown ***)cformats, formatCount, safeArray);
    }
    g_pVBoxFuncs->pfnSafeArrayDestroy(safeArray);
    return result;
}

HRESULT VboxGetMediumFormatId(IMediumFormat* cformat, char** cid) {
    BSTR wid = NULL;
    HRESULT result = IMediumFormat_GetId(cformat, &wid);
    if (FAILED(result)) {
        return result;
    }

    g_pVBoxFuncs->pfnUtf16ToUtf8(wid, cid);
    g_pVBoxFuncs->pfnComUnallocString(wid);
    return result;
}

HRESULT VboxIMediumFormatRelease(IMediumFormat* cformat) {
    return IMediumFormat_Release(cformat);
}
