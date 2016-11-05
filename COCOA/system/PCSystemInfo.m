//
//  PCSystemInfo.m
//  SysUtil
//
//  Created by Almighty Kim on 10/27/16.
//  Copyright © 2016 PocketCluster. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "PCSystemInfo.h"

unsigned long
PCSystemProcessorCount(void) {
    return [[NSProcessInfo processInfo] processorCount];
}

unsigned long
PCSystemActiveProcessorCount(void) {
    return [[NSProcessInfo processInfo] activeProcessorCount];
}

unsigned long long
PCSystemPhysicalMemorySize(void) {
    return [[NSProcessInfo processInfo] physicalMemory];
}
