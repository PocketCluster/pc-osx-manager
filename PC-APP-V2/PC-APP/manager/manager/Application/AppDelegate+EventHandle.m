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

static NSString * const RPATH_EVENT_METHOD_GET    = @"GET";
static NSString * const RPATH_EVENT_METHOD_POST   = @"POST";
static NSString * const RPATH_EVENT_METHOD_PUT    = @"PUT";
static NSString * const RPATH_EVENT_METHOD_DELETE = @"DELETE";

static void
eventHandle(NSString* method, const char* path, const char* payload);

void
PCEventFeedGet(char* path) {
    eventHandle(RPATH_EVENT_METHOD_GET, (const char*)path, NULL);
}

void
PCEventFeedPost(char* path, char* payload) {
    eventHandle(RPATH_EVENT_METHOD_POST, (const char*)path, (const char*)payload);
}

void
PCEventFeedPut(char* path, char* payload) {
    eventHandle(RPATH_EVENT_METHOD_PUT, (const char*)path, (const char*)payload);
}

void
PCEventFeedDelete(char* path) {
    eventHandle(RPATH_EVENT_METHOD_DELETE, (const char*)path, NULL);
}

void
eventHandle(NSString* eventMethod, const char* path, const char* payload) {

    // parse in the background
    dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{

        NSString *eventPath = \
            [[NSString alloc]
                 initWithBytesNoCopy:(void *)path
                 length:strlen((const char*)path)
                 encoding:NSUTF8StringEncoding
                 freeWhenDone:YES];

        NSDictionary* eventPayload = nil;
        if (payload != NULL) {

            NSData *payloadData = \
                [[NSData alloc]
                 initWithBytesNoCopy:(void *)payload
                 length:strlen((const char*)payload)
                 freeWhenDone:YES];

            NSError *error = nil;
            eventPayload = \
                [NSJSONSerialization
                     JSONObjectWithData:payloadData
                     options:NSJSONReadingMutableContainers
                     error:&error];
            if (error != nil) {
                Log(@"%@", [error description]);
                return;
            }
        }

        dispatch_async(dispatch_get_main_queue(), ^{
            [[AppDelegate sharedDelegate]
             HandleEventForMethod:eventMethod
             onPath:eventPath
             withPayload:eventPayload];
        });
    });
}

@implementation AppDelegate (EventHandle)

- (void)HandleEventForMethod:(NSString *)aMethod onPath:(NSString *)aPath withPayload:(NSDictionary *)aPayload {
    Log(@"%@ -> %@ | %@", aMethod, aPath, aPayload);

/*
    NSAlert *alert = [[NSAlert alloc] init];
    [alert addButtonWithTitle:@"OK"];
    [alert addButtonWithTitle:@"Cancel"];
    [alert setMessageText:[message valueForKey:@"feed_msg"]];
    [alert setInformativeText:@"Deleted records cannot be restored."];
    [alert setAlertStyle:NSWarningAlertStyle];
    [alert runModal];
*/
}

@end
