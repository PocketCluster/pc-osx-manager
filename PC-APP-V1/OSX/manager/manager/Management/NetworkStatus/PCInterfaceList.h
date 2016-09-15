//
//  PCInterfaceList.h
//  NETUTIL
//
//  Created by Almighty Kim on 10/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>

#import "LinkInterface.h"

@interface PCInterfaceList : NSObject

+ (NSArray*) all;
+ (BOOL) leaking;
+ (LinkInterface*) interfaceByBSDNumber:(NSInteger)number;

@end
