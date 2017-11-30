//
//  PCSetup2VC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCConstants.h"
#import "StatusCache.h"

#import "ShowAlert.h"
#import "PCRouter.h"
#import "PCSetup2VC.h"

static NSString * const kNameColTag = @"nameCol";
static NSString * const kAddrColTag = @"addrCol";
NSString * const kClusterSetupResult = @"SETUP_SUCCESS";

@interface PCSetup2VC ()<PCRouteRequest>
@property (nonatomic, strong) NSArray *nodeList;
- (void)_enableControls;
- (void)_disableControls;
@end

@implementation PCSetup2VC

- (void) finishConstruction {
    [super finishConstruction];
    [self setTitle:@"Build Cluster"];
    
    WEAK_SELF(self);

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NODE_REG_START)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         if (![[response valueForKeyPath:@"node-reg-start.status"] boolValue]) {
             [ShowAlert
              showWarningAlertWithTitle:@"Unable to add new node"
              message:[response valueForKeyPath:@"node-reg-start.error"]];
         }
     }];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NODE_UNREG_LIST)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {
         NSArray<NSDictionary *>* list = [response valueForKeyPath:@"node-unreged.unreged-list"];
         if (belf != nil) {
             [belf setNodeList:list];
             [[belf nodeTable] reloadData];
         }
     }];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NODE_REG_CANDIDATE)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"path %@ response %@", path, response);

         if (![[response valueForKeyPath:@"node-reg-candidate.status"] boolValue]) {
             [[StatusCache SharedStatusCache] setClusterSetup:NO];

             if (belf != nil) {
                 [belf _enableControls];
             }

             [ShowAlert
              showWarningAlertWithTitle:@"Unable to add new node"
              message:[response valueForKeyPath:@"node-reg-candidate.error"]];
         }
     }];

    [[PCRouter sharedRouter]
     addGetRequest:self
     onPath:@(RPATH_NODE_REG_CONFIRM)
     withHandler:^(NSString *method, NSString *path, NSDictionary *response) {

         Log(@"path %@ response %@", path, response);

         [[StatusCache SharedStatusCache] setClusterSetup:NO];

         if (![[response valueForKeyPath:@"node-reg-confirm.status"] boolValue]) {
             [ShowAlert
              showAlertMessageWithTitle:@"Unable to add new node"
              message:[response valueForKeyPath:@"node-reg-confirm.error"]
              action:^(NSInteger result) {
                  if (belf != nil) {
                      [belf.stageControl shouldControlProgressFrom:belf withParam:@{kClusterSetupResult:@NO}];
                      [belf _enableControls];
                  }
              }];

         } else {
             [ShowAlert
              showAlertMessageWithTitle:@"Setup Successful!"
              message:@"Cluster is successfully setup."
              action:^(NSInteger result) {
                  if (belf != nil) {
                      [belf.stageControl shouldControlProgressFrom:belf withParam:@{kClusterSetupResult:@TRUE}];
                  }
              }];
         }
     }];
}

- (void) prepareDestruction {
    [super prepareDestruction];

    [[PCRouter sharedRouter] delGetRequest:self onPath:@(RPATH_NODE_REG_START)];
    [[PCRouter sharedRouter] delGetRequest:self onPath:@(RPATH_NODE_UNREG_LIST)];
    [[PCRouter sharedRouter] delGetRequest:self onPath:@(RPATH_NODE_REG_CANDIDATE)];
    [[PCRouter sharedRouter] delGetRequest:self onPath:@(RPATH_NODE_REG_CONFIRM)];
}

- (void) viewDidLoad {
    [super viewDidLoad];
    
    [[((BaseBrandView *)self.view) contentBox] addSubview:self.pannel];
    self.pannel = nil;
}

#pragma mark - NSWindowDelegate
- (BOOL)windowShouldClose:(NSWindow *)sender {
    if ([[StatusCache SharedStatusCache] isClusterSetup]) {
        [ShowAlert
         showWarningAlertWithTitle:@"Please do not close Window"
         message:@"Please wait until cluster setup is completed"];
        return NO;
    }
    [PCRouter routeRequestGet:RPATH_NODE_REG_STOP];
    return YES;
}

#pragma mark - NSTableViewDataSourceDelegate
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView {
    return [self.nodeList count];
}

- (nullable id)tableView:(NSTableView *)aTableView objectValueForTableColumn:(nullable NSTableColumn *)aTableColumn row:(NSInteger)row {
    if (aTableView == nil) {
        return nil;
    }
    NSDictionary *node = [self.nodeList objectAtIndex:row];
    if ([[aTableColumn identifier] isEqualToString:kNameColTag]) {
        return [node valueForKey:@"name"];
    }
    if ([[aTableColumn identifier] isEqualToString:kAddrColTag]) {
        return [node valueForKey:@"addr"];
    }
    return nil;
}

#pragma mark - NSTableViewDelegate
- (BOOL)selectionShouldChangeInTableView:(NSTableView *)tableView {
    return NO;
}

- (BOOL)tableView:(NSTableView *)tableView shouldSelectRow:(NSInteger)row {
    return NO;
}

#pragma mark - StageStep
- (void) control:(NSObject<StepControl> *)aControl askedProgressWithParam:(NSDictionary *)aParam {
    // start registration rasker
    [PCRouter routeRequestGet:RPATH_NODE_REG_START];
}
- (void) didControl:(NSObject<StepControl> *)aControl progressedFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult {
}

- (void) control:(NSObject<StepControl> *)aControl askedRevertWithParam:(NSDictionary *)aParam {
}
- (void) didControl:(NSObject<StepControl> *)aControl revertedFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult {
}

#pragma mark - UI Control
- (void)_enableControls {
    [self.btnBuild setEnabled:YES];
    [self.btnCancel setEnabled:YES];

    [self.circularProgress setHidden:YES];
    [self.circularProgress stopAnimation:nil];
    [self.circularProgress displayIfNeeded];
    [self.circularProgress removeFromSuperview];
    [self setCircularProgress:nil];
}

- (void)_disableControls {
    [self.btnBuild setEnabled:NO];
    [self.btnCancel setEnabled:NO];

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

#pragma mark - IBACTION
-(IBAction)build:(id)sender {
    // return if there is no node
    if ([self.nodeList count] == 0){
        return;
    }
    [self _disableControls];

    // enter into cluster setup status
    [[StatusCache SharedStatusCache] setClusterSetup:YES];

    // send registration signal to engine
    [PCRouter routeRequestGet:RPATH_NODE_REG_CANDIDATE];
}

-(IBAction)cancel:(id)sender {
    [PCRouter routeRequestGet:RPATH_NODE_REG_STOP];
    [self.stageControl shouldControlRevertFrom:self withParam:nil];
}

@end
