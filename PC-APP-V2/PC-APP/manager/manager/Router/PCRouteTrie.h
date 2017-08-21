//
//  PCRouteTrie.h
//  manager
//
//  Created by Almighty Kim on 8/13/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "PCRouteRequest.h"

@interface PCRequestHolder : NSObject
@property (nonatomic, strong, readonly) NSObject<PCRouteRequest>* request;
@property (nonatomic, strong, readonly) ResponseHandler handler;
@end

@interface PCRouteTrie : NSObject
- (instancetype) initWithPathComponent:(NSString *)aComponent;

- (void) addRequest:(NSObject<PCRouteRequest> *)aRequest forMethod:(NSString *)aMethod onPath:(NSString *)aPath withHandler:(ResponseHandler)aHandler;
- (void) delRequest:(NSObject<PCRouteRequest> *)aRequest forMethod:(NSString *)aMethod onPath:(NSString *)aPath;
- (PCRequestHolder *)findRequestForMethod:(NSString *)aMethod onPath:(NSString *)aPath;
- (NSArray<PCRequestHolder *> *)findAllRequestForMethod:(NSString*)aMethod onPath:(NSString*)aPath;
@end
