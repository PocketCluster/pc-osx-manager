//
//  AppDelegate+Window.h
//  manager
//
//  Created by Almighty Kim on 8/12/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#include "AppDelegate.h"
#import "UpdateProgress.h"

@interface AppDelegate(Window) <UpdateProgress>
- (NSObject *) activeWindowByClassName:(NSString *)aClassName withResponder:(id)aResponder;
- (void)addOpenWindow:(BaseWindowController *)window;
- (void)removeOpenWindow:(id)window;
@end
