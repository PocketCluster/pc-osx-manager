#include "../include/common.h"
#include "../include/storage_controller.h"

HRESULT VboxGetStorageControllerName(IStorageController* ccontroller, char** cname) {
    BSTR wname = NULL;
    HRESULT result = IStorageController_GetName(ccontroller, &wname);
    if (FAILED(result)) {
        return result;
    }

    g_pVBoxFuncs->pfnUtf16ToUtf8(wname, cname);
    g_pVBoxFuncs->pfnComUnallocString(wname);
    return result;
}

HRESULT VboxGetStorageControllerBus(IStorageController* ccontroller, PRUint32* cbus) {
    return IStorageController_GetBus(ccontroller, cbus);
}

HRESULT VboxGetStorageControllerType(IStorageController* ccontroller, PRUint32* ctype) {
    return IStorageController_GetControllerType(ccontroller, ctype);
}

HRESULT VboxSetStorageControllerType(IStorageController* ccontroller, PRUint32 ctype) {
    return IStorageController_SetControllerType(ccontroller, ctype);
}

HRESULT VboxIStorageControllerRelease(IStorageController* ccontroller) {
    HRESULT result = IStorageController_Release(ccontroller);
    ccontroller = NULL;
    return result;
}

HRESULT VboxMachineAddStorageController(IMachine* cmachine, char* cname, PRUint32 connectionType, IStorageController** ccontroller) {
    BSTR wname;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cname, &wname);
    if (FAILED(result)) {
        return result;
    }

    result = IMachine_AddStorageController(cmachine, wname, connectionType, ccontroller);
    g_pVBoxFuncs->pfnUtf16Free(wname);

    return result;
}
