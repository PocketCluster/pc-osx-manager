//
//  AppDelegate+EventHandle.h
//  manager
//
//  Created by Almighty Kim on 3/24/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "AppDelegate.h"

@interface AppDelegate (EventHandle)
- (void)HandleEventForMethod:(NSString *)aMethod onPath:(NSString *)aPath withPayload:(NSDictionary *)aPayload;
@end
