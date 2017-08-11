//
//  NativeMenu.h
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

@interface NativeMenu : NSObject <NSMenuDelegate>{
@private
    NSStatusItem     *_statusItem;
}
@property (nonatomic, strong) NSStatusItem *statusItem;

- (void) addCommonMenu:(NSMenu *)menuRoot;
@end
