//
//  host.h
//  VBoxGlue
//
//  Created by Almighty Kim on 8/4/17.
//  Copyright © 2017 PocketCluster. All rights reserved.
//

#ifndef __HOST_H__
#define __HOST_H__

HRESULT VboxGetHostFromVirtualbox(IVirtualBox *cVbox, IHost** host);

HRESULT VboxHostSearchNetworkInterfacesFromList(IHost* cHost, const char* cQueryName, char** cFullname);

HRESULT VboxHostFindNetworkInterfaceByID(IHost* cHost, const char* cIfID, IHostNetworkInterface** interfaces);

// this needs exact name to match
HRESULT VboxHostFindNetworkInterfaceByName(IHost* cHost, const char* cName, IHostNetworkInterface** interfaces);

HRESULT VboxGetHostNetworkInterfaceName(IHostNetworkInterface* cIface, char** cFullname);

#if defined(VBOX_BIND_BUG_FIXED)

HRESULT VboxHostFindNetworkInterfaceByType(IHost* cHost, HostNetworkInterfaceType type, IHostNetworkInterface** interfaces);
#endif

#endif /* host.h */
