//
//  UserCheckVC.m
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "UserCheckVC.h"

@interface UserCheckVC ()
-(void)_enableControls;
-(void)_disableControls;
@end

@implementation UserCheckVC

-(void) finishConstruction {
    [super finishConstruction];
    [self setTitle:@"Invitation Check"];
}

- (void)viewDidLoad {
    [super viewDidLoad];

    [[((BaseBrandView *)self.view) contentBox] removeFromSuperview];
    ((BaseBrandView *)self.view).contentBox = nil;
}

- (void)viewDidAppear {
    [super viewDidAppear];

    [self.fieldEmail becomeFirstResponder];
    [self.fieldEmail selectText:self];
    [[self.fieldEmail currentEditor] setSelectedRange:NSMakeRange([[self.fieldEmail stringValue] length], 0)];
}

-(IBAction)check:(id)sender {
    [self.stageControl shouldControlProgressFrom:self withParam:nil];
    
//    [self _disableControls];
}

-(IBAction)cancel:(id)sender {
    [self.stageControl shouldControlRevertFrom:self withParam:nil];
}

#pragma mark - StageStep
-(void)didControl:(NSObject<StepControl> *)aControl progressFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult {
    if (aStep != self) {
        Log(@"this is not identical step!");
        return;
    }
}

-(void)didControl:(NSObject<StepControl> *)aControl revertFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult {
    if (aStep != self) {
        Log(@"this is not identical step!");
        return;
    }
}

#pragma mark - Controls
-(void)_enableControls {
    [self.progress setHidden:YES];
    [self.progress stopAnimation:nil];
    [self.btnCancel setEnabled:YES];
    [self.btnCheck setEnabled:YES];
}

-(void)_disableControls {
    [self.progress setHidden:NO];
    [self.progress startAnimation:nil];
    [self.btnCancel setEnabled:NO];
    [self.btnCheck setEnabled:NO];
}

@end
