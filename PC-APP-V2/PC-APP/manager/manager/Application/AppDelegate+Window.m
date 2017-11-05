//
//  AppDelegate+Window.m
//  manager
//
//  Created by Almighty Kim on 8/12/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+Window.h"
#import "BaseWindowController.h"
#import "NullStringChecker.h"
#include "pc-core.h"

@interface AppDelegate(WindowPrivate)
@end

@implementation AppDelegate(Window)

// add & open new window if there is no same class.
// bring forward and activate if there is a window with same class
- (NSObject *) activeWindowByClassName:(NSString *)aClassName withResponder:(id)aResponder {

    if (ISNULL_STRING(aClassName)) {
        Log(@"cannot find class from invalid string");
        return nil;
    }
    
    // if this isn't main thread and gets called, re-invoke it in main thread.
    if (![NSThread isMainThread]) {
        WEAK_SELF(self);
        [[NSOperationQueue mainQueue]
         addOperationWithBlock:^{
             if(belf){
                 [belf activeWindowByClassName:aClassName withResponder:aResponder];
             }
         }];
        return nil;
    }

    @synchronized(_openWindows) {

        for (BaseWindowController* window in _openWindows) {

            if ([[window class] isSubclassOfClass:[BaseWindowController class]] &&\
                [[window className] isEqualToString:aClassName]) {
                Log(@"found an obj by the class name %@", aClassName);
                [[NSApplication sharedApplication] activateIgnoringOtherApps:YES];
                [window showWindow:aResponder];
                [window bringToFront];
                return window;
            }
        }

        Class class = NSClassFromString(aClassName);
        if (class == nil) {
            Log(@"Unable to find a class for the given name %@", aClassName);
            return nil;
        }
        if (![class isSubclassOfClass:[BaseWindowController class]]) {
            Log(@"Unable to spwan a non-child class of BaseWindowController for the given name %@", aClassName);
            return nil;
        }

        BaseWindowController* window = [class alloc];
        if (window == nil) {
            Log(@"Unable to spawn a class for the given name %@", aClassName);
            return nil;
        }

        Log(@"spawned an obj by the class name %@", aClassName);
        
        window = [window initWithWindowNibName:aClassName];
        [[NSApplication sharedApplication] activateIgnoringOtherApps:YES];
        [window showWindow:aResponder];
        [window bringToFront];

        // add window to managed list
        [_openWindows addObject:window];
        [self updateProcessType];

        return window;
    }
}

- (void)addOpenWindow:(id)window {

    // if this isn't main thread and gets called, re-invoke it in main thread.
    if (![NSThread isMainThread]) {
        WEAK_SELF(self);
        [[NSOperationQueue mainQueue]
         addOperationWithBlock:^{
             if(belf){
                 [belf addOpenWindow:window];
             }
         }];
        return;
    }

    if (![[window class] isSubclassOfClass:[BaseWindowController class]]) {
        Log(@"Unable to add a non-child class of BaseWindowController: %@", [window className]);
        return;
    }

    @synchronized(_openWindows) {
        [_openWindows addObject:window];
        [self updateProcessType];
    }
}

- (void)removeOpenWindow:(id)window {

    // if this isn't main thread and gets called, re-invoke it in main thread.
    if (![NSThread isMainThread]) {
        WEAK_SELF(self);
        [[NSOperationQueue mainQueue]
         addOperationWithBlock:^{
             if(belf){
                 [belf removeOpenWindow:window];
             }
         }];
        return;
    }
    
    @synchronized(_openWindows) {
        [_openWindows removeObject:window];
        [self updateProcessType];
    }
}

/*
 * lifecycleAlive(); should have been called in 'windowWillClose' of 
 * BaseWindowController or AppDelegate. But, we need to handle more windows 
 * than 1, and it is most appropriate to place the call here in updateProcessType
 *
 * - (void)windowWillClose:(NSNotification *)notification {
 *     Log(@"%s", __PRETTY_FUNCTION__);
 *     lifecycleAlive();
 * }
 */
- (void)updateProcessType {
    if([_openWindows count] == 0) {
        lifecycleAlive();
        ProcessSerialNumber psn = { 0, kCurrentProcess };
        TransformProcessType(&psn, kProcessTransformToBackgroundApplication);
    } else {
        ProcessSerialNumber psn = { 0, kCurrentProcess };
        TransformProcessType(&psn, kProcessTransformToForegroundApplication);
#if 0
        SetFrontProcess(&psn);
#else
        //[NSApp activateIgnoringOtherApps:YES];
        [[NSApplication sharedApplication] activateIgnoringOtherApps: YES];
#endif
    }
}


@end
