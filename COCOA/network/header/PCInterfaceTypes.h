//
//  LinkInterfaceTypes.h
//  LinkDetector
//
//  Created by Almighty Kim on 10/20/16.
//  Copyright (c) 2016 PocketCluster.io. All rights reserved.
//

#ifndef __PCINTERFACETYPES_H__
#define __PCINTERFACETYPES_H__

#include "SCNetworkTypes.h"

typedef struct _PCNetworkInterface {
    bool                wifiPowerOff;
    bool                isActive;
    bool                isPrimary;
    unsigned int        addrCount;
    SCNIAddress**       address;

    const char*         bsdName;
    const char*         displayName;
    const char*         macAddress;
    const char*         mediaType;
} PCNetworkInterface;

typedef bool (*pc_interface_callback)(PCNetworkInterface**, unsigned int);

extern void
interface_status_with_callback(pc_interface_callback);

typedef bool (*scni_gateway_callback)(SCNIGateway**, unsigned int);

extern void
gateway_status_with_callback(scni_gateway_callback);

#endif
