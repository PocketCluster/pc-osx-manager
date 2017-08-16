//
//  AgreementVC.h
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <Cocoa/Cocoa.h>
#import "StepStage.h"

@interface AgreementVC : NSViewController<StageStep>
@property (nonatomic, weak) NSObject<StepControl> *stageControl;
@property (nonatomic, assign) IBOutlet NSTextView *agreement;

- (instancetype) initWithStageControl:(NSObject<StepControl> *)aControl nibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil;
-(IBAction)agreed:(id)sender;
-(IBAction)declined:(id)sender;
@end
