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
}
@property (nonatomic, strong, readonly) NSStatusItem *statusItem;

#pragma mark - update notification
- (void) updateNewVersionAvailability:(BOOL)IsAvailable;

#pragma mark - Common menus
- (void) setupCheckupMenu;
- (void) setupOperationMenu;
@end
