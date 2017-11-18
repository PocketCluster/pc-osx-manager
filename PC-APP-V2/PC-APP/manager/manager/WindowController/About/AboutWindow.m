//
//  AboutWindowController.m
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#import "AboutWindow.h"

@interface AboutWindow ()
@end

@implementation AboutWindow
@synthesize webView;

- (id)initWithWindow:(NSWindow *)window {
    self = [super initWithWindow:window];

    return self;
}

- (void)windowDidLoad {
    [super windowDidLoad];

    [self.window setBackgroundColor:[NSColor whiteColor]];
    [self.window setTitleVisibility:NSWindowTitleHidden];
    [self.window setTitlebarAppearsTransparent:YES];
    self.window.styleMask |= NSFullSizeContentViewWindowMask;

    NSString *str = @"<div style=\"text-align:left;font-family:'Helvetica Neue';font-weight:100;font-size:38px;\">PocketCluster</div><div style=\"text-align:left;font-family:Arial;font-size:13px\">Version {VERSION} - Early Evaluation<br/><br/>Copyright &copy;2015,{YEAR} PocketCluster<br/><a href=\"{URL}\">{URL}</a></div>";

    NSString *dateString = [NSString stringWithCString:__DATE__ encoding:NSASCIIStringEncoding];
    NSString *yearString = [dateString substringWithRange:NSMakeRange([dateString length] - 4, 4)];

    str = [str stringByReplacingOccurrencesOfString:@"{YEAR}" withString:yearString];
    str = [str stringByReplacingOccurrencesOfString:@"{VERSION}" withString:[[[NSBundle mainBundle] infoDictionary] valueForKey:@"CFBundleShortVersionString"]];
    str = [str stringByReplacingOccurrencesOfString:@"{URL}" withString:@"https://github.com/pocketcluster/pocketcluster"];
    str = [str stringByReplacingOccurrencesOfString:@"\n" withString:@"<br>"];

    //self.webView.policyDelegate = self;
    [self.webView setDrawsBackground:NO];
    [self.webView.mainFrame loadHTMLString:str baseURL:nil];
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

@end
