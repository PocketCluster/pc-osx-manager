//
//  AppDelegate.h
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

@class NativeMenu;
@class BaseWindowController;

#import "PCConstants.h"

@interface AppDelegate : NSObject <NSApplicationDelegate> {
    NSMutableArray<BaseWindowController *>* _openWindows;
}
@property (nonatomic, strong) NSMutableArray<BaseWindowController *>* openWindows;

+ (AppDelegate*)sharedDelegate;
@end

