//
//  main.m
//  LinkDetector
//
//  Created by Almighty Kim on 10/17/16.
//  Copyright (c) 2016 PocketCluster.io. All rights reserved.
//

#import "PCInterfaceStatus.h"

bool
pc_interface_list(PCNetworkInterface** inetfaces, unsigned int count) {
    NSLog(@"total address count %d", count);
    return true;
}

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        interface_status(&pc_interface_list);
        sleep(1);
    }
    return 0;
}
