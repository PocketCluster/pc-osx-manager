//
//  PCSetup1VC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCSetup1VC.h"
#import "PCSetup2VVVC.h"
#import "PCSetup2RPVC.h"
#import "PCTask.h"
#import "PCConstants.h"

#import "PCSetup3VC.h"

@interface PCSetup1VC ()<PCTaskDelegate>
@property (readwrite, nonatomic) BOOL hideContinue;
@property (readwrite, nonatomic) BOOL hideGoBack;

@property (nonatomic, strong) PCTask *taskLibChecker;
@property (nonatomic, readwrite) int libraryCheckupResult;

- (void)warnLibraryDeficiency;
@end

@implementation PCSetup1VC

@synthesize hideContinue;
@synthesize hideGoBack;

-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if(self){
        [self setHideContinue:YES];
        [self setHideGoBack:YES];
        
        
        // check basic libary status
        PCTask *lc = [[PCTask alloc] init];
        lc.taskCommand = [NSString stringWithFormat:@"sh %@/setup/check_base_library.sh",[[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"]];
        lc.delegate = self;
        self.taskLibChecker = lc;
        [lc launchTask];

    }
    return self;
}

#pragma mark - DPSetupWindowDelegate
- (void)resetToInitialState {
    [self setHideContinue:YES];
    [self setHideGoBack:YES];
}

- (IBAction)setupVagrantCluster:(id)sender {
    
    if(self.libraryCheckupResult != 0){
        [self warnLibraryDeficiency];
        return;
    }

    
    NSViewController *vc3 = [[PCSetup3VC alloc] initWithNibName:@"PCSetup3VC" bundle:[NSBundle mainBundle]];
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kDPNotification_addNextViewControllerAndProceed
     object:self
     userInfo:@{kDPNotification_key_viewController:vc3}];

    
    return;
    
    
    
    
    NSViewController *vc2v =
        [[PCSetup2VVVC alloc]
         initWithNibName:@"PCSetup2VVVC"
         bundle:[NSBundle mainBundle]];
    
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kDPNotification_addNextViewControllerAndProceed
     object:self
     userInfo:@{kDPNotification_key_viewController:vc2v}];
}

- (IBAction)setupRaspberryCluster:(id)sender {
    NSViewController *vc2r =
        [[PCSetup2RPVC alloc]
         initWithNibName:@"PCSetup2RPVC"
         bundle:[NSBundle mainBundle]];
    
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kDPNotification_addNextViewControllerAndProceed
     object:self
     userInfo:@{kDPNotification_key_viewController:vc2r}];
}

#pragma mark - Warning
-(void)warnLibraryDeficiency {
    switch (self.libraryCheckupResult) {
        case PC_LIB_VAGRANT:{
            
            [[NSAlert
              alertWithMessageText:@"Vagrant is not found in the system. Please install Vagrant and restart."
              defaultButton:@"OK"
              alternateButton:nil
              otherButton:nil
              informativeTextWithFormat:@""] runModal];
            
            break;
        }
        case PC_LIB_VIRTUABOX:{
            
            [[NSAlert
              alertWithMessageText:@"Virtualbox is not found in the system. Please install Virtualbox and restart."
              defaultButton:@"OK"
              alternateButton:nil
              otherButton:nil
              informativeTextWithFormat:@""] runModal];
            break;
        }
        default:{
            break;
        }
    }
}

#pragma mark - PCTaskDelegate
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    int term = [aTask terminationStatus];
    [self setLibraryCheckupResult:term];
    [self setTaskLibChecker:nil];
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
}

-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {
    return NO;
}


@end
