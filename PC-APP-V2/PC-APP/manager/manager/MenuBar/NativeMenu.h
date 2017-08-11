//
//  NativeMenu.h
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#import "AboutWindow.h"

@interface NativeMenu : NSObject <NSMenuDelegate>{
@private
    AboutWindow      *_aboutWindow;
    NSStatusItem     *_statusItem;
}
@property (nonatomic, strong) AboutWindow *aboutWindow;
@property (nonatomic, strong) NSStatusItem *statusItem;

- (void) addCommonMenu:(NSMenu *)menuRoot;
@end
