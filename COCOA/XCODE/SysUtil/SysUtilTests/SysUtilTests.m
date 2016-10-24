//
//  SysUtilTests.m
//  SysUtilTests
//
//  Created by Almighty Kim on 10/24/16.
//  Copyright © 2016 PocketCluster. All rights reserved.
//

#import <XCTest/XCTest.h>
#import <string.h>
#import "NSResourcePath.h"
#import "PCDeviceSerial.h"
#import "PCUserEnvironment.h"

@interface SysUtilTests : XCTestCase

@end

@implementation SysUtilTests

- (void)setUp {
    [super setUp];
}

- (void)tearDown {
    [super tearDown];
}

- (void)testPathes {
    [self setUp];
    
    static const char* app_support_path = "/Users/almightykim/Library/Application Support/SysUtil";
    XCTAssert(strcmp(app_support_path, PCApplicationSupportDirectory()) == 0, @"Application Support path is incorrect");

    static const char* app_document_path = "/Users/almightykim/Documents";
    XCTAssert(strcmp(app_document_path, PCApplicationDocumentsDirectory()) == 0, @"Application Document path is incorrect");
    
    NSLog(@"%s", PCApplicationTemporaryDirectory());
    
    NSLog(@"%s", PCApplicationLibraryCacheDirectory());
    
    XCTAssert([[NSString stringWithUTF8String:PCApplicationResourceDirectory()] hasSuffix:@"SysUtil.app/Contents/Resources"], @"Application Resource path is incorrect");
    
    XCTAssert([[NSString stringWithUTF8String:PCApplicationExecutableDirectory()] hasSuffix:@"/SysUtil"], @"Applicatiopn Executable path is incorrect");
    
    [self tearDown];
}

- (void)testDeviceSerial {
    [self setUp];

    NSLog(@"%s", PCDeviceSerialNumber());
    static const char* expected_device_serial = "G8815052XYL";
    XCTAssert(strcmp(PCDeviceSerialNumber(), expected_device_serial) == 0, @"Incorrect Serial number. Expecting G8815052XYL");
    
    [self tearDown];
}

- (void)testUserEnvironment {
    [self setUp];
    
    NSLog(@"%s", PCEnvironmentCocoaHomeDirectory());
    const char* home_directory = "/Users/almightykim";
    XCTAssert(strcmp(home_directory, PCEnvironmentCocoaHomeDirectory()) == 0, @"Incorrect Home directory");

    NSLog(@"%s", PCEnvironmentPosixHomeDirectory());
    const char* posix_home_directory = "/Users/almightykim";
    XCTAssert(strcmp(posix_home_directory, PCEnvironmentPosixHomeDirectory()) == 0, @"Incorrect Home directory");
    
    NSLog(@"%s", PCEnvironmentFullUserName());
    const char* full_name = "Almighty Kim";
    XCTAssert(strcmp(full_name, PCEnvironmentFullUserName()) == 0, @"Incorrect Full username");
    
    NSLog(@"%s", PCEnvironmentUserTemporaryDirectory());
    
    NSLog(@"%s", PCEnvironmentLoginUserName());
    const char* login_name = "almightykim";
    XCTAssert(strcmp(login_name, PCEnvironmentLoginUserName()) == 0, @"Incorrect login username");
    
    [self tearDown];
}

@end
