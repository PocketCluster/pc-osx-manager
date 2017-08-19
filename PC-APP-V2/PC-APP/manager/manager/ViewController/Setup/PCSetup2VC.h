//
//  PCSetup2VC.h
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "BaseSetupVC.h"

@interface PCSetup2VC : BaseSetupVC <NSTableViewDataSource, NSTableViewDelegate>
@property (nonatomic, strong) IBOutlet NSView *pannel;
@property (nonatomic, weak) IBOutlet NSProgressIndicator *progress;
@property (nonatomic, weak) IBOutlet NSTableView *nodeTable;
@property (nonatomic, weak) IBOutlet NSButton *btnBuild;
@property (nonatomic, weak) IBOutlet NSButton *btnCancel;

-(IBAction)build:(id)sender;
-(IBAction)cancel:(id)sender;
@end
