//
//  PCSetup3VC.h
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "BaseSetupVC.h"

@interface PCSetup3VC : BaseSetupVC <NSTableViewDataSource, NSTableViewDelegate>
@property (nonatomic, strong) IBOutlet NSView *pannel;

@property (nonatomic, weak) IBOutlet NSTableView *packageTable;
@property (nonatomic, weak) IBOutlet NSProgressIndicator *progressBar;
@property (nonatomic, weak) IBOutlet NSTextField *progressLabel;
@property (nonatomic, weak) NSProgressIndicator *circularProgress;

@property (nonatomic, weak) IBOutlet NSButton *btnInstall;
@property (nonatomic, weak) IBOutlet NSButton *btnCancel;

-(IBAction)install:(id)sender;
-(IBAction)cancel:(id)sender;
@end
