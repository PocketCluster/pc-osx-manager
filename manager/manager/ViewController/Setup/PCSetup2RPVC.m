//
//  PCSetup2RPVC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup2RPVC.h"

#import "PCSetup3VC.h"
#import "RaspberryManager.h"
#import "PCProcManager.h"
#import "PCTask.h"
#import "Util.h"
#import "PCConstants.h"

@interface PCSetup2RPVC ()<PCTaskDelegate, RaspberryAgentDelegate>
@property (atomic, strong) NSMutableArray *nodeList;
@property (strong, nonatomic) NSDictionary *progDict;

@property (strong, nonatomic) PCTask *sudoTask;
@property (strong, nonatomic) PCTask *saltTask;
@property (strong, nonatomic) PCTask *userTask;
@property (strong, nonatomic) PCTask *skeyTask;
@property (nonatomic, strong) NSString *hostName;

@property (readwrite, nonatomic) BOOL canContinue;
@property (readwrite, nonatomic) BOOL canGoBack;

- (void)setUIToProceedState;
- (void)resetUIForFailure;
- (void)setToNextStage;

- (void)removeViewControler;
@end

@implementation PCSetup2RPVC
@synthesize canContinue;
@synthesize canGoBack;

-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {

    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    
    if(self){
        self.nodeList = [NSMutableArray arrayWithCapacity:0];
        
        self.progDict = @{@"SUDO_SETUP_STEP_0":@[@"Base config done...",@10.0]
                          ,@"SUDO_SETUP_DONE":@[@"Start setting up Vagrant",@20.0]
                          ,@"USER_SETUP_STEP_0":@[@"USER_SETUP_STEP_0",@30.0]
                          ,@"USER_SETUP_STEP_1":@[@"USER_SETUP_STEP_1",@50.0]
                          ,@"USER_SETUP_STEP_2":@[@"USER_SETUP_STEP_2",@90.0]
                          ,@"USER_SETUP_DONE":@[@"USER_SETUP_DONE",@100.0]};
        
        self.hostName = [[[NSHost currentHost] localizedName] lowercaseString];
        
        [self resetToInitialState];
        
        [[RaspberryManager sharedManager] addAgentDelegateToQueue:self];
        [[RaspberryManager sharedManager] refreshInterface];
    }
    
    return self;
}

-(void)viewDidAppear {
    if ([[RaspberryManager sharedManager] ethernetInterface]){
        [self.warningLabel setHidden:YES];
    }else{
        [self.warningLabel setHidden:NO];
    }
}

-(void)dealloc {
    Log(@"%s",__PRETTY_FUNCTION__);
}

#pragma mark - GCDAsyncUdpSocketDelegate

- (void)didReceiveUnboundedAgentData:(NSDictionary *)anAgentData {
    
    BOOL doesNodeExist = false;
    for (NSDictionary *node in self.nodeList){
        if([[node valueForKey:SLAVE_NODE_MACADDR] isEqualToString:[anAgentData valueForKey:SLAVE_NODE_MACADDR]]){
            doesNodeExist = true;
            break;
        }
    }

    if (!doesNodeExist && self.nodeList.count < MAX_TRIAL_RASP_NODE_COUNT){
        [self.nodeList addObject:anAgentData];
        [self.nodeList sortUsingComparator:^NSComparisonResult(NSDictionary*  _Nonnull node1, NSDictionary* _Nonnull node2) {
            return [[node1 valueForKey:ADDRESS] compare:[node2 valueForKey:ADDRESS] options:NSNumericSearch];
        }];

        for (int i = 0; i < [self.nodeList count]; i++){
            NSMutableDictionary *nd = [self.nodeList objectAtIndex:i];
            [nd setValue:[NSString stringWithFormat:@"pc-node%d",(i + 1)] forKey:SLAVE_NODE_NAME];
        }

        [self.nodeTable reloadData];
    }
}


#pragma mark - NSTableViewDataSourceDelegate
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView {
    return [self.nodeList count];
}

- (nullable id)tableView:(NSTableView *)tableView objectValueForTableColumn:(nullable NSTableColumn *)tableColumn row:(NSInteger)row {
    return [self.nodeList objectAtIndex:row];
}

#pragma mark - NSTableViewDelegate
-(NSView *)tableView:(NSTableView *)aTableView viewForTableColumn:(NSTableColumn *)aTableColumn row:(NSInteger)row{

    NSDictionary *nd = [self.nodeList objectAtIndex:row];
    NSTableCellView *nv = [aTableView makeViewWithIdentifier:@"nodeview" owner:self];
    
    if([aTableColumn.identifier isEqualToString:@"nodename"]){
        [nv.textField setStringValue:[nd valueForKey:SLAVE_NODE_NAME]];
    }else{
        [nv.textField setStringValue:[nd valueForKey:ADDRESS]];
    }
    
    return nv;
}

- (BOOL)selectionShouldChangeInTableView:(NSTableView *)tableView {
    return NO;
}

- (BOOL)tableView:(NSTableView *)tableView shouldSelectRow:(NSInteger)row {
    return NO;
}

