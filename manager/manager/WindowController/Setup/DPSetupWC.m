//
//  DPSetupWC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "DPSetupWindow.h"
#import "DPSetupWC.h"
#import "PCSetup1VC.h"
#import "Util.h"

@implementation DPSetupWC

- (void)windowDidLoad {
    [super windowDidLoad];
    
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
         if (!completed) {
             Log(@"Cancelled setup process");
         } else {
             Log(@"Completed setup process");
         }
         [setupFlow orderOut:self];
     }];
    
    [setupFlow setBackgroundImage:[NSImage imageNamed:@"AppIcon"]];
    
    [setupFlow resetToZeroStage];
    [setupFlow makeKeyAndOrderFront:self];

}

-(void)resetSetupStage {
    [(DPSetupWindow *)self.window resetToZeroStage];
}

-(void)bringToFront {
    [self.window makeKeyAndOrderFront:[Util getApp]];
}
@end
