//
//  UserCheckVC.m
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "UserCheckVC.h"
#import "BaseBrandView.h"

@interface UserCheckVC ()
@end

@implementation UserCheckVC

- (instancetype) initWithStageControl:(NSObject<StepControl> *)aControl nibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil {
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if (self != nil) {
        [self setTitle:@"Inviatation Check"];
        self.stageControl = aControl;
    }
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];

    [[((BaseBrandView *)self.view) contentBox] removeFromSuperview];
    ((BaseBrandView *)self.view).contentBox = nil;
}

-(IBAction)check:(id)sender {
    [self.stageControl shouldControlProgressFrom:self withParam:nil];
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
@end
