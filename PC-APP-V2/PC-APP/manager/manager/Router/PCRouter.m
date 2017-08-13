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

#import "PCRoutePathConst.h"
#import "PCRouteTrie.h"
#import "PCRouter.h"

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

- (void) addGetRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath {
    [self.rootNode addRequest:aRequest forMethod:RPATH_EVENT_METHOD_GET onPath:aPath];
}

- (void) delGetRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath {
    [self.rootNode delRequest:aRequest forMethod:RPATH_EVENT_METHOD_GET onPath:aPath];
}


- (void) addPostRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath {
    [self.rootNode addRequest:aRequest forMethod:RPATH_EVENT_METHOD_POST onPath:aPath];
}

- (void) delPostRequest:(NSObject<PCRouteRequest> *)aRequest onPath:(NSString*)aPath {
    [self.rootNode delRequest:aRequest forMethod:RPATH_EVENT_METHOD_POST onPath:aPath];
}

- (void) responseFor:(NSString *)aMethod onPath:(NSString *)aPath withPayload:(NSDictionary *)aPayload {
    NSObject<PCRouteRequest>* request = [self.rootNode findRequestForMethod:aMethod onPath:aMethod];
    [request responseFor:aMethod onPath:aPath withPayload:aPayload];
}

@end
