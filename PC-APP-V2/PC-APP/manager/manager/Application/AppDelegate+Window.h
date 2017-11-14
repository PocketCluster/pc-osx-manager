//
//  AppDelegate+Window.h
//  manager
//
//  Created by Almighty Kim on 8/12/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#include "AppDelegate.h"

@interface AppDelegate(Window)
- (NSObject *) activeWindowByClassName:(NSString *)aClassName withResponder:(id)aResponder;
- (BaseWindowController *)findWindowControllerByClassName:(NSString *)aClassName withResponder:(id)aResponder;
- (void)addOpenWindow:(BaseWindowController *)window;
- (void)removeOpenWindow:(id)window;
- (void)updateProcessType;
@end
