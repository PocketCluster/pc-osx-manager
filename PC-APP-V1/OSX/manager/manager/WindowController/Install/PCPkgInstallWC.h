//
//  PCPkgInstallWC.h
//  manager
//
//  Created by Almighty Kim on 11/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "BaseWindowController.h"

@interface PCPkgInstallWC : BaseWindowController <NSTableViewDataSource, NSTableViewDelegate>

@property (nonatomic, weak) IBOutlet NSProgressIndicator *circularProgress;
@property (nonatomic, weak) IBOutlet NSProgressIndicator *progressBar;
@property (nonatomic, weak) IBOutlet NSTableView *packageTable;
@property (nonatomic, weak) IBOutlet NSTextField *progressLabel;
@property (nonatomic, weak) IBOutlet NSButton *installBtn;
@property (nonatomic, weak) IBOutlet NSButton *closeBtn;

- (IBAction)install:(id)sender;

@end
