//
//  TaskOutputVC.h
//  manager
//
//  Created by Almighty Kim on 10/23/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <Cocoa/Cocoa.h>

@interface TaskOutputVC : NSViewController

@property (unsafe_unretained) IBOutlet NSTextView *outputTextView;
@property (weak) IBOutlet NSProgressIndicator *progressBar;
@property (weak) IBOutlet NSTextField *taskCommandLabel;
@property (weak) IBOutlet NSTextField *taskStatusLabel;
@property (weak) IBOutlet NSButton *closeWindowButton;
@property (weak) IBOutlet NSButton *cancelButton;

@property (strong, nonatomic) id target;
@property (strong, nonatomic) NSString *taskCommand;
@property (strong, nonatomic) NSString *taskAction;
@property (strong, nonatomic) NSTask *task;
@property (strong, nonatomic) NSString *windowUUID;

- (IBAction)closeButtonClicked:(id)sender;
- (IBAction)cancelButtonClicked:(id)sender;

@end
