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
    static const char* app_support_path = "/Users/almightykim/Library/Application Support/SysUtil";
    XCTAssert(strcmp(app_support_path, PCApplicationSupportDirectory()) == 0, @"Application Support path is incorrect");

    static const char* app_document_path = "/Users/almightykim/Documents";
    XCTAssert(strcmp(app_document_path, PCApplicationDocumentsDirectory()) == 0, @"Application Document path is incorrect");
    
    NSLog(@"%s", PCApplicationTemporaryDirectory());
    
    NSLog(@"%s", PCApplicationLibraryCacheDirectory());
    
    XCTAssert([[NSString stringWithUTF8String:PCApplicationResourceDirectory()] hasSuffix:@"SysUtil.app/Contents/Resources"], @"Application Resource path is incorrect");
    
    XCTAssert([[NSString stringWithUTF8String:PCApplicationExecutableDirectory()] hasSuffix:@"/SysUtil"], @"Applicatiopn Executable path is incorrect");
}

@end
