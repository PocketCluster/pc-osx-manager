//
//  NSResourcePath.m
//  SystemUtil
//
//  Created by Almighty Kim on 10/24/16.
//  Copyright © 2016 PocketCluster. All rights reserved.
//

#import "NSResourcePath.h"
#import "NSFileManager+DirectoryLocations.h"
#import "NSBundle+DirectoryPath.h"

const char*
PCApplicationSupportDirectory(void) {
    return [[[NSFileManager defaultManager] applicationSupportDirectory] UTF8String];
}

const char*
PCApplicationDocumentsDirectory(void) {
    return [[[NSFileManager defaultManager] applicationDocumentsDirectory] UTF8String];
}

const char*
PCApplicationTemporaryDirectory(void) {
    return [[[NSFileManager defaultManager] applicationTemporaryDirectory] UTF8String];
}

const char*
PCApplicationLibraryCacheDirectory(void) {
    return [[[NSFileManager defaultManager] applicationLibraryCacheDirectory] UTF8String];
}

const char*
PCApplicationResourceDirectory(void) {
    return [[[NSBundle mainBundle] resourcePath] UTF8String];
}

const char*
PCApplicationExecutableDirectory(void) {
    return [[[NSBundle mainBundle] executablePath] UTF8String];
}
