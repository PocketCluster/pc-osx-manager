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

-(void)test_Search_Host_Interface {
    char* fullName = NULL;
    const char* foundName = "en1: Wi-Fi (AirPort)";
    VBGlueResult result;
    
    result = VBoxHostSearchNetworkInterfaceByName(vboxGlue, "en1", &fullName);
    XCTAssertTrue( VBGlue_Ok == result, @"find network address should not generate error");
    XCTAssertTrue( strcmp(foundName, fullName) == 0, @"full name should match");
    
    if (result == VBGlue_Fail) {
        NSLog(@"Failure reason %s", VBoxGetErrorMessage(vboxGlue));
    } else {
        NSLog(@"Full network interface name | %s", fullName);
    }
    free(fullName);
}

-(void)test_Error_Message {
    VBGlueResult result = VBoxTestErrorMessage(vboxGlue);
    XCTAssertTrue( VBGlue_Fail == result, @"Should generate error");
    NSLog(@"Error Message %s", VBoxGetErrorMessage(vboxGlue));
}

- (void)test_Create_Release_Machine {
    static const char* TEST_MACHINE_NAME = "TEST_MACHINE_NAME";

    VBGlueResult result = VBoxMachineCreateByName(vboxGlue, [userhome UTF8String], TEST_MACHINE_NAME);
    XCTAssertTrue( VBGlue_Ok == result, @"Machine creation should return true");
    NSLog(@"setting file path %s", VBoxGetSettingFilePath(vboxGlue));
    NSLog(@"MachineID Origin %s", VBoxGetMachineID(vboxGlue));
    
    result = VBoxMachineRelease(vboxGlue);
    XCTAssertTrue( VBGlue_Ok == result, @"Releasing machine should not generate error");
}

-(void)test_Build_machine {
    VBoxBuildOption *options =
        VBoxMakeBuildOption(2, 2048,
                            "en1: Wi-Fi (AirPort)",
                            "/Users/almightykim/Workspace/VBOX-IMAGE/pc-core.iso",
                            "/Users/almightykim/Workspace/VBOX-IMAGE/pc-core-hdd.vmdk",
                            "/tmp", "/temp");

    VBGlueResult result = VBoxMachineCreateByName(vboxGlue, [userhome UTF8String], TARGET_MACHINE_NAME);
    XCTAssertTrue( VBGlue_Ok == result, @"Machine creation should return true");
    
    result = VBoxMachineBuildWithOption(vboxGlue, options);
    XCTAssertTrue( VBGlue_Ok == result, @"Machine building should return true");
    if (result == VBGlue_Ok) {
        NSLog(@"Setting file path %s", VBoxGetSettingFilePath(vboxGlue));
        NSLog(@"MachineID Origin %s", VBoxGetMachineID(vboxGlue));
        XCTAssertTrue(VBGlueMachine_PoweredOff == VBoxMachineGetCurrentState(vboxGlue));
    } else {
        NSLog(@"Failed reason %s", VBoxGetErrorMessage(vboxGlue));
    }
 
    size_t len = strlen(VBoxGetMachineID(vboxGlue));
    XCTAssertTrue( len != 0, @"Machine id should not be zero");
}

// we need to preserve session to start & stop
-(void)test_Start_Stop_Machine {
    VBGlueResult result = VBoxMachineFindByNameOrID(vboxGlue, TARGET_MACHINE_NAME);
    XCTAssertTrue( VBGlue_Ok == result, @"find machine should return true");
    if (result == VBGlue_Ok) {
        NSLog(@"Setting file path %s", VBoxGetSettingFilePath(vboxGlue));
        NSLog(@"MachineID Origin %s", VBoxGetMachineID(vboxGlue));
        XCTAssertTrue(VBGlueMachine_PoweredOff == VBoxMachineGetCurrentState(vboxGlue));
    } else {
        NSLog(@"Failed reason %s", VBoxGetErrorMessage(vboxGlue));
        return;
    }

    result = VBoxMachineHeadlessStart(vboxGlue);
    XCTAssertTrue( VBGlue_Ok == result, @"machine start should return true");
    if (result == VBGlue_Ok) {
        XCTAssertTrue(VBGlueMachine_Running == VBoxMachineGetCurrentState(vboxGlue));
    } else {
        NSLog(@"Failed reason %s", VBoxGetErrorMessage(vboxGlue));
    }
   
    sleep(30);

    result = VBoxMachineAcpiDown(vboxGlue);
    XCTAssertTrue( VBGlue_Ok == result, @"machine stop should return true");
    if (result == VBGlue_Ok) {
        XCTAssertTrue(VBGlueMachine_PoweredOff == VBoxMachineGetCurrentState(vboxGlue));
    } else {
        NSLog(@"Failed reason %s", VBoxGetErrorMessage(vboxGlue));
    }
}

-(void)test_Destory_Machine {
    VBGlueResult result = VBoxMachineFindByNameOrID(vboxGlue, TARGET_MACHINE_NAME);
    XCTAssertTrue( VBGlue_Ok == result, @"Should be able to find existing machine return true");
    
    bool isChanged = true;
    result = VBoxMachineIsSettingChanged(vboxGlue, &isChanged);
    XCTAssertTrue( VBGlue_Ok == result, @"Machine setting checking should return true");
    XCTAssertTrue( !isChanged, @"Machine setting should have not changed");
    
    result = VBoxMachineDestory(vboxGlue);
    XCTAssertTrue( VBGlue_Ok == result, @"Machine Destruction should go fine");
    
    result = VBoxMachineRelease(vboxGlue);
    XCTAssertTrue( VBGlue_Ok == result, @"Releasing machine should not generate error");
}

@end