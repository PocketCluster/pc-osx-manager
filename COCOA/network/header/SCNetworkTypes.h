//
//  SCNetworkStruct.h
//  LinkDetector
//
//  Created by Sung-Taek,Kim on 10/20/16.
//  Copyright (c) 2016 PocketCluster.io. All rights reserved.
//

#ifndef __SCNETWORKTYPES_H__
#define __SCNETWORKTYPES_H__

#include <stdbool.h>

// TODO : we're supposed to a platform agnostic enum for family as the values of AF_INET/AF_INET6 differs from platform to platform
typedef struct _SCNIAddress {
    unsigned int       flags;
    unsigned char      family;
    bool               is_primary;
    const char*        addr;
    const char*        netmask;
    const char*        broadcast;
    const char*        peer;
} SCNIAddress;

typedef struct _SCNIGateway {
    unsigned char      family;
    bool               is_default;
    const char*        ifname;
    const char*        addr;
} SCNIGateway;


#endif
