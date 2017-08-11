//
//  AppDelegate.h
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

@class NativeMenu;

#import "PCConstants.h"

@interface AppDelegate : NSObject <NSApplicationDelegate> {
    NSMutableArray* _openWindows;
}
@property (nonatomic, strong) NSMutableArray *openWindows;
@property (readonly) NativeMenu *nativeMenu;

+ (AppDelegate*)sharedDelegate;
@end

