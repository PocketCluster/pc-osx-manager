//
//  PCProcManager.h
//  manager
//
//  Created by Almighty Kim on 11/2/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

@interface PCProcManager : NSObject
+(instancetype)sharedManager;

- (void)startSalt;
- (void)stopSalt;
- (void)freshStart;

- (void)startWebServer;
- (void)stopWebServer;
@end