#pragma mark - PCTaskDelegate
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
    if(aTask.terminationStatus != 0) {
        
        Log(@"installation error ! %d",aTask.terminationStatus);
        
        [self resetUIForFailure];
        [self.progressLabel setStringValue:@"Installation Error. Please try again."];
        
        self.sudoTask = nil;
        self.saltTask = nil;
        self.userTask = nil;
        return;
    }
    
    [self setUIToProceedState];

    if(self.sudoTask == aPCTask ){
        
        PCTask *st = [PCTask new];
        st.taskCommand = @"brew install saltstack 2>&1";
        st.delegate = self;
        self.saltTask = st;
        [st launchTask];
        
        self.sudoTask = nil;
    }
    
    if(self.saltTask == aPCTask) {
        
        [[PCProcManager sharedManager] freshSaltStart];
        
        NSString *basePath  = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
        NSString *userSetup = [NSString stringWithFormat:@"%@/setup/raspberry_user_setup.sh",basePath];
        
        PCTask *userTask = [PCTask new];
        userTask.taskCommand = [NSString stringWithFormat:@"bash %@ %@ %ld", userSetup, basePath, [self.nodeList count]];
        userTask.delegate = self;
        self.userTask = userTask;
        [userTask launchTask];
        
        self.saltTask = nil;
    }
    
    if(self.userTask == aPCTask) {

        PCTask *kt = [PCTask new];
        kt.taskCommand = @"salt-key -L 2>&1";
        kt.delegate = self;
        self.skeyTask = kt;
        [kt launchTask];
        
        self.userTask = nil;
    }
    
    if(self.skeyTask == aPCTask) {
        
        

        
        //[self setToNextStage];
        
        
        
        
        
        self.skeyTask = nil;
    }
    
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    
    NSData *data = [aFileHandler availableData];
    NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    
    Log(@"%@",str);
    
    
    
    if (self.skeyTask == aPCTask){
        
        
        
        
        
        
        
    }
    
    
    
    
    
    NSArray *p = nil;
    for (NSString *key in self.progDict) {
        if ([str containsString:key]){
            p = [self.progDict valueForKey:key];
            break;
        }
    }

    if(p != nil){
        [self.progressLabel setStringValue:[p objectAtIndex:0]];
        [self.progressBar setDoubleValue:[[p objectAtIndex:1] doubleValue]];
        [self.progressBar displayIfNeeded];
    }
}

-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {
    return NO;
}

#pragma mark - IBACTION
-(IBAction)build:(id)sender
{
    // update interface status
    [[RaspberryManager sharedManager] refreshInterface];
    
    // if there is no Ethernet, do not proceed.
    if ([[RaspberryManager sharedManager] ethernetInterface]){
        [self.warningLabel setHidden:YES];
    }else{
        [self resetUIForFailure];
        [self.warningLabel setHidden:NO];
        return;
    }

    // return if there is no node
    NSUInteger nodeCount = MIN([self.nodeList count], MAX_TRIAL_RASP_NODE_COUNT);
    if (nodeCount <= 0){
        // NSAlert
        
        [self resetUIForFailure];
        return;
    }

    [self setUIToProceedState];
    
    // setup actual raspberry nodes
    [[RaspberryManager sharedManager] setupRaspberryNodes:self.nodeList];
    
    // setup hosts address with this
    NSMutableString *nodeip = [NSMutableString new];
    for (NSUInteger i = 0; i < nodeCount; ++i){
        NSDictionary *node = [self.nodeList objectAtIndex:i];
        [nodeip appendString:[NSString stringWithFormat:@"%@ ", [node valueForKey:ADDRESS]]];
    }
    
    NSString *basePath = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
    NSString *sudoSetup =
        [NSString
         stringWithFormat:@"%@/setup/raspberry_sudo_setup.sh %@ %@ %@",
         basePath, basePath,
         [[RaspberryManager sharedManager] ethernetInterface].ip4Address,
         nodeip];

    PCTask *sudoTask = [PCTask new];
    sudoTask.taskCommand = [NSString stringWithFormat:@"bash %@",sudoSetup];
    sudoTask.sudoCommand = YES;
    sudoTask.delegate = self;
    self.sudoTask = sudoTask;
    
    [sudoTask launchTask];
}


#pragma mark - Setup UI status
- (void)setUIToProceedState {
    self.canContinue = NO;
    self.canGoBack = NO;
    [self.buildBtn setEnabled:NO];
    [self.circularProgress startAnimation:nil];
}

-(void)resetUIForFailure {
    [self resetToInitialState];
    [self.circularProgress stopAnimation:nil];
    [self.progressBar setDoubleValue:0.0];
    [self.progressBar displayIfNeeded];
    [self.buildBtn setEnabled:YES];
}

-(void)setToNextStage {
    self.canContinue = YES;
    self.canGoBack = NO;
    
    [self.circularProgress stopAnimation:nil];
    [self.progressBar setDoubleValue:100.0];
    [self.progressBar displayIfNeeded];
    [self.buildBtn setEnabled:NO];
    
    [[Util getApp] setClusterType:PC_CLUSTER_RASPBERRY];
    
    NSViewController *vc3 = [[PCSetup3VC alloc] initWithNibName:@"PCSetup3VC" bundle:[NSBundle mainBundle]];
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kDPNotification_addFinalViewController
     object:self
     userInfo:@{kDPNotification_key_viewController:vc3}];
}

#pragma mark - DPSetupWindowDelegate
-(void)resetToInitialState {
    self.canContinue = NO;
    self.canGoBack = YES;
}

- (void)didRevertToPreviousStage {
    WEAK_SELF(self);
    [[RaspberryManager sharedManager] removeAgentDelegateFromQueue:self];
    [[NSOperationQueue mainQueue]
     addOperationWithBlock:^{
         if(belf){
             [belf removeViewControler];
         }
     }];
}

- (void)removeViewControler {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kDPNotification_deleteViewController
     object:self
     userInfo:@{kDPNotification_key_viewControllerClass:[PCSetup2RPVC class]}];
}

@end
