//
//  AppDelegate+EventHandle.m
//  manager
//
//  Created by Almighty Kim on 3/24/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+EventHandle.h"

/*
 * Event message from engine is most likely a feedback for api call
 * The event message then should be processed and delivered to RunLoop in
 * Default mode to display sucess and failure message to UI in seqeunce
 */

void
PCEventHandle(const char* engineMessage) {

    dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        // Parse in the background
        NSData *msgData = \
            [[NSData alloc]
                 initWithBytesNoCopy:(void *)engineMessage
                 length:strlen((const char*)engineMessage)
                 freeWhenDone:YES];
        
        NSError *error = nil;
        NSDictionary* message = \
            [NSJSONSerialization
                 JSONObjectWithData:msgData
                 options:NSJSONReadingMutableContainers
                 error:&error];

        if (error != nil) {
            Log(@"%@", [error description]);
            return;
        }

        dispatch_async(dispatch_get_main_queue(), ^{
            [(AppDelegate *)[[NSApplication sharedApplication] delegate]
                 performSelectorOnMainThread:@selector(HandleEventMessage:)
                 withObject:message
                 waitUntilDone:NO];
        });
    });

}

@implementation AppDelegate (EventHandle)

- (void)HandleEventMessage:(NSDictionary *)message {
    Log(@"message %@", message);
    
    NSAlert *alert = [[NSAlert alloc] init];
    [alert addButtonWithTitle:@"OK"];
    [alert addButtonWithTitle:@"Cancel"];
    [alert setMessageText:[message valueForKey:@"feed_msg"]];
    [alert setInformativeText:@"Deleted records cannot be restored."];
    [alert setAlertStyle:NSWarningAlertStyle];
    [alert runModal];
}

@end
