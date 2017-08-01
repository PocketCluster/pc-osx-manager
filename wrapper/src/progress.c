#include "VBoxCAPIGlue.h"
#include "common.h"

HRESULT VboxProgressWaitForCompletion(IProgress* cprogress, int timeout) {
    return IProgress_WaitForCompletion(cprogress, timeout);
}

HRESULT VboxGetProgressPercent(IProgress* cprogress, PRUint32* cpercent) {
    return IProgress_GetPercent(cprogress, cpercent);
}

HRESULT VboxGetProgressResultCode(IProgress* cprogress, PRInt32* code) {
    return IProgress_GetResultCode(cprogress, code);
}

HRESULT VboxIProgressRelease(IProgress* cprogress) {
    HRESULT result = IProgress_Release(cprogress);
    cprogress = NULL;
    return result;
}
