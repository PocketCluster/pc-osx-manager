//
//  DPSetupWC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "DPSetupWC.h"

#import "DPSetupWindow.h"
#import "PCSetup1VC.h"
#import "PCSetup2VVVC.h"
#import "PCSetup2RPVC.h"
#import "PCSetup3VC.h"
#import "PCSetup4VC.h"

#import "Util.h"

@interface DPSetupWC ()
@end

@implementation DPSetupWC

- (void)windowDidLoad {
    [super windowDidLoad];
    
    DPSetupWindow *setupFlow = (DPSetupWindow *)[self window];
    [setupFlow setTitle:@"New Cluster"];

    if (setupFlow == nil){
        return;
    }
    
    NSViewController *vc1 =
        [[PCSetup1VC alloc]
         initWithNibName:@"PCSetup1VC"
         bundle:[NSBundle mainBundle]];
    
    NSViewController *vc2v =
        [[PCSetup2VVVC alloc]
         initWithNibName:@"PCSetup2VVVC"
         bundle:[NSBundle mainBundle]];
    
    NSViewController *vc2r =
        [[PCSetup2RPVC alloc]
         initWithNibName:@"PCSetup2RPVC"
         bundle:[NSBundle mainBundle]];
    
    NSViewController *vc3 =
        [[PCSetup3VC alloc]
         initWithNibName:@"PCSetup3VC"
         bundle:[NSBundle mainBundle]];
    
    NSViewController *vc4 =
        [[PCSetup4VC alloc]
         initWithNibName:@"PCSetup4VC"
         bundle:[NSBundle mainBundle]];

    [setupFlow
     initWithViewControllers:@[vc1,
                               //vc2v,
                               vc2r,
                               vc3,
                               vc4]
     completionHandler:^(BOOL completed) {
         if (!completed) {
             NSLog(@"Cancelled setup process");
         } else {
             NSLog(@"Completed setup process");
         }
         [setupFlow orderOut:self];
     }];
    
    [setupFlow setBackgroundImage:[NSImage imageNamed:@"AppIcon"]];
    
    [setupFlow resetToZeroStage];
    [setupFlow makeKeyAndOrderFront:self];
    
    [[Util getApp] startMulticastSocket];
}

@end
