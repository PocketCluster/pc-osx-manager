//
//  AppDelegate+Window.m
//  manager
//
//  Created by Almighty Kim on 8/12/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+Window.h"

@interface AppDelegate(WindowPrivate)
- (void)updateProcessType;
@end

@implementation AppDelegate(Window)

- (void)addOpenWindow:(id)window {
    @synchronized(_openWindows) {
        [_openWindows addObject:window];
        [self updateProcessType];
    }
}

- (void)removeOpenWindow:(id)window {
    @synchronized(_openWindows) {
        [_openWindows removeObject:window];
        [self updateProcessType];
    }
}

- (void)updateProcessType {
    if([_openWindows count] == 0) {
        ProcessSerialNumber psn = { 0, kCurrentProcess };
        TransformProcessType(&psn, kProcessTransformToBackgroundApplication);
    } else {
        ProcessSerialNumber psn = { 0, kCurrentProcess };
        TransformProcessType(&psn, kProcessTransformToForegroundApplication);
#if USE_DEPRECATED_CF_METHOD
        SetFrontProcess(&psn);
#else
        //[NSApp activateIgnoringOtherApps:YES];
        [[NSApplication sharedApplication] activateIgnoringOtherApps: YES];
#endif
    }
}

@end
