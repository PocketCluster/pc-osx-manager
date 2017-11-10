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
handleResponse(NSString* method, const char* path, const char* response);

void
PCFeedResponseForGet(char* path, char* response) {
    handleResponse(RPATH_EVENT_METHOD_GET, (const char*)path, (const char*)response);
}

void
PCFeedResponseForPost(char* path, char* response) {
    handleResponse(RPATH_EVENT_METHOD_POST, (const char*)path, (const char*)response);
}

void
PCFeedResponseForPut(char* path, char* response) {
    handleResponse(RPATH_EVENT_METHOD_PUT, (const char*)path, (const char*)response);
}

void
PCFeedResponseForDelete(char* path, char* response) {
    handleResponse(RPATH_EVENT_METHOD_DELETE, (const char*)path, (const char*)response);
}

void
handleResponse(NSString* eventMethod, const char* path, const char* response) {

    // parse in the background
    dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        @autoreleasepool {
            NSString *eventPath = \
                [[NSString alloc]
                     initWithBytesNoCopy:(void *)path
                     length:strlen((const char*)path)
                     encoding:NSUTF8StringEncoding
                     freeWhenDone:YES];

            NSDictionary* eventResponse = nil;
            if (response != NULL) {

                NSData *payloadData = \
                    [[NSData alloc]
                     initWithBytesNoCopy:(void *)response
                     length:strlen((const char*)response)
                     freeWhenDone:YES];

                NSError *error = nil;
                eventResponse = \
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
                [[PCRouter sharedRouter]
                 responseFor:eventMethod
                 onPath:eventPath
                 withPayload:eventResponse];
            });
        }
    });
}

@implementation AppDelegate (ResponseHandle)

// (2017/08/16)
// We don't need to unnecessarily call another just to have singleton point on *main thread* !
// Nonetheless, without a valid method, this category might be eradicated by compiler.
// So... let's just leave it here, but it won't be used.
- (void)HandleResponseForMethod:(NSString *)aMethod onPath:(NSString *)aPath withPayload:(NSDictionary *)aResponse {
    [[PCRouter sharedRouter] responseFor:aMethod onPath:aPath withPayload:aResponse];
}

@end
