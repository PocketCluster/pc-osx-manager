//
//  MessageAlert.h
//  manager
//
//  Created by Almighty Kim on 8/16/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface ShowAlert : NSObject
+ (void) showWarningAlertWithTitle:(NSString *)aTitle message:(NSString *)aMessage;
@end
