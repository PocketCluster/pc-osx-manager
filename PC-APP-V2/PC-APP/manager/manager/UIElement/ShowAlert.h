//
//  MessageAlert.h
//  manager
//
//  Created by Almighty Kim on 8/16/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>

extern NSString * const ALRT_MESSAGE_TEXT;
extern NSString * const ALRT_INFORMATIVE_TEXT;;

@interface ShowAlert : NSObject
+ (void) showWarningAlertFromMeta:(NSDictionary *)aMeta;
@end
