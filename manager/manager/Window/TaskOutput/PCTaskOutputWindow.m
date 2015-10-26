//
//  PCTaskOutputWindow.m
//  manager
//
//  Created by Almighty Kim on 10/26/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCTaskOutputWindow.h"
#import "AppDelegate.h"
#import "Util.h"

@interface PCTaskOutputWindow()
@property (strong, nonatomic) PCTask *taskOperator;
@end


@implementation PCTaskOutputWindow
{
    BOOL _isClosed;
}

@dynamic target;
@dynamic taskCommand;
@dynamic taskAction;
@dynamic task;
@dynamic sudoCommand;

- (void)finishConstruction {
    [super finishConstruction];
    _isClosed = false;
    self.taskOperator = [PCTask new];
    self.taskOperator.delegate = self;
}

- (void)launchTask
{
    /*
     NSString *name = [self.target isKindOfClass:[VagrantMachine class]] ?
     [NSString stringWithFormat:@"%@ - %@",((VagrantMachine*)self.target).instance.displayName, ((VagrantMachine*)self.target).name] :
     ((VagrantInstance*)self.target).displayName;
     
     self.window.title = [NSString stringWithFormat:@"%@ %@", name, self.taskAction];
     */
    
    self.title = @"Task Window";
    
    self.taskCommandLabel.stringValue = self.taskOperator.taskCommand;
    self.taskStatusLabel.stringValue = @"Running task...";
    [self.progressBar startAnimation:self];
    
    [self.taskOperator launchTask];
}

- (void)windowWillClose:(NSNotification *)notification {
    
    Log(@"%s", __PRETTY_FUNCTION__);
    
    @synchronized(self) {
        _isClosed = YES;
    }

    [super windowWillClose:notification];
}

- (IBAction)closeButtonClicked:(id)sender {
    [self close];
}

- (IBAction)cancelButtonClicked:(id)sender {
    NSAlert *confirmAlert = [NSAlert alertWithMessageText:@"Are you sure you want to cancel the running task?" defaultButton:@"Confirm" alternateButton:@"Cancel" otherButton:nil informativeTextWithFormat:@""];
    NSInteger button = [confirmAlert runModal];
    
    if(button == NSAlertDefaultReturn) {
        [self.taskOperator cancelTask];
    }
}

#pragma mark - PCTaskDelegate
-(void)setTarget:(id)aTarget {
    [self.taskOperator setTarget:aTarget];
}

-(void)setTaskCommand:(NSString *)aTaskCommand {
    [self.taskOperator setTaskCommand:aTaskCommand];
}

-(void)setTaskAction:(NSString *)aTaskAction {
    [self.taskOperator setTaskAction:aTaskAction];
}

-(void)setTask:(NSTask *)aTask {
    [self.taskOperator setTask:aTask];
}

-(void)setSudoCommand:(BOOL)aSudoCommand {
    [self.taskOperator setSudoCommand:aSudoCommand];
}

-(void)task:(PCTask *)aPCTask taskCompletion:(NSTask *)aTask {
    
    [self.progressBar stopAnimation:self];
    [self.progressBar setIndeterminate:NO];
    [self.progressBar setDoubleValue:self.progressBar.maxValue];
    
    NSButton *closeButton = [self standardWindowButton:NSWindowCloseButton];
    [closeButton setEnabled:YES];
    
    [self.closeWindowButton setEnabled:YES];
    [self.cancelButton setHidden:YES];
    
    NSString *notificationText;
    
    if(aTask.terminationStatus != 0) {
        self.taskStatusLabel.stringValue = @"Completed with errors";
        notificationText = @"Task completed with errors";
        
    } else {
        self.taskStatusLabel.stringValue = @"Completed successfully";
        notificationText = @"Task completed successfully";
    }
    
    /*
     NSString *name = [self.target isKindOfClass:[VagrantMachine class]] ? [NSString stringWithFormat:@"%@ - %@",((VagrantMachine*)self.target).instance.displayName, ((VagrantMachine*)self.target).name] : ((VagrantInstance*)self.target).displayName;
     
     [[Util getApp] showUserNotificationWithTitle:notificationText informativeText:[NSString stringWithFormat:@"%@ %@", name, self.taskAction] taskWindowUUID:self.taskOperator.taskUUID];
     */
        
    if([[NSUserDefaults standardUserDefaults] boolForKey:@"autoCloseTaskWindows"] && aTask.terminationStatus == 0) {
        dispatch_async(dispatch_get_global_queue(0,0), ^{
            [self close];
        });
    }
}

-(void)task:(PCTask *)aPCTask recievedOutput:(NSFileHandle *)aFileHandler {
    
    NSData *data = [aFileHandler availableData];
    NSString *str = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    
    @synchronized(self) {
        if (!_isClosed) {
            //smart scrolling logic for command output
            BOOL scroll = (NSMaxY(self.outputTextView.visibleRect) == NSMaxY(self.outputTextView.bounds));
            
            [self.outputTextView.textStorage appendAttributedString:[[NSAttributedString alloc] initWithString:str]];
            
            if([NSFont fontWithName:@"Menlo" size:11])
            {
                [self.outputTextView.textStorage setFont:[NSFont fontWithName:@"Menlo" size:11]];
            }
            
            if (scroll)
            {
                [self.outputTextView scrollRangeToVisible: NSMakeRange(self.outputTextView.string.length, 0)];
            }
        }
    }
}

-(BOOL)task:(PCTask *)aPCTask isOutputClosed:(id<PCTaskDelegate>)aDelegate {
    return _isClosed;
}


@end
