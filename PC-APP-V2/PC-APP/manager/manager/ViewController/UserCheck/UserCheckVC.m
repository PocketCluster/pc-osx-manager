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

@interface UserCheckVC ()
-(void)_disableControls;
@end

static bool
is_valid_invitation(NSString* invitation) {
    static NSString * const check_pattern = @"^[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}$";
    if ISNULL_STRING(invitation) {
        return false;
    }
    return [[NSPredicate
             predicateWithFormat:@"SELF MATCHES %@", check_pattern]
            evaluateWithObject:invitation];
}

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

    [self.view.window setInitialFirstResponder:[self fieldCode]];
    [self.fieldCode becomeFirstResponder];
    [self.fieldCode selectText:self];
    [[self.fieldCode currentEditor] setSelectedRange:NSMakeRange([[self.fieldCode stringValue] length], 0)];
}

-(IBAction)check:(id)sender {
    NSString *invitation =  [self.fieldCode stringValue];
    if (!is_valid_invitation(invitation)) {
        [ShowAlert
         showWarningAlertWithTitle:@"Invalid Code"
         message:@"Please provide valid invitation code"];
        return;
    }

    [self _disableControls];

    [PCRouter
     routeRequestPost:RPATH_USER_AUTHED
     withRequestBody:@{@"invitation":invitation}];
}

-(IBAction) cancel:(id)sender {
    [self.stageControl shouldControlRevertFrom:self withParam:nil];
}

#pragma mark - StageStep
-(void)didControl:(NSObject<StepControl> *)aControl progressedFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult {
    if (aStep != self) {
        Log(@"this is not identical step!");
        return;
    }
}

-(void)didControl:(NSObject<StepControl> *)aControl revertedFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult {
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

    [self.fieldCode setEnabled:YES];
    [self.btnCancel setEnabled:YES];
    [self.btnCheck setEnabled:YES];
}

-(void) _disableControls {
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
