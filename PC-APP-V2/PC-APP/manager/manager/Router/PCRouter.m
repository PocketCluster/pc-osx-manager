//
//  PCRouter.m
//  manager
//
//  Created by Almighty Kim on 8/14/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "SynthesizeSingleton.h"
#import "PCConstants.h"
#import "NullStringChecker.h"

#import "Node.h"
#import "Package.h"
#import "Cluster.h"

#import "PCRoutePathConst.h"
#import "PCRouteTrie.h"
#import "PCRouter.h"
#import "pc-core.h"

@interface PCRouter() {
    __strong PCRouteTrie *_rootNode;
}
@property (nonatomic, strong, readonly) PCRouteTrie *rootNode;
@end

@implementation PCRouter
@synthesize rootNode = _rootNode;
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(PCRouter, sharedRouter);

- (id)init {
    self = [super init];
    if(self) {
        _rootNode = [[PCRouteTrie alloc] initWithPathComponent:@"/"];
    }
    return self;
}

- (void)dealloc {
    _rootNode = nil;
}

- (void) addGetRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath withHandler:(ResponseHandler)aHandler {
    NSAssert([NSThread isMainThread], @"Request should only be added in Main Thread!");

    [self.rootNode addRequest:aRequest forMethod:RPATH_EVENT_METHOD_GET onPath:aPath withHandler:aHandler];
}

- (void) delGetRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath {
    NSAssert([NSThread isMainThread], @"Request should only be removed in Main Thread!");

    [self.rootNode delRequest:aRequest forMethod:RPATH_EVENT_METHOD_GET onPath:aPath];
}

- (void) addPostRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath withHandler:(ResponseHandler)aHandler {
    NSAssert([NSThread isMainThread], @"Request should only be added in Main Thread!");

    [self.rootNode addRequest:aRequest forMethod:RPATH_EVENT_METHOD_POST onPath:aPath withHandler:aHandler];
}

- (void) delPostRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath {
    NSAssert([NSThread isMainThread], @"Request should only be removed in Main Thread!");

    [self.rootNode delRequest:aRequest forMethod:RPATH_EVENT_METHOD_POST onPath:aPath];
}

- (void) responseFor:(NSString *)aMethod onPath:(NSString *)aPath withPayload:(NSDictionary *)aResponse {
    NSAssert([NSThread isMainThread], @"Reponse should only be received and handled in Main Thread!");

    PCRequestHolder *holder = [self.rootNode findRequestForMethod:aMethod onPath:aPath];
    if (holder != nil) {
        holder.handler(aMethod, aPath, aResponse);
    }
}

+ (void) routeRequestGet:(const char*)aPath {
    NSAssert([NSThread isMainThread], @"Request should only be made in Main Thread!");

    if (aPath == NULL || strlen(aPath) == 0) {
        return;
    }

    RouteRequestGet((char *)aPath);
}

#if 0
+ (NSDictionary*) _percentEncodeDictionayValues:(NSDictionary*) dict {
    NSMutableDictionary* edict=[NSMutableDictionary dictionaryWithDictionary:dict];
    
    NSMutableCharacterSet* URLQueryPartAllowedCharacterSet = [[NSCharacterSet URLQueryAllowedCharacterSet] mutableCopy];
    [URLQueryPartAllowedCharacterSet removeCharactersInString:@"?&=@+/'"];
    
    for(NSString* key in [dict allKeys]) {
        if([dict[key] isKindOfClass:[NSString class]]) {
            edict[key] = [dict[key] stringByAddingPercentEncodingWithAllowedCharacters:URLQueryPartAllowedCharacterSet];
        }
    }
    return edict;
}
#endif

+ (void) routeRequestPost:(const char*)aPath withRequestBody:(NSDictionary *)aRequestBody {
    NSAssert([NSThread isMainThread], @"Request should only be made in Main Thread!");

    if (aPath == NULL || strlen(aPath) == 0) {
        return;
    }
    if (aRequestBody == nil || [aRequestBody count] == 0) {
        return;
    }

    NSError *error = nil;
    NSData *data = [NSJSONSerialization dataWithJSONObject:aRequestBody options:0 error:&error];
    if (error != nil) {
        Log(@"%@", [error description]);
        return;
    }

    Log(@"routeRequestPost aPath[%s], aRequestBody[%s]\n", aPath, (char *)[data bytes]);

    RouteRequestPost((char *)aPath, (char *)[data bytes]);
}

@end
