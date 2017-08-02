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
    VBRESULT ret = NewVBoxGlue(&vboxGlue);
    if ( ret != GOOD ) {
        NSLog(@"error!!!");
    }
}

- (void)tearDown {
    VBRESULT ret = CloseVBoxGlue(vboxGlue);
    if ( ret != GOOD ) {
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
@end