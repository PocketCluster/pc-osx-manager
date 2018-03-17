//
//  TransitionWC.m
//  manager
//
//  Created by Almighty Kim on 11/5/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "TransitionWC.h"
#import "NullStringChecker.h"

@interface TransitionWC ()
@property (nonatomic, strong) NSString *packageTransition;
@end

@implementation TransitionWC

- (instancetype) initWithPackageExecution:(NSString *)aTransition {
    self = [super initWithWindowNibName:@"TransitionWC"];
    if (self != nil) {
        self.packageTransition = aTransition;
    }
    return self;
}

- (void)windowDidLoad {
    [super windowDidLoad];
    [self.window setBackgroundColor:[NSColor whiteColor]];
    [self.window setTitleVisibility:NSWindowTitleVisible];
    [self.window setTitlebarAppearsTransparent:YES];
    [self.window setShowsResizeIndicator:NO];
    [self.window setShowsToolbarButton:NO];
    self.window.styleMask |= NSFullSizeContentViewWindowMask;

    [self.packageLabel setStringValue:[self packageTransition]];
    [self.circularProgress startAnimation:self];
    [self.circularProgress displayIfNeeded];
}

#pragma mark - MonitorExecution
- (void) onExecutionStartup:(Package *)aPackage {}
- (void) didExecutionStartup:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    if (isSuccess) {
        [self close];
    } else {
        [self.circularProgress setHidden:YES];
        [self.circularProgress stopAnimation:nil];
        [self.circularProgress displayIfNeeded];
        [self.circularProgress removeFromSuperview];
        [self setCircularProgress:nil];

        [self.errorLabel setStringValue:anErrMsg];
        [self.closeBtn setHidden:NO];
    }
}

- (void) onExecutionKill:(Package *)aPackage {}
- (void) didExecutionKill:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    if (isSuccess) {
        [self close];
    } else {
        [self.circularProgress setHidden:YES];
        [self.circularProgress stopAnimation:nil];
        [self.circularProgress displayIfNeeded];
        [self.circularProgress removeFromSuperview];
        [self setCircularProgress:nil];

        [self.errorLabel setStringValue:anErrMsg];
        [self.closeBtn setHidden:NO];
    }
}

- (void) onExecutionProcess:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg {}

#pragma mark - Button Method
- (IBAction) closeWindow:(id)sender {
    [self close];
}
@end
