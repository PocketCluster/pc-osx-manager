//
//  NativeMenu+Raspberry.h
//  manager
//
//  Created by Almighty Kim on 11/1/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "NativeMenu.h"
#import "RaspberryMenuItem.h"

@interface NativeMenu(Raspberry)<RaspberryMenuItemDelegate>
-(void)raspberryRegisterNotifications;
@end