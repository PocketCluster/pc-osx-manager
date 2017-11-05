//
//  IntroVC.m
//  manager
//
//  Created by Almighty Kim on 10/19/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "IntroVC.h"

@implementation IntroVC
@synthesize webView;
@synthesize progressLabel;
@synthesize circularProgress;

- (void)finishConstruction {
    [super finishConstruction];
    [self setTitle:@"PocketCluster - v0.1.4 Early Evaulation Version"];
}

- (void)viewDidLoad {
    [super viewDidLoad];

    NSString *str = @"<div style=\"text-align:center;font-family:'Helvetica Neue';font-weight:100;font-size:38px;width:100%;\">PocketCluster</div><div style=\"text-align:center;font-family:Arial;font-size:13px;width:100%;\">Version {VERSION} - Early Evaluation</div>";
    str = [str stringByReplacingOccurrencesOfString:@"{VERSION}" withString:[[[NSBundle mainBundle] infoDictionary] valueForKey:@"CFBundleShortVersionString"]];
    str = [str stringByReplacingOccurrencesOfString:@"\n" withString:@"<br>"];

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

@end
