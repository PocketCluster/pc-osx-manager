//
//  main.c
//  netifaces
//
//  Created by Almighty Kim on 10/2/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <string.h>
#include "../../netifaces.h"

#ifdef TEST_NETIFACES_GATEWAY_ALLOCATION

// make sure add_to_gatways is not static!
extern bool add_to_gatways(Gateway** results, Gateway* gateway);

void test_gateway_allocation() {
    Gateway *results = NULL;
    char buf[256];
    for (int i = 0; i < 100; i++) {
        Gateway *gw = (Gateway *) calloc(1, sizeof(Gateway));
        
        sprintf(buf, "addr %d", i);
        gw->addr = malloc(sizeof(char) * strlen(buf));
        memcpy(gw->addr, buf, strlen(buf));
        
        sprintf(buf, "ifname %d", i);
        gw->ifname = malloc(sizeof(char) * strlen(buf));
        memcpy(gw->ifname, buf, strlen(buf));
        
        if  ((i % 3) == 0) {
            gw->family = 30; // ip6
            if (i == 81) {
                gw->is_default = true;
            }
        } else {
            gw->family = 2; // ip4
            if (i == 22) {
                gw->is_default = true;
            }
        }
        
        add_to_gatways(&results, gw);
    }
    
    assert(strcmp(find_default_ip4_gw(&results)->ifname, "ifname 22") == 0);
    assert(strcmp(find_default_ip6_gw(&results)->ifname, "ifname 81") == 0);
    release_gateways(&results);
    printf("check complete!\n");
}
#endif

int main(int argc, const char * argv[]) {

    Gateway *results = NULL;
    find_system_gateways(&results);
    if (results == NULL) {
        printf("result is null\n");
        return 0;
    }
    Gateway* default_gw = find_default_ip4_gw(&results);
    printf("default ip4 gw %s ifname %s\n", default_gw->addr, default_gw->ifname);
    release_gateways_info(&results);
    
    
    Interface *interfaces = NULL, *iface = NULL;
    Address *address = NULL;
    int err = find_system_interfaces(&interfaces);
    if (err != 0) {
        printf("error in interface info acquisition\n");
    }
    if (interfaces == NULL) {
        printf("interfaces interface is null\n");
        return 0;
    }
    
    iface = interfaces;
    while (iface != NULL) {
        printf("iface name %s\n",iface->name);
        
        address = iface->address;
        while (address != NULL) {
            printf("\taddr %s | flags 0x%X\n", address->addr, address->flags);
            address = address->next;
        }

        iface = iface->next;
    }
    
    release_interfaces_info(&interfaces);

    return 0;
}

