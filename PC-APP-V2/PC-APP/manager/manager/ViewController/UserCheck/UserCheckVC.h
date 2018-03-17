//
//  UserCheckVC.h
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "BaseSetupVC.h"

@interface UserCheckVC : BaseSetupVC
@property (nonatomic, strong) IBOutlet NSView *pannel;
@property (nonatomic, weak) IBOutlet NSTextField *fieldCode;
@property (nonatomic, weak) IBOutlet NSButton *btnCheck;
@property (nonatomic, weak) IBOutlet NSButton *btnCancel;
@property (nonatomic, weak) NSProgressIndicator *circularProgress;

- (IBAction) check:(id)sender;
- (IBAction) cancel:(id)sender;
- (void) enableControls;
@end
