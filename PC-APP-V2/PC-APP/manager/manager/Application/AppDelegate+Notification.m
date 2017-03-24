//
//  AppDelegate+Notification.m
//  manager
//
//  Created by Almighty Kim on 3/24/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+Notification.h"
#include "_cgo_export.h"

@interface AppDelegate (Notification_Private)
- (void) receiveSleepNote: (NSNotification*) notification;
- (void) receiveWakeNote: (NSNotification*) notification;
@end

@implementation AppDelegate (Notification)
- (void) receiveSleepNote: (NSNotification*) notification {
    Log(@"receiveSleepNote: %@", [notification name]);
}

- (void) receiveWakeNote: (NSNotification*) notification {
    Log(@"receiveWakeNote: %@", [notification name]);
}

- (void) addSleepNotifications {
    // These notifications are filed on NSWorkspace's notification center, not
    // the default notification center. You will not receive sleep/wake
    // notifications if you file with the default notification center.
    [[[NSWorkspace sharedWorkspace] notificationCenter]
     addObserver: self
     selector: @selector(receiveSleepNote:)
     name: NSWorkspaceWillSleepNotification object: NULL];
    
    [[[NSWorkspace sharedWorkspace] notificationCenter]
     addObserver: self
     selector: @selector(receiveWakeNote:)
     name: NSWorkspaceDidWakeNotification object: NULL];
}

- (void) removeSleepNotifications {
    // These notifications are filed on NSWorkspace's notification center, not
    // the default notification center. You will not receive sleep/wake
    // notifications if you file with the default notification center.
    [[[NSWorkspace sharedWorkspace] notificationCenter]
     removeObserver:self
     name:NSWorkspaceWillSleepNotification object:NULL];
    
    [[[NSWorkspace sharedWorkspace] notificationCenter]
     removeObserver:self
     name:NSWorkspaceDidWakeNotification object:NULL];
}
@end
