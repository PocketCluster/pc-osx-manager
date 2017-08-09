#include "VBoxCAPIGlue.h"
#include "common.h"

HRESULT VboxGetUsbControllerName(IUSBController* ccontroller, char** cname) {
    BSTR wname = NULL;
    HRESULT result = IUSBController_GetName(ccontroller, &wname);
    if (FAILED(result)) {
        return result;
    }

    g_pVBoxFuncs->pfnUtf16ToUtf8(wname, cname);
    g_pVBoxFuncs->pfnComUnallocString(wname);
    return result;
}

HRESULT VboxGetUsbControllerStandard(IUSBController* ccontroller, PRUint16* cstandard) {
    return IUSBController_GetUSBStandard(ccontroller, cstandard);
}

HRESULT VboxGetUsbControllerType(IUSBController* ccontroller, PRUint32* ctype) {
    return IUSBController_GetType(ccontroller, ctype);
}

HRESULT VboxIUSBControllerRelease(IUSBController* ccontroller) {
    return IUSBController_Release(ccontroller);
}

HRESULT VboxMachineAddUsbController(IMachine* cmachine, char* cname, PRUint32 ccontrollerType, IUSBController** ccontroller) {
    BSTR wname;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cname, &wname);
    if (FAILED(result)) {
        return result;
    }

    result = IMachine_AddUSBController(cmachine, wname, ccontrollerType, ccontroller);
    g_pVBoxFuncs->pfnUtf16Free(wname);

    return result;
}
