//
//  MonitorAppCheck.h
//  manager
//
//  Created by Almighty Kim on 11/13/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

@protocol MonitorAppCheck <NSObject>

@required
// check system readiness
- (void) didAppCheckSystemReadiness:(BOOL)isReady;

// check app has been expried
- (void) didAppCheckAppExpiration:(BOOL)isExpired;

// check if first time run
- (void) didAppCheckIsFirstRun:(BOOL)isFirstRun;

// check if user is authed
- (void) didAppCheckUserAuthed:(BOOL)isUserAuthed;
@end