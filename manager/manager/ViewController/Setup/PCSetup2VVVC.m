//
//  PCSetup2VVVC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup2VVVC.h"

#import "Util.h"
#import "PCTaskOutputWindow.h"

@interface PCSetup2VVVC ()

@end

@implementation PCSetup2VVVC

-(IBAction)vagrantUp:(id)sender
{
    
    NSString *basePath    = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
    NSString *setupScript = [NSString stringWithFormat:@"%@/setup/setup_vagrant_base.sh",basePath];
    
    __block PCTaskOutputWindow *tow = [[PCTaskOutputWindow alloc] initWithNibName:nil bundle:nil];
    tow.taskCommand = [NSString stringWithFormat:@"sh %@ %@", setupScript, basePath];
    tow.sudoCommand = YES;
    
    WEAK_SELF(self);
    [belf.view.window
     beginSheet:tow
     completionHandler:^(NSModalResponse returnCode) {
         [belf.view.window endSheet:tow];
     }];
    
    [tow launchTask];
    
    //[[Util getApp] addOpenWindow:tow];
}

@end
