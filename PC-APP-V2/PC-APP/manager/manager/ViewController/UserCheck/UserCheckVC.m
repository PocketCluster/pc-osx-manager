//
//  UserCheckVC.m
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "ShowAlert.h"
#import "UserCheckVC.h"
#import "PCRouter.h"
#import "NullStringChecker.h"
#import "NSString+EmailForm.h"

@interface UserCheckVC ()
-(void)_disableControls;
@end

@implementation UserCheckVC

-(void) finishConstruction {
    [super finishConstruction];
    [self setTitle:@"PocketCluster Invitation"];
}

- (void)viewDidLoad {
    [super viewDidLoad];

    [[((BaseBrandView *)self.view) contentBox] addSubview:self.pannel];
    self.pannel = nil;
}

- (void)viewDidAppear {
    [super viewDidAppear];

    [self.fieldEmail becomeFirstResponder];
    [self.fieldEmail selectText:self];
    [[self.fieldEmail currentEditor] setSelectedRange:NSMakeRange([[self.fieldEmail stringValue] length], 0)];
}

-(IBAction)check:(id)sender {
    NSString *email = [self.fieldEmail stringValue];
    NSString *code =  [self.fieldCode stringValue];

    if (ISNULL_STRING(email) || ![email isValidEmailForm]) {
        [ShowAlert
         showWarningAlertWithTitle:@"Invalid Email Address"
         message:@"Please provide valid email address"];
        return;
    }
    if (ISNULL_STRING(code)) {
        [ShowAlert
         showWarningAlertWithTitle:@"Invalid Invitation"
         message:@"Please provide valid invitation code"];
        return;
    }

    [self _disableControls];

    [PCRouter
     routeRequestPost:RPATH_USER_AUTHED
     withRequestBody:@{@"email":email, @"code":code}];
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
-(void) enableControls {
    [self.circularProgress setHidden:YES];
    [self.circularProgress stopAnimation:nil];
    [self.circularProgress displayIfNeeded];
    [self.circularProgress removeFromSuperview];
    [self setCircularProgress:nil];

    [self.fieldEmail setEnabled:YES];
    [self.fieldCode setEnabled:YES];
    [self.btnCancel setEnabled:YES];
    [self.btnCheck setEnabled:YES];
}

-(void) _disableControls {
    [self.fieldEmail setEnabled:NO];
    [self.fieldCode setEnabled:NO];
    [self.btnCancel setEnabled:NO];
    [self.btnCheck setEnabled:NO];

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
