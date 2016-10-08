//
//  libvboxcom_create_test.m
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/16/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#import <XCTest/XCTest.h>
#import "VboxCom.h"

#import <XCTest/XCTest.h>
#import "VboxCom.h"

@interface libvboxcom_create_test : XCTestCase
@property (nonatomic, strong) VboxCom* vboxInterface;
@end

@implementation libvboxcom_create_test

- (void)setUp {
    [super setUp];
    NSError *error;
    self.vboxInterface = [[VboxCom alloc] initWithError:&error];
    XCTAssertNil(error, @"Initialization should not generate any error");
    [_vboxInterface openSession:&error];
    XCTAssertNil(error, @"Openning session should not generate any error");
}

- (void)tearDown {
    NSError *error;
    [_vboxInterface closeSession:&error];
    XCTAssertNil(error, @"Terminating session should not generate any error");
    self.vboxInterface = nil;
    [super tearDown];
}

- (void)testCreateAndReleaseMachine {
    __autoreleasing NSError *error = nil;
    XCTAssertTrue([self.vboxInterface createMachineWithName:@"pc-master-1" error:&error],@"Machine creation should return true");
    XCTAssertNil(error, @"Machine creation should not generate any error");
    
    // reelase machine
    [self.vboxInterface releaseMachine:&error];
    XCTAssertNil(error, @"Releasing machine should not generate error");
}

- (void)testCreateAndFindMachine {
    __autoreleasing NSError *error = nil;
    XCTAssertTrue([self.vboxInterface createMachineWithName:@"pc-master-2" error:&error],@"Machine creation should return true");
    XCTAssertNil(error, @"Machine creation should not generate any error");
    
    NSString *machineID = [self.vboxInterface retrieveMachineId:&error];
    XCTAssertNotNil(machineID, @"Machine id should not be empty");
    XCTAssertNil(error, @"Finding machine id should not generate any error");
    NSLog(@"MachineID %@", machineID);
    
    // relase machine
    [self.vboxInterface releaseMachine:&error];
    XCTAssertNil(error, @"Releasing machine should not generate error");
}

- (void)testFindMachineBeforeRegistration {
    __autoreleasing NSError *error = nil;
    XCTAssertFalse([self.vboxInterface acquireMachineByNameOrID:@"89be88f7-fc05-4aed-b3c7-9cf553be16a4" error:&error]);
    XCTAssertNotNil(error, @"Finding machine id should not generate any error");
    NSLog(@"%@", [error localizedFailureReason]);
}

@end
