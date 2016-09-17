//
//  PCSetup2RPVC.h
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "DPSetupWindow.h"

@interface PCSetup2RPVC : NSViewController  <DPSetupWindowStageViewController, NSTableViewDataSource, NSTableViewDelegate>

@property (nonatomic, weak) IBOutlet NSProgressIndicator *circularProgress;
@property (nonatomic, weak) IBOutlet NSTableView *nodeTable;
@property (nonatomic, weak) IBOutlet NSTextField *warningLabel;
@property (nonatomic, weak) IBOutlet NSProgressIndicator *progressBar;
@property (nonatomic, weak) IBOutlet NSTextField *progressLabel;
@property (nonatomic, weak) IBOutlet NSButton *buildBtn;

-(IBAction)build:(id)sender;
@end