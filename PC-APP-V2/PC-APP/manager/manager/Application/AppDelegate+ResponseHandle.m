//
//  AppDelegate+EventHandle.m
//  manager
//
//  Created by Almighty Kim on 3/24/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+ResponseHandle.h"
#import "PCRoutePathConst.h"
#import "PCRouter.h"

/*
 * Event message from engine is most likely a feedback for api call
 * The event message then should be processed and delivered to RunLoop in
 * Default mode to display sucess and failure message to UI in seqeunce
 */

static void
handleResponse(NSString* method, const char* path, const char* payload);

void
PCFeedResponseForGet(char* path, char* payload) {
    handleResponse(RPATH_EVENT_METHOD_GET, (const char*)path, (const char*)payload);
}

void
PCFeedResponseForPost(char* path, char* payload) {
    handleResponse(RPATH_EVENT_METHOD_POST, (const char*)path, (const char*)payload);
}

void
PCFeedResponseForPut(char* path, char* payload) {
    handleResponse(RPATH_EVENT_METHOD_PUT, (const char*)path, (const char*)payload);
}

void
PCFeedResponseForDelete(char* path, char* payload) {
    handleResponse(RPATH_EVENT_METHOD_DELETE, (const char*)path, (const char*)payload);
}

void
handleResponse(NSString* eventMethod, const char* path, const char* payload) {

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
             HandleResponseForMethod:eventMethod
             onPath:eventPath
             withPayload:eventPayload];
        });
    });
}

@implementation AppDelegate (ResponseHandle)

- (void)HandleResponseForMethod:(NSString *)aMethod onPath:(NSString *)aPath withPayload:(NSDictionary *)aPayload {
    [[PCRouter sharedRouter] responseFor:aMethod onPath:aPath withPayload:aPayload];
}

@end
