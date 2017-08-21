//
//  AboutWindowController.h
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#import <WebKit/WebKit.h>
#import "BaseWindowController.h"

@interface AboutWindow : BaseWindowController <WebPolicyDecisionListener
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= __MAC_10_11
, WebPolicyDelegate, WebFrameLoadDelegate, WebUIDelegate
#endif
>

@property (weak) IBOutlet WebView *webView;

@end
