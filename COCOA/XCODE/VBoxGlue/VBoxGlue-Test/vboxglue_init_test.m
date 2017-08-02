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
    NSString *userhome;
}
@end

@implementation vboxglue_init_test

- (void)setUp {
    [super setUp];
    userhome = NSHomeDirectory();
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
    XCTAssertTrue( 5001022 <= VBoxAppVersion(), @"Virtualbox version should be greater than or equal to 5001022");
    NSLog(@"API Version %d", VBoxApiVersion());
    XCTAssertTrue( 5001 <= VBoxApiVersion(), @"VBox API version should be greater than or equal to  5001");
}

-(void)test_Error_Message {
    XCTAssertTrue( VBGlue_Fail == VBoxTestErrorMessage(vboxGlue), @"Should generate error");
    NSLog(@"Error Message %s", VBoxGetErrorMessage(vboxGlue));
}

- (void)testCreateAndReleaseMachine {
    XCTAssertTrue( VBGlue_Ok == VBoxCreateMachineByName(vboxGlue, "pc-master-1", [userhome UTF8String]), @"Machine creation should return ok");
    NSLog(@"setting file path %s", VboxGetSettingFilePath(vboxGlue));
    
    // release machine
    XCTAssertTrue( VBGlue_Ok == VBoxReleaseMachine(vboxGlue), @"Releasing machine should not generate error");
}

- (void)test_Create_Find_Release_Machine {
    XCTAssertTrue( VBGlue_Ok == VBoxCreateMachineByName(vboxGlue, "pc-master-2", [userhome UTF8String]), @"Machine creation should return true");
    NSLog(@"setting file path %s", VboxGetSettingFilePath(vboxGlue));
    NSLog(@"MachineID %s", VboxGetMachineID(vboxGlue));
    
    // release machine
    XCTAssertTrue( VBGlue_Ok == VBoxReleaseMachine(vboxGlue), @"Releasing machine should not generate error");
}

#if 0
- (void)testFindMachineBeforeRegistration {
    __autoreleasing NSError *error = nil;
    XCTAssertFalse([self.vboxInterface acquireMachineByNameOrID:@"89be88f7-fc05-4aed-b3c7-9cf553be16a4" error:&error]);
    XCTAssertNotNil(error, @"Finding machine id should not generate any error");
    NSLog(@"%@", [error localizedFailureReason]);
}
#endif

@end