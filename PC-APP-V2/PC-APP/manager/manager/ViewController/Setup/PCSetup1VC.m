//
//  PCSetup1VC.m
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCConstants.h"

#import "PCSetup1VC.h"
#import "PCSetup2VC.h"
#import "PCSetup3VC.h"

@interface PCSetup1VC ()
//<PCTaskDelegate>
@property (readwrite, nonatomic) BOOL hideContinue;
@property (readwrite, nonatomic) BOOL hideGoBack;

@property (nonatomic, strong) id taskLibChecker;
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
    
#if 0
        // check basic libary status
        PCTask *lc = [[PCTask alloc] init];
        lc.taskCommand = [NSString stringWithFormat:@"bash %@/setup/check_vagrant_library.sh",[[[NSBundle mainBundle] resourcePath] stringByAppendingPathComponent:@"Resources.bundle/"]];
        lc.delegate = self;
        self.taskLibChecker = lc;
        [lc launchTask];
#endif
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
    
#if 0
    NSViewController *vc3c =
        [[PCSetup3VC alloc]
         initWithNibName:@"PCSetup3VC"
         bundle:[NSBundle mainBundle]];
    
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kDPNotification_addNextViewControllerAndProceed
     object:self
     userInfo:@{kDPNotification_key_viewController:vc3c}];
    
    return;
    
    NSViewController *vc2v =
        [[PCSetup2VVVC alloc]
         initWithNibName:@"PCSetup2VVVC"
         bundle:[NSBundle mainBundle]];
    
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kDPNotification_addNextViewControllerAndProceed
     object:self
     userInfo:@{kDPNotification_key_viewController:vc2v}];
#endif
}

- (IBAction)setupRaspberryCluster:(id)sender {

#if 0
    NSViewController *vc3c =
        [[PCSetup3VC alloc]
         initWithNibName:@"PCSetup3VC"
         bundle:[NSBundle mainBundle]];
    
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kDPNotification_addNextViewControllerAndProceed
     object:self
     userInfo:@{kDPNotification_key_viewController:vc3c}];
    
    return;
#endif
    
    NSViewController *vc2r =
        [[PCSetup2VC alloc]
         initWithNibName:@"PCSetup2VC"
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

#if 0
#pragma mark - PCTaskDelegate
-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    int term = [aTask terminationStatus];
    [self setLibraryCheckupResult:term];
    [self setTaskLibChecker:nil];
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {}
-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {return NO;}
#endif

@end
