//
//  PCSetup3VC.h
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "DPSetupWindow.h"

@interface PCSetup3VC : NSViewController <DPSetupWindowStageViewController, NSTableViewDataSource, NSTableViewDelegate>
@property (nonatomic, weak) IBOutlet NSTableView *packageTable;
-(IBAction)install:(id)sender;
@end
