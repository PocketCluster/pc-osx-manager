//
//  AppDelegate.h
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

@class NativeMenu;

#import "PCConstants.h"

@interface AppDelegate : NSObject <NSApplicationDelegate>
@property (readonly) NativeMenu *nativeMenu;
- (void)addOpenWindow:(id)window;
- (void)removeOpenWindow:(id)window;
@end

