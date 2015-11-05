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

#import "PCSetup3VC.h"


@interface PCSetup1VC ()
@property (readwrite, nonatomic) BOOL hideContinue;
@property (readwrite, nonatomic) BOOL hideGoBack;

@end

@implementation PCSetup1VC

@synthesize hideContinue;
@synthesize hideGoBack;

-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if(self){
        [self setHideContinue:YES];
        [self setHideGoBack:YES];
    }
    return self;
}

#pragma mark - DPSetupWindowDelegate
- (void)resetToInitialState {
    [self setHideContinue:YES];
    [self setHideGoBack:YES];
}

- (IBAction)setupVagrantCluster:(id)sender {
    
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




@end
