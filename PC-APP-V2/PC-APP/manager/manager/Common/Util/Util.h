//
//  Util.h
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#import "AppDelegate.h"

@class Package;

@interface Util : NSObject
+ (AppDelegate*)getApp;
+ (NSString*)escapeShellArg:(NSString*)arg;
+ (NSString*)trimTrailingSlash:(NSString*)path;
+ (void)redirectConsoleLogToDocumentFolder;
+ (NSString*)getMachineId;
+ (BOOL)shouldSendProfileData;
+ (void)runTerminalCommand:(NSString*)command;
+ (void)openTerminalCommand:(Package*)aPackage;
@end
