//
//  PCRouteTrie.m
//  manager
//
//  Created by Almighty Kim on 8/13/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

// test classes in managerTests

#import "PCRouteTrie.h"

@interface PCRouteTrie() {
    __strong NSMutableArray<PCRouteTrie* >* _children;
    __strong NSString* _component;
    __strong NSMutableDictionary<NSString*, ResponseHandler>* _methods;
}
@property (nonatomic, strong, readonly) NSMutableArray<PCRouteTrie* >* children;
@property (nonatomic, strong, readonly) NSString* component;
@property (nonatomic, strong, readonly) NSMutableDictionary<NSString*, ResponseHandler>* methods;
@end

@implementation PCRouteTrie
@synthesize children = _children;
@synthesize component = _component;
@synthesize methods = _methods;

- (instancetype) initWithPathComponent:(NSString *)aComponent {
    self = [super init];
    if (self != nil) {
        _children = [NSMutableArray new];
        _component = aComponent;
        _methods = [NSMutableDictionary new];
    }
    return self;
}

- (void)dealloc {
    _component = nil;
    
    [_children removeAllObjects];
    _children = nil;

    [_methods removeAllObjects];
    _methods = nil;
}

- (NSString *)description {
    return [NSString stringWithFormat:@"PCTrieNode.component : %@ child[%ld], methods[%ld]", [self component], [self.children count], [self.methods count]];
}

- (void) addNode:(NSString*)aMethod forPath:(NSString*)aPath withHandlerBlock:(ResponseHandler)aHandler {

    NSArray<NSString *>* components = [aPath componentsSeparatedByString:@"/"];
    components = [components subarrayWithRange:NSMakeRange(1, [components count] - 1)];

    for (NSUInteger count = [components count]; 0 < count;) {
        
        PCRouteTrie *aNode = nil;
        NSString *component = nil;
        
        [self traseverse:components toNode:&aNode forComponent:&component];
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

- (void) traseverse:(NSArray<NSString *>*)components toNode:(PCRouteTrie **)toNode forComponent:(NSString **)toComponent {

    NSString *component = [components objectAtIndex:0];

    // if no children, bailout
    if (0 < [self.children count]) {
        for (PCRouteTrie *child in self.children) {
            if ([component isEqualToString:child.component]) {
                NSArray<NSString *>* next = [components subarrayWithRange:NSMakeRange(1, [components count] - 1)];
                if (0 < [next count]) { // http://xkcd.com/1270/
                    return [child traseverse:next toNode:toNode forComponent:toComponent]; // tail recursion is it's own reward.
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

@end
