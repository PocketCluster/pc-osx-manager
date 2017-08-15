//
//  MessageAlert.m
//  manager
//
//  Created by Almighty Kim on 8/16/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "ShowAlert.h"

NSString * const ALRT_MESSAGE_TEXT     = @"alrt_message_text";
NSString * const ALRT_INFORMATIVE_TEXT = @"alrt_informative_text";

@implementation ShowAlert
+ (void) showWarningAlertFromMeta:(NSDictionary *)aMeta {
    if (aMeta == nil) {
        return;
    }
    if ([aMeta objectForKey:ALRT_MESSAGE_TEXT] == nil) {
        return;
    }
    if ([aMeta objectForKey:ALRT_INFORMATIVE_TEXT] == nil) {
        return;
    }

    NSAlert *alert = [[NSAlert alloc] init];
    [alert addButtonWithTitle:@"OK"];
    [alert addButtonWithTitle:@"Cancel"];
    [alert setMessageText:[aMeta objectForKey:ALRT_MESSAGE_TEXT]];
    [alert setInformativeText:[aMeta objectForKey:ALRT_INFORMATIVE_TEXT]];
    [alert setAlertStyle:NSWarningAlertStyle];
    [alert runModal];
}
@end
