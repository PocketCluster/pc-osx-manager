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
    int                 bsdNumber;
    bool                wifiPowerOff;
    bool                isActive;
    
    unsigned int        addrCount;
    SCNIAddress**       address;

    const char*         bsdName;
    const char*         displayName;
    const char*         macAddress;
    const char*         mediaType;
} PCNetworkInterface;

typedef bool (*pc_interface_callback)(PCNetworkInterface**, unsigned int);

CF_EXPORT void
interface_status(pc_interface_callback callback);

#endif
