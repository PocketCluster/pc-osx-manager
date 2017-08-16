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
@property (nonatomic, weak) IBOutlet NSTextField *fieldEmail;
@property (nonatomic, weak) IBOutlet NSTextField *fieldCode;
@property (nonatomic, weak) IBOutlet NSButton *btnCheck;
@property (nonatomic, weak) IBOutlet NSButton *btnCancel;

-(instancetype) initWithStageControl:(NSObject<StepControl> *)aControl nibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil;

-(IBAction)check:(id)sender;
-(IBAction)cancel:(id)sender;
@end
