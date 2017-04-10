//
//  VersionComparison.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "VersionComparison.h"

@implementation VersionComparison

- (NSComparisonResult)compareVersion:(NSString *)versionA toVersion:(NSString *)versionB {

    //return [Util compareVersion:versionA toVersion:versionB];
    //    if ([[json valueForKeyPath:@"osx.version"] compare:[[NSBundle mainBundle] objectForInfoDictionaryKey:(NSString *)kCFBundleVersionKey] options:NSNumericSearch] == NSOrderedDescending)

    return [versionB compare:versionA options:NSNumericSearch];
}

@end
