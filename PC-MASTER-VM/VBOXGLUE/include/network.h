//
//  network.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __NETWORK_H__
#define __NETWORK_H__

HRESULT VboxNetworkAdapterSetEnabled(INetworkAdapter *cAdapter, BOOL isEnabled);

HRESULT VboxNetworkAdapterSetAttachmentType(INetworkAdapter *cAdapter, NetworkAttachmentType attachementType);

HRESULT VboxNetworkAdapterSetBridgedHostInterface(INetworkAdapter *cAdapter, const char* cHostInterface);

HRESULT VboxNetworkAdapterSetAdapterType(INetworkAdapter *cAdapter, NetworkAdapterType adapterType);

HRESULT VboxNetworkAdapterSetPromiscModePolicy(INetworkAdapter *cAdapter, NetworkAdapterPromiscModePolicy policyType);

HRESULT VboxNetworkAdapterSetCableConnected(INetworkAdapter *cAdapter, BOOL isConnected);

HRESULT VboxNetworkAdapterRelease(INetworkAdapter *cAdapter);

#endif /* network_h */
