//
//  MessageAlert.m
//  manager
//
//  Created by Almighty Kim on 8/16/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "ShowAlert.h"
#import "NullStringChecker.h"

@implementation ShowAlert
+ (void) showWarningAlertWithTitle:(NSString *)aTitle message:(NSString *)aMessage {
    if (ISNULL_STRING(aTitle)) {
        return;
    }
    if (ISNULL_STRING(aMessage)) {
        return;
    }

    NSAlert *alert = [[NSAlert alloc] init];
    [alert addButtonWithTitle:@"OK"];
    //[alert addButtonWithTitle:@"Cancel"];
    [alert setMessageText:aTitle];
    [alert setInformativeText:aMessage];
    [alert setAlertStyle:NSWarningAlertStyle];
    [alert runModal];
}
@end
