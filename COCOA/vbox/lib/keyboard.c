#include "../include/common.h"
#include "../include/keyboard.h"

HRESULT VboxKeyboardPutScancodes(IKeyboard* ckeyboard, PRUint32 scancodesCount, PRInt32* cscancodes, PRUint32* ccodesStored) {
    SAFEARRAY *pSafeArray = g_pVBoxFuncs->pfnSafeArrayCreateVector(VT_I4, 0, scancodesCount);
    g_pVBoxFuncs->pfnSafeArrayCopyInParamHelper(pSafeArray, cscancodes, sizeof(PRInt32) * scancodesCount);
    HRESULT result = IKeyboard_PutScancodes(ckeyboard, ComSafeArrayAsInParam(pSafeArray), ccodesStored);
    g_pVBoxFuncs->pfnSafeArrayDestroy(pSafeArray);
    return result;
}

HRESULT VboxIKeyboardRelease(IKeyboard* ckeyboard) {
    return IKeyboard_Release(ckeyboard);
}
