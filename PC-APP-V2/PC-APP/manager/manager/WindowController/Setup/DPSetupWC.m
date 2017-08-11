//
//  DPSetupWC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+Window.h"

#import "DPSetupWindow.h"
#import "DPSetupWC.h"

#import "PCSetup1VC.h"
#import "PCSetup2RPVC.h"

//#import "RaspberryManager.h"

@implementation DPSetupWC

- (void)windowDidLoad {
    [super windowDidLoad];
    
    WEAK_SELF(self);
    
    DPSetupWindow *setupFlow = (DPSetupWindow *)[self window];
    [setupFlow setTitle:@"New Cluster"];

    if (setupFlow == nil || ![setupFlow isKindOfClass:[DPSetupWindow class]]){
        return;
    }
    
    NSViewController *vc1 =
        [[PCSetup1VC alloc]
         initWithNibName:@"PCSetup1VC"
         bundle:[NSBundle mainBundle]];
    
    [setupFlow
     initWithViewControllers:@[vc1]
     completionHandler:^(BOOL completed) {

#if 0
         if (!completed) {
             Log(@"Cancelled setup process");
         } else {
             Log(@"Completed setup process");
         }
#endif

         [setupFlow orderOut:belf];
         [belf close];
     }];

    [setupFlow setBackgroundImage:[NSImage imageNamed:@"AppIcon"]];
    
    [setupFlow resetToZeroStage];
    [setupFlow makeKeyAndOrderFront:self];

}

-(void)resetSetupStage {
    [(DPSetupWindow *)self.window resetToZeroStage];
}

-(void)windowWillClose:(NSNotification *)notification {

    DPSetupWindow *dsw = (DPSetupWindow *)self.window;
    NSArray *dpwvc = [NSArray arrayWithArray:[dsw viewControllers]];
    for(NSViewController *vc in dpwvc){
        if([vc isKindOfClass:[PCSetup2RPVC class]]){
            //[[RaspberryManager sharedManager] removeAgentDelegateFromQueue:vc];
        }
    }
    [dsw removeNotifications];
    [super windowWillClose:notification];
}

@end
