//
//  DeviceSerialNumber.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "DeviceSerialNumber.h"

#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/IOKitLib.h>

// Returns the serial number as a CFString.
// It is the caller's responsibility to release the returned CFString when done with it.
CFTypeRef CopySerialNumber()
{
    CFTypeRef serialNumberCFString;
    
    io_service_t    platformExpert = IOServiceGetMatchingService(kIOMasterPortDefault, IOServiceMatching("IOPlatformExpertDevice"));
    if (platformExpert) {
        serialNumberCFString = IORegistryEntryCreateCFProperty(platformExpert, CFSTR(kIOPlatformSerialNumberKey), kCFAllocatorDefault, 0);
        IOObjectRelease(platformExpert);
    }
    
    return serialNumberCFString;
}
