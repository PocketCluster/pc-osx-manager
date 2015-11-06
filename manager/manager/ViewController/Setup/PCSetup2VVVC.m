//
//  PCSetup2VVVC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup2VVVC.h"

#import "Util.h"
#import "PCTask.h"
#import "PCSetup3VC.h"
#import "PCProcManager.h"
#import "VagrantManager.h"

@interface PCSetup2VVVC ()<PCTaskDelegate>

@property (strong, nonatomic) PCTask *vagInitTask;
@property (strong, nonatomic) PCTask *sudoTask;
@property (strong, nonatomic) PCTask *saltTask;
@property (strong, nonatomic) PCTask *userTask;

@property (strong, nonatomic) NSDictionary *progDict;

@property (readwrite, nonatomic) BOOL canContinue;
@property (readwrite, nonatomic) BOOL canGoBack;

- (void)setUIToProceedState;
- (void)resetUIForFailure;
- (void)setToNextStage;

- (void)removeViewControler;
@end

@implementation PCSetup2VVVC

@synthesize canContinue;
@synthesize canGoBack;

-(instancetype)initWithNibName:(NSString *)aNibNameOrNil bundle:(NSBundle *)aNibBundleOrNil {
    
    self = [super initWithNibName:aNibNameOrNil bundle:aNibBundleOrNil];
    
    if(self){
        self.progDict = @{@"SUDO_SETUP_STEP_0":@[@"Setting up base configuration.",@10.0]
                           ,@"SUDO_SETUP_DONE":@[@"Finishing configuration.",@20.0]
                           ,@"USER_SETUP_STEP_0":@[@"Starting Vagrant.",@30.0]
                           ,@"USER_SETUP_STEP_1":@[@"Setting up connection.",@70.0]
                           ,@"USER_SETUP_STEP_2":@[@"Finalizing...",@90.0]
                           ,@"USER_SETUP_DONE":@[@"Done!",@100.0]};

        [self resetToInitialState];
    }

    return self;
}



#pragma mark - PCTaskDelegate

-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
    if(aTask.terminationStatus != 0) {
        [self resetUIForFailure];
        [self.progressLabel setStringValue:@"Installation Error. Please try again."];
        
        self.vagInitTask = nil;
        self.sudoTask = nil;
        self.saltTask = nil;
        self.userTask = nil;
        return;
    }
    
    [self setUIToProceedState];
    
    if(self.vagInitTask == aPCTask){
        
        NSString *basePath  = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
        NSString *sudoSetup = [NSString stringWithFormat:@"%@/setup/vagrant_sudo_setup.sh",basePath];
        PCTask *sudoTask = [PCTask new];
        sudoTask.taskCommand = [NSString stringWithFormat:@"sh %@ %@", sudoSetup, basePath];
        sudoTask.sudoCommand = YES;
        sudoTask.delegate = self;
        self.sudoTask = sudoTask;
        [sudoTask launchTask];

        self.vagInitTask = nil;
    }
    
    if (self.sudoTask == aPCTask){

        PCTask *st = [PCTask new];
        st.taskCommand = @"brew install saltstack 2>&1";
        st.delegate = self;
        self.saltTask = st;
        [st launchTask];
        
        self.sudoTask = nil;
    }
    
    if (self.saltTask == aPCTask){
        
        // initiate salt
        [[PCProcManager sharedManager] freshSaltStart];
        sleep(3);
        
        // initiate user task
        NSString *basePath = [[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"];
        NSString *userSetup = [NSString stringWithFormat:@"%@/setup/vagrant_user_setup.sh",basePath];
        PCTask *userTask = [PCTask new];
        userTask.taskCommand = [NSString stringWithFormat:@"sh %@ %@", userSetup, basePath];
        userTask.delegate = self;
        self.userTask = userTask;
        [userTask launchTask];

        self.saltTask = nil;
    }
    
    if(self.userTask == aPCTask){
        
        [self setToNextStage];
        
        self.userTask = nil;
        
    }
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    
    NSData *data = [aFileHandler availableData];
    NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    
    
    Log(@"%@",str);
    
    // save vagrant interface
    if (self.vagInitTask == aPCTask) {
        [[VagrantManager sharedManager] setVboxInterface:str];
        return;
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
    return false;
}

#pragma mark - IBACTION
-(IBAction)build:(id)sender {
    [self setUIToProceedState];

    PCTask *vt = [PCTask new];
    vt.taskCommand = [NSString stringWithFormat:@"python %@/setup/vbox_netiface_create.py",[[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"]];
    vt.delegate = self;
    self.vagInitTask = vt;
    [vt launchTask];
    
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

    [[Util getApp] stopBasicServices];
    [[Util getApp] setClusterType:PC_CLUTER_VAGRANT];

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
     userInfo:@{kDPNotification_key_viewControllerClass:[PCSetup2VVVC class]}];
}

@end
