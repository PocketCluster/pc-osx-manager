//
//  PCRouteTrie.h
//  manager
//
//  Created by Almighty Kim on 8/13/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <Foundation/Foundation.h>

typedef void (^ResponseHandler)(NSDictionary *payload);

@interface PCRouteTrie : NSObject
- (instancetype) initWithPathComponent:(NSString *)aComponent;
- (void) addNode:(NSString*)aMethod forPath:(NSString*)aPath withHandlerBlock:(ResponseHandler)aHandler;
- (void) traseverse:(NSArray<NSString *>*)components toNode:(PCRouteTrie **)toNode forComponent:(NSString **)toComponent;
@end
