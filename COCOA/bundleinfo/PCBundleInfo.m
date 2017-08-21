//
//  PCBundleInfo.m
//  static-core
//
//  Created by Almighty Kim on 8/19/17.
//  Copyright © 2017 PocketCluster. All rights reserved.
//

#import "PCBundleInfo.h"
#import <Foundation/Foundation.h>

const char*
PCBundleVersionString() {
    return [[[[NSBundle mainBundle] infoDictionary] objectForKey:@"CFBundleShortVersionString"] UTF8String];
}

const char*
PCBundleExpirationString() {
    return [[[[NSBundle mainBundle] infoDictionary] objectForKey:@"PCBundleExpirationString"] UTF8String];
}