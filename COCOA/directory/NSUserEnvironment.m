//
//  NSUserEnvironment.m
//  SysUtil
//
//  Created by Almighty Kim on 10/24/16.
//  Copyright © 2016 PocketCluster. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "PCUSerEnvironment.h"

const char*
PCEnvironmentHomeDirectory(void) {
    NSString* path = NSHomeDirectory();
    return [path UTF8String];
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