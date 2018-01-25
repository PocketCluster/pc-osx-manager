//
//  TransitionWC.h
//  manager
//
//  Created by Almighty Kim on 11/5/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "BaseWindowController.h"
#import "MonitorExecution.h"

@interface TransitionWC : BaseWindowController <MonitorExecution>
@property (weak) IBOutlet NSTextField *packageLabel;
@property (weak) IBOutlet NSTextField *errorLabel;
@property (weak) IBOutlet NSButton *closeBtn;
@property (weak) IBOutlet NSProgressIndicator *circularProgress;
- (instancetype) initWithPackageExecution:(NSString *)aTransition;
@end
