//
//  vboxglue_init_test.m
//  VBoxGlue
//
//  Created by Almighty Kim on 8/2/17.
//  Copyright © 2017 PocketCluster. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <XCTest/XCTest.h>
#import <string.h>
#include "libvboxcom.h"

static const char* TARGET_MACHINE_NAME = "POCKET_VBOX_TEST";

@interface vboxglue_init_test : XCTestCase {
    VBoxGlue vboxGlue;
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
    XCTAssertTrue( VBGlue_Fail == VBoxTestErrorMessage(vboxGlue),               @"Should generate error");
    NSLog(@"Error Message %s", VBoxGetErrorMessage(vboxGlue));
}

- (void)test_Create_Release_Machine {
    static const char* TEST_MACHINE_NAME = "TEST_MACHINE_NAME";

    XCTAssertTrue( VBGlue_Ok == VBoxCreateMachineByName(vboxGlue,
                                                        [userhome UTF8String],
                                                        TEST_MACHINE_NAME),     @"Machine creation should return true");
    NSLog(@"setting file path %s", VboxGetSettingFilePath(vboxGlue));
    NSLog(@"MachineID Origin %s", VboxGetMachineID(vboxGlue));
    

    XCTAssertTrue( VBGlue_Ok == VBoxReleaseMachine(vboxGlue),                   @"Releasing machine should not generate error");
}

-(void)test_Build_machine {
    VBoxBuildOption *options =
        VBoxMakeBuildOption(2, 2048,
                            "en1",
                            "/tmp", "/temp",
                            "/Users/almightykim/Workspace/VBOX-IMAGE/pc-core.iso",
                            "/Users/almightykim/Workspace/VBOX-IMAGE/pc-core-hdd.vmdk");
 
    XCTAssertTrue( VBGlue_Ok == VBoxCreateMachineByName(vboxGlue,
                                                        [userhome UTF8String],
                                                        TARGET_MACHINE_NAME),   @"Machine creation should return true");
    XCTAssertTrue( VBGlue_Ok == VBoxBuildMachine(vboxGlue, options),            @"Machine building should return true");
    NSLog(@"setting file path %s", VboxGetSettingFilePath(vboxGlue));
    NSLog(@"MachineID Origin %s", VboxGetMachineID(vboxGlue));
    XCTAssertTrue( strlen(VboxGetMachineID(vboxGlue)) != 0,                     @"Machine id should not be zero");
}

-(void)test_Destory_Machine {
    bool isChanged = true;
    XCTAssertTrue( VBGlue_Ok == VBoxFindMachineByNameOrID(vboxGlue,
                                                          TARGET_MACHINE_NAME), @"Should be able to find existing machine return true");
    XCTAssertTrue( VBGlue_Ok == VBoxIsMachineSettingChanged(vboxGlue, &isChanged),@"Machine setting checking should return true");
    XCTAssertTrue( !isChanged,                                                  @"Machine setting should have not changed");
    XCTAssertTrue( VBGlue_Ok == VBoxDestoryMachine(vboxGlue),                   @"Machine Destruction should go fine");
    XCTAssertTrue( VBGlue_Ok == VBoxReleaseMachine(vboxGlue),                   @"Releasing machine should not generate error");
}

@end