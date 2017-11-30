//
//  MessageAlert.h
//  manager
//
//  Created by Almighty Kim on 8/16/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

@interface ShowAlert : NSObject
+ (void) showWarningAlertWithTitle:(NSString *)aTitle message:(NSString *)aMessage;
+ (void) showTerminationAlertWithTitle:(NSString *)aTitle message:(NSString *)aMessage;
+ (void) showAlertMessageWithTitle:(NSString *)aTitle message:(NSString *)aMessage action:(void (^)(NSInteger result))anAction;
@end
