//
//  PCSetup2VVVC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup2VVVC.h"

@interface PCSetup2VVVC ()

@end

@implementation PCSetup2VVVC

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do view setup here.
}

-(IBAction)vagrantUp:(id)sender
{
//    [[self setupFlow] resetToZeroStage];
    [[NSApplication sharedApplication]
     beginSheet:nil//[self setupFlow]
     modalForWindow:[[self view] window]
     modalDelegate:nil
     didEndSelector:NULL
     contextInfo:NULL];

}

@end
