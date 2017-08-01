#include "../include/common.h"
#include "../include/network.h"

HRESULT VboxNetworkAdapterSetEnabled(INetworkAdapter *cAdapter, BOOL isEnabled) {
    PRBool adapter_enabled = (PRBool)(isEnabled);
    return INetworkAdapter_SetEnabled(cAdapter, adapter_enabled);
}

HRESULT VboxNetworkAdapterSetAttachmentType(INetworkAdapter *cAdapter, NetworkAttachmentType attachementType) {
    return INetworkAdapter_SetAttachmentType(cAdapter, attachementType);
}

HRESULT VboxNetworkAdapterSetBridgedHostInterface(INetworkAdapter *cAdapter, const char* cHostInterface) {
    BSTR wHostInterface;
    HRESULT result = g_pVBoxFuncs->pfnUtf8ToUtf16(cHostInterface, &wHostInterface);
    if (FAILED(result)) {
        return result;
    }
    result = INetworkAdapter_SetBridgedInterface(cAdapter, NULL);
    g_pVBoxFuncs->pfnUtf16Free(wHostInterface);
    return result;
}

HRESULT VboxNetworkAdapterSetAdapterType(INetworkAdapter *cAdapter, NetworkAdapterType adapterType) {
    return INetworkAdapter_SetAdapterType(cAdapter, adapterType);
}

HRESULT VboxNetworkAdapterSetPromiscModePolicy(INetworkAdapter *cAdapter, NetworkAdapterPromiscModePolicy policyType) {
    return INetworkAdapter_SetPromiscModePolicy(cAdapter, policyType);
}

HRESULT VboxNetworkAdapterSetCableConnected(INetworkAdapter *cAdapter, BOOL isConnected) {
    PRBool cable_connected = (PRBool)(isConnected);
    return INetworkAdapter_SetCableConnected(cAdapter, cable_connected);
}

HRESULT VboxNetworkAdapterRelease(INetworkAdapter *cAdapter) {
    HRESULT result = INetworkAdapter_Release(cAdapter);
    cAdapter = NULL;
    return result;
}
