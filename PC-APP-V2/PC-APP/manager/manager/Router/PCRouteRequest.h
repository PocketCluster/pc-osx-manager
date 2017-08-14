//
//  PCRouteRequest.h
//  manager
//
//  Created by Almighty Kim on 8/14/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "routepath.h"

typedef void (^ResponseHandler)(NSString *method, NSString *path, NSDictionary *payload);

@protocol PCRouteRequest <NSObject>
@end