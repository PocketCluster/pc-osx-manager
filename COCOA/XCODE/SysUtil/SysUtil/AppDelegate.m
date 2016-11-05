//
//  AppDelegate.m
//  SysUtil
//
//  Created by Almighty Kim on 10/24/16.
//  Copyright © 2016 PocketCluster. All rights reserved.
//

#import "AppDelegate.h"
#import "NSResourcePath.h"
#import "PCDeviceSerial.h"
#import "PCUserEnvironment.h"


@interface AppDelegate ()

@property (weak) IBOutlet NSWindow *window;
@end

@implementation AppDelegate

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {

    printf("%s\n", PCEnvironmentCocoaHomeDirectory());
    printf("%s\n", PCEnvironmentPosixHomeDirectory());
    printf("%s\n", PCEnvironmentFullUserName());
    printf("%s\n", PCEnvironmentLoginUserName());
    printf("%s\n", PCEnvironmentUserTemporaryDirectory());
    
    printf("%s\n", PCApplicationSupportDirectory());
    printf("%s\n", PCApplicationDocumentsDirectory());
    printf("%s\n", PCApplicationTemporaryDirectory());
    printf("%s\n", PCApplicationLibraryCacheDirectory());
    printf("%s\n", PCApplicationResourceDirectory());
    printf("%s\n", PCApplicationExecutableDirectory());
    
    printf("%s\n", PCDeviceSerialNumber());
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    // Insert code here to tear down your application
}

@end
