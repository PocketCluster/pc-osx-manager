//
//  DeviceSerialNumber.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <CoreFoundation/CoreFoundation.h>
#import <IOKit/IOKitLib.h>
#import "DeviceSerialNumber.h"
#import "PCDeviceSerial.h"

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

@implementation DeviceSerialNumber

const char*
PCDeviceSerialNumber(void) {
    __autoreleasing NSString *sn = [(__bridge NSString*)((CFStringRef)CopySerialNumber()) uppercaseString];
    return [sn UTF8String];
}

+ (NSString *)deviceSerialNumber {
    NSString *sn = [(__bridge NSString*)((CFStringRef)CopySerialNumber()) lowercaseString];
    return sn;
}

+ (NSString *)UUIDString {
    CFUUIDRef theUUID = CFUUIDCreate(NULL);
    CFStringRef string = CFUUIDCreateString(NULL, theUUID);
    CFRelease(theUUID);

    // TODO: __bridge_transfer?
    return (__bridge NSString*)string;
}

+(CFUUIDBytes)UUIDBytes {
    CFUUIDRef theUUID = CFUUIDCreate(NULL);
    CFUUIDBytes bytes = CFUUIDGetUUIDBytes(theUUID);
    CFRelease(theUUID);
    return bytes;
}
@end
