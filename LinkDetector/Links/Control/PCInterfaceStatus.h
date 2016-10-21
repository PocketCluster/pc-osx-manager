//
//  PCInterfaceList.h
//  NETUTIL
//
//  Created by Almighty Kim on 10/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCInterfaceTypes.h"

@interface PCInterfaceStatus : NSObject
- (void) startMonitoring;
- (void) stopMonitoring;
@end
