//
//  AppDelegate.h
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

@interface AppDelegate : NSObject <NSApplicationDelegate, NSMenuDelegate, NSUserNotificationCenterDelegate>
- (void)addOpenWindow:(id)window;
- (void)removeOpenWindow:(id)window;
- (void)updateProcessType;
- (NSImage*)getThemedImage:(NSString*)imageName;
- (NSString*)getCurrentTheme;


- (void)startPCMultSrv;
- (void)stopPCMultSrv;
@end

