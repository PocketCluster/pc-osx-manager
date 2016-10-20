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

typedef struct _SCNIAddress {
    unsigned int       flags;
    unsigned char      family;
    char*              addr;
    char*              netmask;
    char*              broadcast;
    char*              peer;
} SCNIAddress;

typedef struct _SCNIGateway {
    unsigned char      family;
    bool               is_default;
    char*              ifname;
    char*              addr;
} SCNIGateway;


#endif
