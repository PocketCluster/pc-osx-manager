//
//  OutputWindow.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import <Cocoa/Cocoa.h>
//#import "VirtualMachineInfo.h"
#import "PCTask.h"


@interface TaskOutputWindow : NSWindowController <NSWindowDelegate, PCTaskDelegate>

@property (strong, nonatomic) id target;
@property (strong, nonatomic) NSString *taskCommand;
@property (strong, nonatomic) NSString *taskAction;
@property (strong, nonatomic) NSTask *task;
@property (nonatomic) BOOL sudoCommand;

@property (unsafe_unretained) IBOutlet NSTextView *outputTextView;
@property (weak) IBOutlet NSProgressIndicator *progressBar;
@property (weak) IBOutlet NSTextField *taskCommandLabel;
@property (weak) IBOutlet NSTextField *taskStatusLabel;
@property (weak) IBOutlet NSButton *closeWindowButton;
@property (weak) IBOutlet NSButton *cancelButton;

- (IBAction)closeButtonClicked:(id)sender;
- (IBAction)cancelButtonClicked:(id)sender;

@end
