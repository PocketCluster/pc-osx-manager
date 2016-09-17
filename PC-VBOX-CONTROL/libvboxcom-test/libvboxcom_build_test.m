//
//  libvboxcom_test.m
//  libvboxcom-test
//
//  Created by Almighty Kim on 7/12/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#import <XCTest/XCTest.h>
#import "VboxCom.h"

@interface libvboxcom_build_test : XCTestCase
@property (nonatomic, strong) VboxCom* vboxInterface;
@end

@implementation libvboxcom_build_test

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

-(void)testBuildMachine {
    __autoreleasing NSError *error = nil;
    XCTAssertTrue([self.vboxInterface createMachineWithName:@"pc-master" error:&error],@"Machine creation should return true");
    XCTAssertNil(error, @"Machine creation should not generate any error");

    BOOL result = [self.vboxInterface
                   buildMachineWithCPU:2
                   memorySize:2048
                   hostInterface:@"en4"
                   sharedFolderName:@"/pocket"
                   sharedFolderPath:@"/Users/almightykim/Workspace"
                   bootImagePath:@"/Users/almightykim/Downloads/boot2docker.iso"
                   hardDiskSize:(1 << 26)
                   progress:^(int progress) {
                       NSLog(@"progress %d", progress);
                   } error:&error];

    NSLog(@"error message %@", [error localizedDescription]);
    XCTAssertTrue(result);
    XCTAssertNil(error, @"Machine building should not generate any error");

    NSString* machineID = [self.vboxInterface retrieveMachineId:&error];
    XCTAssertNotNil(machineID, @"machineID cannot be empty after creation");
    NSLog(@"ID : %@", machineID);
    
    XCTAssertFalse([self.vboxInterface isMachineSettingChanged:&error]);
    XCTAssertNil(error, @"Checking machine setting change should not generate any error");
}

-(void)testDestoryMachine {
    __autoreleasing NSError *error = nil;
    XCTAssertTrue([self.vboxInterface acquireMachineByNameOrID:@"pc-master" error:&error] ,@"Machine retrieval should return true");
    XCTAssertNil(error, @"Machine retrieval should not generate any error");

    XCTAssertFalse([self.vboxInterface isMachineSettingChanged:&error]);
    XCTAssertNil(error, @"Checking machine setting change should not generate any error");
    
    XCTAssertTrue([self.vboxInterface destoryCurrentMachine:&error]);
    XCTAssertNil(error, @"Machine destruction should not generate any error");    
}

#ifdef CHECK_PERFORMANCE
-(void)testPerformanceExample {
    // This is an example of a performance test case.
    [self measureBlock:^{
        // Put the code you want to measure the time of here.
    }];
}
#endif

@end
