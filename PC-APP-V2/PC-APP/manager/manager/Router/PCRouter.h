//
//  PCRouter.h
//  manager
//
//  Created by Almighty Kim on 8/14/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "PCRouteRequest.h"

@interface PCRouter : NSObject

+ (instancetype)sharedRouter;

// All these method should only be invoked on MAIN THREAD!
- (void) addGetRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath withHandler:(ResponseHandler)aHandler;
- (void) delGetRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath;

- (void) addPostRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath withHandler:(ResponseHandler)aHandler;
- (void) delPostRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath;

- (void) responseFor:(NSString *)aMethod onPath:(NSString *)aPath withPayload:(NSDictionary *)aResponse;

+ (void) routeRequestGet:(const char*)aPath;
+ (void) routeRequestPost:(const char*)aPath withRequestBody:(NSDictionary *)aRequestBody;

@end
