//
//  NativeMenu+Raspberry.m
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "NativeMenu+Raspberry.h"
#import "RaspberryManager.h"

@interface NativeMenu(RaspberryPrivate)
-(void)raspberryRefreshingStarted:(NSNotification *)aNotification;
-(void)raspberryRefreshingEnded:(NSNotification *)aNotification;
-(void)raspberryUpdateRunningNodeCount:(NSNotification *)aNotification;
-(void)raspberryUpdateNodeCount:(NSNotification *)aNotification;
-(void)raspberryNodeUp:(NSNotification *)aNotification;
-(void)raspberryNodeDown:(NSNotification *)aNotification;
-(void)raspberryNodeAdded:(NSNotification *)aNotification;
-(void)raspberryNodeRemoved:(NSNotification *)aNotification;
-(void)raspberryNodeUpdated:(NSNotification *)aNotification;
@end

@implementation NativeMenu(Raspberry)

-(void)raspberryRefreshingStarted:(NSNotification *)aNotification {
}

-(void)raspberryRefreshingEnded:(NSNotification *)aNotification {
}

-(void)raspberryUpdateRunningNodeCount:(NSNotification *)aNotification {
}

-(void)raspberryUpdateNodeCount:(NSNotification *)aNotification {
}

-(void)raspberryNodeUp:(NSNotification *)aNotification {
}

-(void)raspberryNodeDown:(NSNotification *)aNotification {
}

-(void)raspberryNodeAdded:(NSNotification *)aNotification {
}

-(void)raspberryNodeRemoved:(NSNotification *)aNotification {
}

-(void)raspberryNodeUpdated:(NSNotification *)aNotification {
}

-(void)raspberryRegisterNotifications {
    
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryRefreshingStarted:)      name:kRASPBERRY_MANAGER_REFRESHING_STARTED          object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryRefreshingEnded:)        name:kRASPBERRY_MANAGER_REFRESHING_ENDED            object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryUpdateRunningNodeCount:) name:kRASPBERRY_MANAGER_UPDATE_RUNNING_NODE_COUNT   object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryUpdateNodeCount:)        name:kRASPBERRY_MANAGER_UPDATE_NODE_COUNT           object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryNodeUp:)                 name:kRASPBERRY_MANAGER_NODE_UP                     object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryNodeDown:)               name:kRASPBERRY_MANAGER_NODE_DOWN                   object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryNodeAdded:)              name:kRASPBERRY_MANAGER_NODE_ADDED                  object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryNodeRemoved:)            name:kRASPBERRY_MANAGER_NODE_REMOVED                object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(raspberryNodeUpdated:)            name:kRASPBERRY_MANAGER_NODE_UPDATED                object:nil];
    
}

@end
