//
//  managerTests.m
//  managerTests
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <XCTest/XCTest.h>

@interface managerTests : XCTestCase

@end

@implementation managerTests

- (void)setUp {
    [super setUp];
    // Put setup code here. This method is called before the invocation of each test method in the class.
}

- (void)tearDown {
    // Put teardown code here. This method is called after the invocation of each test method in the class.
    [super tearDown];
}

- (void)test_PCTrieNode {
    NSString *path = @"/v1/system/monitor/status/battery";
    NSArray *components = [path componentsSeparatedByString:@"/"];
    
    @try {
        [@[] objectAtIndex:0];
    }
    @catch (NSException *exception) {
        NSLog(@"empty[0] %@", [exception description]);
    }
    @finally {
    }
    
    NSLog(@"components(full) %@", components);
    NSLog(@"components[0, 1] %@", [components subarrayWithRange:NSMakeRange(0, 1)]);
    NSLog(@"components[1, 0] %@", [components subarrayWithRange:NSMakeRange(1, 0)]);
    NSLog(@"components[2, 0] %@", [components subarrayWithRange:NSMakeRange(2, 0)]);
    NSLog(@"components[1, 1] %@", [components subarrayWithRange:NSMakeRange(1, 1)]);
    NSLog(@"components[1, last(lenth - 1)] %@", [components subarrayWithRange:NSMakeRange(1, [components count] - 1)]);
    
    
    PCRouteTrie *trie = [[PCRouteTrie alloc] initWithPathComponent:@"/"];
    [trie addNode:@"GET" forPath:@"/v1/system/monitor/battery" withHandlerBlock:^(NSDictionary *payload) {
        NSLog(@"new Handler");
    }];
    
    PCRouteTrie *node = nil;
    NSString *component;
    [trie traseverse:[@"v1/system/monitor/battery" componentsSeparatedByString:@"/"] toNode:&node forComponent:&component];
    NSLog(@"node %@ component %@",[node description], component);
    [trie traseverse:[@"v1/system/monitor" componentsSeparatedByString:@"/"] toNode:&node forComponent:&component];
    NSLog(@"node %@ component %@",[node description], component);}

- (void)testPerformanceExample {
    // This is an example of a performance test case.
    [self measureBlock:^{
        // Put the code you want to measure the time of here.
    }];
}

@end
