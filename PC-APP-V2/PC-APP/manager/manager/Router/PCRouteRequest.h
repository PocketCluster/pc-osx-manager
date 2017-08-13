//
//  PCRouteRequest.h
//  manager
//
//  Created by Almighty Kim on 8/14/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

@protocol PCRouteRequest <NSObject>
@required
-(void)requestPath:(NSString *)aPath forPayload:(NSDictionary *)aDict;
@end