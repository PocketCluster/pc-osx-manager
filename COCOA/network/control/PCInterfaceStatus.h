//
//  PCInterfaceList.h
//  NETUTIL
//
//  Created by Almighty Kim on 10/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCInterfaceTypes.h"

@class PCInterfaceStatus;
// all the parameters are here to read and C99 variant. Don't release or modify it.
@protocol PCInterfaceStatusNotification <NSObject>
-(void)PCInterfaceStatusChanged:(PCInterfaceStatus *)monitor interfaceStatus:(PCNetworkInterface**)status count:(unsigned int)count;
-(void)PCGatewayStatusChanged:(PCInterfaceStatus *)monitor gatewayStatus:(SCNIGateway**)status count:(unsigned int)count;
@end

@interface PCInterfaceStatus : NSObject
- (instancetype)initWithStatusAudience:(NSObject<PCInterfaceStatusNotification>*)audience;
- (void) startMonitoring;
- (void) stopMonitoring;
@end
