//
//  BaseSetupVC.h
//  manager
//
//  Created by Almighty Kim on 8/19/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "StepStage.h"
#import "NSGeometry.h"
#import "BaseBrandView.h"

@interface BaseSetupVC : NSViewController<StageStep>
@property (nonatomic, weak) NSObject<StepControl> *stageControl;

- (instancetype) initWithStageControl:(NSObject<StepControl> *)aControl nibName:(NSString *)aNibName bundle:(NSBundle *)aBundle;
- (void)finishConstruction;
- (void)prepareDestruction;
@end
