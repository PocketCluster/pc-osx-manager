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

    PCApplicationSupportDirectory();
    PCApplicationDocumentsDirectory();
    PCApplicationTemporaryDirectory();
    PCApplicationLibraryCacheDirectory();
    PCApplicationResourceDirectory();
    PCApplicationExecutableDirectory();
    PCDeviceSerialNumber();
    PCEnvironmentCocoaHomeDirectory();

    PCEnvironmentPosixHomeDirectory();
    PCEnvironmentFullUserName();
    PCEnvironmentUserTemporaryDirectory();
    PCEnvironmentLoginUserName();
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    // Insert code here to tear down your application
}

@end
