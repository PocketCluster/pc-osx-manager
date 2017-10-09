//
//  PCSetup1VC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup1VC.h"

@implementation PCSetup1VC

- (void) finishConstruction {
    [super finishConstruction];
    [self setTitle:@"Setup a New Cluster"];
}

- (void) viewDidLoad {
    [super viewDidLoad];

    [[((BaseBrandView *)self.view) contentBox] addSubview:self.pannel];
    self.pannel = nil;
}

- (IBAction)setup:(id)sender {
    [self.stageControl shouldControlProgressFrom:self withParam:nil];
}

- (IBAction)cancel:(id)sender {
    [self.stageControl shouldControlRevertFrom:self withParam:nil];
}

@end
