//
//  DeviceSerialNumber.h
//  manager
//
//  Created by Almighty Kim on 10/21/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>

extern CFTypeRef CopySerialNumber(void);


@interface DeviceSerialNumber : NSObject
+ (NSString *)deviceSerialNumber;
@end

