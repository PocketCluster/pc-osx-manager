//
//  PCRouter.h
//  manager
//
//  Created by Almighty Kim on 8/14/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "Node.h"
#import "Package.h"
#import "Cluster.h"
#import "PCRouteRequest.h"

@interface PCRouter : NSObject

+ (PCRouter *)sharedRouter;

- (void) addGetRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath;
- (void) delGetRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath;

- (void) addPostRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath;
- (void) delPostRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath;

- (void) responseFor:(NSString *)aMethod onPath:(NSString *)aPath withPayload:(NSDictionary *)aPayload;

@end
