//
//  UserCheckVC.h
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <Cocoa/Cocoa.h>
#import "StepStage.h"

@interface UserCheckVC : NSViewController<StageStep>
@property (nonatomic, weak) NSObject<StepControl> *stageControl;
@property (nonatomic, weak) IBOutlet NSProgressIndicator *progress;

-(instancetype) initWithStageControl:(NSObject<StepControl> *)aControl nibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil;

-(IBAction)check:(id)sender;
-(IBAction)cancel:(id)sender;
@end
