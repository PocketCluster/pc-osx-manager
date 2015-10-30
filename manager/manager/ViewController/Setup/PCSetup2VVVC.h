//
//  PCSetup2VVVC.h
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "DPSetupWindow.h"

@interface PCSetup2VVVC : NSViewController  <DPSetupWindowStageViewController>

@property (nonatomic, weak) IBOutlet NSProgressIndicator *progressBar;
@property (nonatomic, weak) IBOutlet NSTextField *progressLabel;
@property (nonatomic, weak) IBOutlet NSButton *buildBtn;

-(IBAction)build:(id)sender;
@end
