//
//  IntroVC.h
//  manager
//
//  Created by Almighty Kim on 10/19/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <WebKit/WebKit.h>
#import "BaseSetupVC.h"

@interface IntroVC : BaseSetupVC <WebPolicyDecisionListener
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= __MAC_10_11
, WebPolicyDelegate, WebFrameLoadDelegate, WebUIDelegate
#endif
>

@property (weak) IBOutlet WebView *webView;

@end
