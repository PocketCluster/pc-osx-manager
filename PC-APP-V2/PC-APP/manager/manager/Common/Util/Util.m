//
//  Util.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "Util.h"
@implementation Util

+ (AppDelegate*)getApp {
    return (AppDelegate*)[[NSApplication sharedApplication] delegate];
}

+ (NSString*)escapeShellArg:(NSString*)arg {
    NSMutableString *result = [NSMutableString stringWithString:arg];
    [result replaceOccurrencesOfString:@"'" withString:@"'\\''" options:0 range:NSMakeRange(0, [result length])];
    [result insertString:@"'" atIndex:0];
    [result appendString:@"'"];
    return result;
}

+ (NSString*)trimTrailingSlash:(NSString*)path {
    if(path.length > 1 && [[path substringFromIndex:path.length-1] isEqualToString:@"/"]) {
        return [path substringToIndex:path.length-1];
    } else {
        return path;
    }
}

+ (void)redirectConsoleLogToDocumentFolder {
    NSArray *paths = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES);
    
    NSDateFormatter *formatter = [[NSDateFormatter alloc] init];
    [formatter setDateFormat:@"yyyyMMdd-HHmmss"];
    
    NSString *documentsDirectory = [paths objectAtIndex:0];
    NSString *logPath = [documentsDirectory stringByAppendingPathComponent:[NSString stringWithFormat:@"vagrant-manager-%@.log", [formatter stringFromDate:[NSDate date]]]];
    freopen([logPath fileSystemRepresentation],"a+",stderr);
}

+ (NSString*)getMachineId {
    NSString *uuid = [[NSUserDefaults standardUserDefaults] stringForKey:@"machineId"];
    if(!uuid) {
        uuid = [[NSUUID UUID] UUIDString];
        [[NSUserDefaults standardUserDefaults] setObject:uuid forKey:@"machineId"];
        [[NSUserDefaults standardUserDefaults] synchronize];
    }
    
    return uuid;
}

+ (BOOL)shouldSendProfileData {
    if([[NSUserDefaults standardUserDefaults] objectForKey:@"sendProfileData"] == nil) {
        return YES;
    }
    return [[NSUserDefaults standardUserDefaults] boolForKey:@"sendProfileData"];
}

+ (void)runTerminalCommand:(NSString*)command {
}

+ (void)openTerminalCommand:(PCPackageMeta*)aPackage {
}

@end
