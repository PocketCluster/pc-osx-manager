//
//  vboxglue_init_test.m
//  VBoxGlue
//
//  Created by Almighty Kim on 8/2/17.
//  Copyright © 2017 PocketCluster. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <XCTest/XCTest.h>
#include "libvboxcom.h"

@interface vboxglue_init_test : XCTestCase {
    VBoxGlue *vboxGlue;
}
@end

@implementation vboxglue_init_test

- (void)setUp {
    [super setUp];
    VBGlueResult ret = NewVBoxGlue(&vboxGlue);
    if ( ret != VBGlue_Ok ) {
        NSLog(@"error!!!");
    }
}

- (void)tearDown {
    VBGlueResult ret = CloseVBoxGlue(vboxGlue);
    if ( ret != VBGlue_Ok ) {
        NSLog(@"error!!!");
    }
    [super tearDown];
}

-(void)test_App_API_Version {
    NSLog(@"App Version %d", VBoxAppVersion());
    XCTAssertTrue(5001022 <= VBoxAppVersion(), @"Virtualbox version should be greater than or equal to 5001022");
    NSLog(@"API Version %d", VBoxApiVersion());
    XCTAssertTrue(5001 <= VBoxApiVersion(), @"VBox API version should be greater than or equal to  5001");
}

-(void)test_Error_Message {
    XCTAssertTrue(VBGlue_Fail == VBoxTestErrorMessage(vboxGlue), @"Should generate error");
    NSLog(@"Error Message %s", VBoxGetErrorMessage(vboxGlue));
}


#if 0
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
#endif

@end