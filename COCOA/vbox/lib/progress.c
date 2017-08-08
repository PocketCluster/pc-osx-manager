#include "../include/common.h"
#include "../include/progress.h"

HRESULT VboxProgressWaitForCompletion(IProgress* cprogress, int timeout) {
    return IProgress_WaitForCompletion(cprogress, timeout);
}

HRESULT VboxProgressGetPercent(IProgress* cprogress, PRUint32* cpercent) {
    return IProgress_GetPercent(cprogress, cpercent);
}

HRESULT VboxProgressGetResultCode(IProgress* cprogress, PRInt32* code) {
    return IProgress_GetResultCode(cprogress, code);
}

HRESULT VboxProgressGetResultInfo(IProgress* cprogress, char** cErrorMessage) {
    IVirtualBoxErrorInfo *errorInfo = NULL;
    BSTR wErrorMessage = NULL;
    HRESULT result;

    result = IProgress_get_ErrorInfo(cprogress, &errorInfo);
    if (FAILED(result) || errorInfo == NULL) {
        return result;
    }
    result = IVirtualBoxErrorInfo_get_Text(errorInfo, &wErrorMessage);

    g_pVBoxFuncs->pfnUtf16ToUtf8(wErrorMessage, cErrorMessage);
    g_pVBoxFuncs->pfnComUnallocString(wErrorMessage);
    return result;
}

HRESULT VboxProgressRelease(IProgress* cprogress) {
    HRESULT result = IProgress_Release(cprogress);
    cprogress = NULL;
    return result;
}
