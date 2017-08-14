//
//  PCRouteTrie.m
//  manager
//
//  Created by Almighty Kim on 8/13/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

// test classes in managerTests

#import "PCRouteTrie.h"

@interface PCRequestHolder()
@property (nonatomic, strong, readwrite) NSObject<PCRouteRequest>* request;
@property (nonatomic, strong, readwrite) ResponseHandler handler;
@end

@implementation PCRequestHolder
-(void) dealloc {
    self.request = nil;
    self.handler = nil;
}
@end


@interface PCRouteTrie() {
    __strong NSMutableArray<PCRouteTrie* >* _children;
    __strong NSString* _component;
    __strong NSMutableDictionary<NSString*, NSMutableArray<PCRequestHolder *>*>* _methods;
}
@property (nonatomic, strong, readonly) NSMutableArray<PCRouteTrie* >* children;
@property (nonatomic, strong, readonly) NSString* component;
@property (nonatomic, strong, readonly) NSMutableDictionary<NSString*, NSMutableArray<PCRequestHolder *>*>* methods;

// - (void) addNode:(NSString*)aMethod forPath:(NSString*)aPath withHandlerBlock:(ResponseHandler)aHandler;
- (PCRouteTrie *) _findOrAddNodeForPath:(NSString*)aPath;
- (void) _traseverse:(NSArray<NSString *>*)components toNode:(PCRouteTrie **)toNode forComponent:(NSString **)toComponent;
@end

@implementation PCRouteTrie
@synthesize children = _children;
@synthesize component = _component;
@synthesize methods = _methods;

- (instancetype) initWithPathComponent:(NSString *)aComponent {
    self = [super init];
    if (self != nil) {
        _children = [NSMutableArray<PCRouteTrie*> new];
        _component = aComponent;
        _methods = [NSMutableDictionary<NSString*, NSMutableArray<PCRequestHolder *>*> new];
    }
    return self;
}

- (void)dealloc {
    _component = nil;

    [_methods removeAllObjects];
    _methods = nil;
    
    [_children removeAllObjects];
    _children = nil;
}

- (NSString *)description {
    return [NSString stringWithFormat:@"PCTrieNode.component : %@ child[%ld], methods[%ld]", [self component], [self.children count], [self.methods count]];
}

#if 0
- (void) addNode:(NSString*)aMethod forPath:(NSString*)aPath withHandlerBlock:(ResponseHandler)aHandler {
    
    NSArray<NSString *>* components = [aPath componentsSeparatedByString:@"/"];
    components = [components subarrayWithRange:NSMakeRange(1, [components count] - 1)];
    
    for (NSUInteger count = [components count]; 0 < count;) {
        
        PCRouteTrie *aNode = nil;
        NSString *component = nil;
        
        [self _traseverse:components toNode:&aNode forComponent:&component];
        // update an existing node.
        if ([component isEqualToString:aNode.component] && count == 1) {
            [aNode.methods setValue:aHandler forKey:aMethod];
            return;
        }
        PCRouteTrie *newNode = [[PCRouteTrie alloc] initWithPathComponent:component];
        
        // this is the last component of the url resource, so it gets the handler.
        if (count == 1) {
            [newNode.methods setValue:aHandler forKey:aMethod];
        }
        [aNode.children addObject:newNode];
        --count;
    }
}
#endif

- (PCRouteTrie *) _findOrAddNodeForPath:(NSString*)aPath {

    PCRouteTrie *target = nil;
    NSArray<NSString *>* components = [aPath componentsSeparatedByString:@"/"];
    components = [components subarrayWithRange:NSMakeRange(1, [components count] - 1)];

    for (NSUInteger count = [components count]; 0 < count;) {
        
        PCRouteTrie *aNode = nil;
        NSString *component = nil;
        
        [self _traseverse:components toNode:&aNode forComponent:&component];

        // found a candiate node
        if ([component isEqualToString:aNode.component] && count == 1) {
            target = aNode;
            break;
        }
        PCRouteTrie *newNode = [[PCRouteTrie alloc] initWithPathComponent:component];
        [aNode.children addObject:newNode];
        
        // this is the last component of the url resource, so it is the target node
        if (count == 1) {
            target = newNode;
            break;
        }

        count--;
    }

#ifdef DEBUG
    if (target == nil) {
        Log(@"PCRouteTrie [%s] -> %@ target node is null. This should never happen! ", __PRETTY_FUNCTION__, aPath);
    }
#endif
    return target;
}

- (void) _traseverse:(NSArray<NSString *>*)components toNode:(PCRouteTrie **)toNode forComponent:(NSString **)toComponent {

    NSString *component = [components objectAtIndex:0];

    // if no children, bailout
    if (0 < [self.children count]) {
        for (PCRouteTrie *child in self.children) {
            if ([component isEqualToString:child.component]) {
                NSArray<NSString *>* next = [components subarrayWithRange:NSMakeRange(1, [components count] - 1)];
                if (0 < [next count]) { // http://xkcd.com/1270/
                    return [child _traseverse:next toNode:toNode forComponent:toComponent]; // tail recursion is it's own reward.
                }
                *toNode = child;
                *toComponent = component;
                return;
            }
        }
    }
    *toNode = self;
    *toComponent = component;
    return;
}

- (void) addRequest:(NSObject<PCRouteRequest> *)aRequest forMethod:(NSString *)aMethod onPath:(NSString *)aPath withHandler:(ResponseHandler)aHandler {
    PCRouteTrie *node = [self _findOrAddNodeForPath:aPath];
    NSMutableArray<PCRequestHolder *> *reqList = [node.methods objectForKey:aMethod];
    if (reqList == nil) {
        reqList = [NSMutableArray<PCRequestHolder *> new];
        [node.methods setValue:reqList forKey:aMethod];
    }
    PCRequestHolder *holder = [PCRequestHolder new];
    [holder setRequest:aRequest];
    [holder setHandler:aHandler];
    [reqList addObject:holder];
}

- (void) delRequest:(NSObject<PCRouteRequest> *)aRequest forMethod:(NSString *)aMethod onPath:(NSString *)aPath {
    PCRouteTrie *node = [self _findOrAddNodeForPath:aPath];
    NSMutableArray<PCRequestHolder *> *reqList = [node.methods objectForKey:aMethod];
    // nothing to delete
    if (reqList == nil) {
        return;
    }
    
    __weak PCRequestHolder *holder = nil;
    for (PCRequestHolder *h in reqList) {
        if (h.request == aRequest) {
            holder = h;
            break;
        }
    }

    if (holder != nil) {
        Log(@"b4 removal length %ld", [reqList count]);
        [reqList removeObject:holder];
        Log(@"a4 removal length %ld", [reqList count]);
    }
}

// this always returns the last object
- (PCRequestHolder *)findRequestForMethod:(NSString *)aMethod onPath:(NSString *)aPath {
    PCRouteTrie *node = [self _findOrAddNodeForPath:aPath];
    NSMutableArray<PCRequestHolder *> *reqList = [node.methods objectForKey:aMethod];
    return [reqList lastObject];
}

- (NSArray<PCRequestHolder *> *)findAllRequestForMethod:(NSString*)aMethod onPath:(NSString*)aPath {
    PCRouteTrie *node = [self _findOrAddNodeForPath:aPath];
    return [node.methods objectForKey:aMethod];
}

@end
