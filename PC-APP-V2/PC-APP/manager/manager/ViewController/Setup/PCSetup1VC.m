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

    [self _disableControls];
}

- (IBAction)setup:(id)sender {
    [self.stageControl shouldControlProgressFrom:self withParam:nil];
}

- (IBAction)cancel:(id)sender {
    [self.stageControl shouldControlRevertFrom:self withParam:nil];
}

-(void) enableControls {
    [self.circularProgress setHidden:YES];
    [self.circularProgress stopAnimation:nil];
    [self.circularProgress displayIfNeeded];
    [self.circularProgress removeFromSuperview];
    [self setCircularProgress:nil];

    [self.btnCancel setEnabled:YES];
    [self.btnSetup setEnabled:YES];
}

-(void) _disableControls {
    [self.btnCancel setEnabled:NO];
    [self.btnSetup setEnabled:NO];

    NSProgressIndicator *ind = [[NSProgressIndicator alloc] initWithFrame:(NSRect){{20.0, 20.0}, {16.0, 16.0}}];
    [ind setControlSize:NSSmallControlSize];
    [ind setStyle:NSProgressIndicatorSpinningStyle];
    [self.view addSubview:ind];
    [ind setHidden:NO];
    [ind setIndeterminate:YES];
    [ind startAnimation:self];
    [ind displayIfNeeded];

    [self setCircularProgress:ind];
}

@end
