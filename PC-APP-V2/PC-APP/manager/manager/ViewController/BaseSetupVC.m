//
//  BaseSetupVC.m
//  manager
//
//  Created by Almighty Kim on 8/19/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "BaseSetupVC.h"

@interface BaseSetupVC ()
@end

@implementation BaseSetupVC

- (instancetype) initWithStageControl:(NSObject<StepControl> *)aControl nibName:(NSString *)aNibName bundle:(NSBundle *)aBundle {
    self = [super initWithNibName:aNibName bundle:aBundle];
    if (self != nil) {
        self.stageControl = aControl;
        [self finishConstruction];
    }
    return self;
}

- (void)finishConstruction {

}

- (void)dealloc {
    self.stageControl = nil;
}

#pragma mark - StageStep
-(void)didControl:(NSObject<StepControl> *)aControl progressFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult {
}

-(void)didControl:(NSObject<StepControl> *)aControl revertFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult {
}


@end
