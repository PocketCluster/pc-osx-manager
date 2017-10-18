//
//  NativeMenu.h
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#import "PCRouteRequest.h"

@interface NativeMenu : NSObject <NSMenuDelegate, PCRouteRequest>{
@private
    NSStatusItem     *_statusItem;
    NSMenuItem       *_updateAvail;
}
@property (nonatomic, strong, readonly) NSStatusItem *statusItem;
@property (nonatomic, strong, readonly) NSMenuItem *updateAvail;

- (void) clusterStatusOn;
- (void) clusterStatusOff;

- (void) updateNewVersionAvailability:(BOOL)IsAvailable;

#pragma mark - menu updator
- (void) updateMenuWithCondition;

#pragma mark - Common menus
- (void) addCommonMenu:(NSMenu *)menuRoot;
- (void) addInitCommonMenu:(NSMenu *)menuRoot;
@end
