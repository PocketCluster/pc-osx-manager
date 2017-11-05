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

    NSString *str = @"<div style=\"display:-webkit-flexbox;display:-webkit-flex;display:flex;-webkit-flex-align:center;-webkit-align-items:center;align-items:center;vertical-align:middle;height:100%;\"><span style=\"font-family:'Helvetica Neue';font-weight:100;font-size:22px;width:100%;text-align:center;\">{PACKAGE_STATE}</span></div>";
    if (!ISNULL_STRING(self.packageTransition)) {
        str = [str stringByReplacingOccurrencesOfString:@"{PACKAGE_STATE}" withString:self.packageTransition];
    }

    [self.webView setDrawsBackground:NO];
    [self.webView.mainFrame loadHTMLString:str baseURL:nil];
    [self.circularProgress startAnimation:self];
    [self.circularProgress displayIfNeeded];
}

#pragma mark - WebView Delegate
- (void)webView:(WebView*)webView decidePolicyForNavigationAction:(NSDictionary*)actionInformation request:(NSURLRequest*)request frame:(WebFrame*)frame decisionListener:(id<WebPolicyDecisionListener>)listener {
    NSString *host = [[request URL] host];
    if(host) {
        [[NSWorkspace sharedWorkspace] openURL:[request URL]];
    } else {
        [listener use];
    }
}

- (void)use {
}

- (void)download {
}

- (void)ignore {
}

#pragma mark - MonitorExecution
- (void) onExecutionStartup:(Package *)aPackage {}
- (void) didExecutionStartup:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    [self close];
}

- (void) onExecutionKill:(Package *)aPackage {}
- (void) didExecutionKill:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg {
    [self close];
}

- (void) onExecutionProcess:(Package *)aPackage success:(BOOL)isSuccess error:(NSString *)anErrMsg {}
@end
