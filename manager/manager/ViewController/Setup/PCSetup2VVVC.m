//
//  PCSetup2VVVC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup2VVVC.h"
#import "TaskOutputWindow.h"
#import "Util.h"

@interface PCSetup2VVVC ()
@property (nonatomic, strong) TaskOutputWindow *task;
@end

@implementation PCSetup2VVVC

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do view setup here.
}

-(IBAction)vagrantUp:(id)sender
{
    
    TaskOutputWindow *task = [[TaskOutputWindow alloc] initWithWindowNibName:@"TaskOutputWindow"];

    task.taskCommand = @"ls /Users/almightykim/";
//    task.target = machine;
//    task.taskAction = command;

    task.isSudoCommand = YES;
    
    [[NSApplication sharedApplication]
     beginSheet:[task window]
     modalForWindow:[self.view window]
     modalDelegate:nil
     didEndSelector:NULL
     contextInfo:NULL];
    
    [[Util getApp] addOpenWindow:task];
    
    
}

@end
