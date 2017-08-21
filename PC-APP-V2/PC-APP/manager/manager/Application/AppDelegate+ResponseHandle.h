//
//  AppDelegate+ResponseHandle.h
//  manager
//
//  Created by Almighty Kim on 3/24/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate.h"

@interface AppDelegate (ResponseHandle)
- (void)HandleResponseForMethod:(NSString *)aMethod onPath:(NSString *)aPath withPayload:(NSDictionary *)aResponse;
@end
