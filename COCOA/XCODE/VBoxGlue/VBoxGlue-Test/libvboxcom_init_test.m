//
//  libvboxcom_init_test.m
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/16/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//


#import <XCTest/XCTest.h>
#import "VboxCom.h"

@interface libvboxcom_init_test : XCTestCase
@property (nonatomic, strong) VboxCom* vboxInterface;
@end

@implementation libvboxcom_init_test

- (void)setUp {
    [super setUp];
    NSError *error;
    self.vboxInterface = [[VboxCom alloc] initWithError:&error];
    XCTAssertNil(error, @"Initialization should not generate any error");
}

- (void)tearDown {
    self.vboxInterface = nil;
    [super tearDown];
}

-(void)testAppVersion {
    XCTAssertTrue([self.vboxInterface checkAppVersion]);
}

@end
