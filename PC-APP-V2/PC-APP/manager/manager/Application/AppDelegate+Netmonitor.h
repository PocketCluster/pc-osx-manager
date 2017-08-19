//
//  AppDelegate+Netmonitor.h
//  manager
//
//  Created by Almighty Kim on 4/10/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "PCInterfaceStatus.h"
#import "AppDelegate.h"

extern bool
PCUpdateInterfaceList(PCNetworkInterface** interfaces, unsigned int count);

extern bool
PCUpdateGatewayList(SCNIGateway** gateways, unsigned int count);

@interface AppDelegate(Netmonitor)<PCInterfaceStatusNotification>
@end
