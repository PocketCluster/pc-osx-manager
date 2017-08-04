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
    result = INetworkAdapter_SetBridgedInterface(cAdapter, wHostInterface);
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

HRESULT VboxNetworkAddPortForwardingRule(INetworkAdapter *cAdapter, const char* cRuleName, NATProtocol protocol, const char* cHostIp, unsigned short hostPort, const char* cGuestIp, unsigned short guestPort) {
    INATEngine *natEngine = NULL;
    HRESULT result;
    
    BSTR wRuleName;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cRuleName, &wRuleName);
    if (FAILED(result)) {
        return result;
    }
    
    BSTR wHostIp;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cHostIp, &wHostIp);
    if (FAILED(result)) {
        return result;
    }
    
    BSTR wGuestIp;
    result = g_pVBoxFuncs->pfnUtf8ToUtf16(cGuestIp, &wGuestIp);
    if (FAILED(result)) {
        return result;
    }
    
    // get NAT engine
    result = INetworkAdapter_get_NATEngine(cAdapter, &natEngine);
    if (FAILED(result)) {
        return result;
    }

    // add port forwarding rule.
    // Even if this fails we need to go on to free memories
    result = INATEngine_AddRedirect(natEngine, wRuleName, protocol, wHostIp, hostPort, wGuestIp, guestPort);

    // the DNS Host Resolver doesn't support SRV records
    // the DNS proxy has performance issues
    // direct DNS pass-through doesn't support roaming laptops well
    // we can't win, so let's go direct and at least get performance
    result = INATEngine_SetDNSProxy(natEngine, PR_FALSE);
    result = INATEngine_SetDNSUseHostResolver(natEngine, PR_FALSE);
    
    // release NAT engine (we'll skip the result to pass important results)
    INATEngine_Release(natEngine);

    g_pVBoxFuncs->pfnUtf16Free(wGuestIp);
    g_pVBoxFuncs->pfnUtf16Free(wHostIp);
    g_pVBoxFuncs->pfnUtf16Free(wRuleName);
    
    return result;
}

HRESULT VboxNetworkAdapterRelease(INetworkAdapter *cAdapter) {
    HRESULT result = INetworkAdapter_Release(cAdapter);
    cAdapter = NULL;
    return result;
}
