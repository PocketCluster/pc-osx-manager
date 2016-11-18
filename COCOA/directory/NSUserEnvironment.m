//
//  NSUserEnvironment.m
//  SysUtil
//
//  Created by Almighty Kim on 10/24/16.
//  Copyright © 2016 PocketCluster. All rights reserved.
//

#include <unistd.h>
#include <sys/types.h>
#include <pwd.h>
#import <Foundation/Foundation.h>
#import "PCUSerEnvironment.h"

const char*
PCEnvironmentCocoaHomeDirectory(void) {
    NSString* path = NSHomeDirectory();
    return [path UTF8String];
}

extern const char*
PCEnvironmentPosixHomeDirectory(void) {
    struct passwd *pw = getpwuid(getuid());
    return pw->pw_dir;
}

const char*
PCEnvironmentFullUserName(void) {
    NSString* user = NSFullUserName();
    return [user UTF8String];
}

const char*
PCEnvironmentUserTemporaryDirectory(void) {
    NSString* path = NSTemporaryDirectory();
    return [path UTF8String];
}

const char*
PCEnvironmentLoginUserName(void) {
    NSString* user = NSUserName();
    return [user UTF8String];
}

const char*
PCEnvironmentCurrentCountryCode(void) {
    return [[[NSLocale currentLocale] objectForKey:NSLocaleCountryCode] UTF8String];
}

const char*
PCEnvironmentCurrentLanguageCode(void) {
    return [[[[NSLocale currentLocale] objectForKey:NSLocaleLanguageCode] uppercaseString] UTF8String];
}