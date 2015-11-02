//
//  PCProcManager.h
//  manager
//
//  Created by Almighty Kim on 11/2/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface PCProcManager : NSObject
+(instancetype)sharedManager;

- (void)startSalt;
- (void)stopSalt;
- (void)freshStart;
@end
