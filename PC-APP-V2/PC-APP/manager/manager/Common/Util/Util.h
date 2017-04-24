//
//  Util.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "AppDelegate.h"

@class PCPackageMeta;

@interface Util : NSObject
+ (AppDelegate*)getApp;
+ (NSString*)escapeShellArg:(NSString*)arg;
+ (NSString*)trimTrailingSlash:(NSString*)path;
+ (void)redirectConsoleLogToDocumentFolder;
+ (NSString*)getMachineId;
+ (BOOL)shouldSendProfileData;
+ (void)runTerminalCommand:(NSString*)command;
+ (void)openTerminalCommand:(PCPackageMeta*)aPackage;
@end
