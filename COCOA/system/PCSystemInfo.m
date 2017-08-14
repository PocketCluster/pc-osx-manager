//
//  PCSystemInfo.m
//  SysUtil
//
//  Created by Almighty Kim on 10/27/16.
//  Copyright © 2016 PocketCluster. All rights reserved.
//

#import <Foundation/Foundation.h>
#include <sys/sysctl.h>
#import "PCSystemInfo.h"

#import <stdio.h>

unsigned long
PCSystemProcessorCount(void) {
    return [[NSProcessInfo processInfo] processorCount];
}

unsigned long
PCSystemActiveProcessorCount(void) {
    return [[NSProcessInfo processInfo] activeProcessorCount];
}

/*
 *   More information is available on <sys/sysctl.h>
 *   
 *   hw.memsize                - The number of bytes of physical memory in the system.
 *
 */

unsigned long long
PCSystemPhysicalMemorySize(void) {
#ifdef USE_NSPROCESSINFO
    return [[NSProcessInfo processInfo] physicalMemory];
#else
    unsigned long long size = 0;
    size_t size_len = sizeof(size);
    sysctlbyname("hw.memsize", &size, &size_len, NULL, 0);
    return size;
#endif
}

/*
 *   More information is available on <sys/sysctl.h>.
 *   We can also get information with 'sysctl -h hw'
 *
 *   hw.ncpu                   - The maximum number of processors that could be available this boot.
 *                               Use this value for sizing of static per processor arrays; i.e. processor load statistics.
 *
 *   hw.activecpu              - The number of processors currently available for executing threads.
 *                               Use this number to determine the number threads to create in SMP aware applications.
 *                               This number can change when power management modes are changed.
 *
 *   hw.physicalcpu            - The number of physical processors available in the current power management mode.
 *   hw.physicalcpu_max        - The maximum number of physical processors that could be available this boot
 *
 *   hw.logicalcpu             - The number of logical processors available in the current power management mode.
 *   hw.logicalcpu_max         - The maximum number of logical processors that could be available this boot
 *
 *   Reference                 - https://stackoverflow.com/questions/1715580/how-to-discover-number-of-logical-cores-on-mac-os-x
 */

unsigned long
PCSystemPhysicalCoreCount(void) {
    unsigned long count = 0;
    size_t count_len = sizeof(count);
    sysctlbyname("hw.physicalcpu_max", &count, &count_len, NULL, 0);
    return count;
}

// https://stackoverflow.com/questions/150355/programmatically-find-the-number-of-cores-on-a-machine
// https://stackoverflow.com/questions/2901694/programmatically-detect-number-of-physical-processors-cores-or-if-hyper-threadin/2921632
#ifdef SYSTEM_AGONISTIC
    #ifdef _WIN32
    #include <windows.h>
    #elif MACOS
    #include <sys/param.h>
    #include <sys/sysctl.h>
    #else
    #include <unistd.h>
    #endif

    int getNumCores() {
    #ifdef WIN32
        SYSTEM_INFO sysinfo;
        GetSystemInfo(&sysinfo);
        return sysinfo.dwNumberOfProcessors;
    #elif MACOS
        int nm[2];
        size_t len = 4;
        uint32_t count;
        
        nm[0] = CTL_HW; nm[1] = HW_AVAILCPU;
        sysctl(nm, 2, &count, &len, NULL, 0);
        
        if(count < 1) {
            nm[1] = HW_NCPU;
            sysctl(nm, 2, &count, &len, NULL, 0);
            if(count < 1) { count = 1; }
        }
        return count;
    #else
        return sysconf(_SC_NPROCESSORS_ONLN);
    #endif
    }
#endif