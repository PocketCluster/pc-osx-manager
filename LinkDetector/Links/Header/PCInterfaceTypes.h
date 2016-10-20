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

typedef struct _pc_interface {
    int                 bsd_number;
    bool                wifi_power_off;
    
    unsigned int        address_length;
    SCNIAddress**       address;

    char*               bsd_name;
    char*               display_name;
    char*               mac_address;
    char*               media_type;
} pc_interface;

typedef bool (*pc_interface_callback)(pc_interface**, unsigned int);

extern void interface_status(pc_interface_callback callback);

#endif
