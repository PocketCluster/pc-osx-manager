//
//  PCRouteRequest.h
//  manager
//
//  Created by Almighty Kim on 8/14/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

typedef void (^ResponseHandler)(NSDictionary *payload);

@protocol PCRouteRequest <NSObject>
@required
-(void)responseFor:(NSString *)aMethod onPath:(NSString *)aPath withPayload:(NSDictionary *)aPayload;
@end