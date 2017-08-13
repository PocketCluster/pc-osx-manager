//
//  PCRouteTrie.h
//  manager
//
//  Created by Almighty Kim on 8/13/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "PCRouteRequest.h"

@interface PCRouteTrie : NSObject
- (instancetype) initWithPathComponent:(NSString *)aComponent;

- (void) addRequest:(NSObject<PCRouteRequest> *)aRequest forMethod:(NSString*)aMethod onPath:(NSString*)aPath;
- (void) delRequest:(NSObject<PCRouteRequest> *)aRequest forMethod:(NSString*)aMethod onPath:(NSString*)aPath;
- (NSObject<PCRouteRequest> *)findRequestForMethod:(NSString*)aMethod onPath:(NSString*)aPath;
- (NSArray *)findAllRequestForMethod:(NSString*)aMethod onPath:(NSString*)aPath;
@end
