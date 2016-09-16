//
//  AppDelegate.m
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <Sparkle/Sparkle.h>

#ifdef USE_LIBSSH2
    #import <NMSSH/NMSSH.h>
#endif

#import "Util.h"
#import "AppDelegate.h"
#import "NativeMenu.h"

@interface AppDelegate ()<SUUpdaterDelegate, NSUserNotificationCenterDelegate>
@property (nonatomic, strong, readwrite) NativeMenu *nativeMenu;
@property (nonatomic, strong) NSMutableArray *openWindows;
@end

@implementation AppDelegate

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    
    // opened window list
    self.openWindows = [[NSMutableArray alloc] init];
    
    //initialize updates
    [[SUUpdater sharedUpdater] setDelegate:self];
    [[SUUpdater sharedUpdater] setSendsSystemProfile:[Util shouldSendProfileData]];
    [[SUUpdater sharedUpdater] checkForUpdateInformation];
    
    //create popup and status menu item
    self.nativeMenu = [[NativeMenu alloc] init];
    
    Log(@"Application Started");
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
}

- (void)application:(NSApplication *)application didReceiveRemoteNotification:(NSDictionary *)userInfo {
}

#pragma mark - WINDOW MANAGEMENT
- (void)addOpenWindow:(id)window {
}

- (void)removeOpenWindow:(id)window {
}
@end
