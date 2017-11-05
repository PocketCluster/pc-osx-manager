//
//  TransitionWC.h
//  manager
//
//  Created by Almighty Kim on 11/5/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <WebKit/WebKit.h>
#import "BaseWindowController.h"
#import "MonitorExecution.h"

@interface TransitionWC : BaseWindowController <MonitorExecution, WebPolicyDecisionListener
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= __MAC_10_11
, WebPolicyDelegate, WebFrameLoadDelegate, WebUIDelegate
#endif
>
@property (weak) IBOutlet WebView *webView;
@property (weak) IBOutlet NSProgressIndicator *circularProgress;
- (instancetype) initWithPackageExecution:(NSString *)aTransition;
@end
