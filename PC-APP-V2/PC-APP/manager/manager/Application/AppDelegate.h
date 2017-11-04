//
//  AppDelegate.h
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

@class NativeMenu;
@class BaseWindowController;

@interface AppDelegate : NSObject <NSApplicationDelegate> {
    NSMutableArray<BaseWindowController *>* _openWindows;

@private
    BOOL _isSystemReady;
    BOOL _isAppExpired;
    BOOL _isFirstTime;
    BOOL _isUserAuthed;
}
@property (nonatomic, strong) NSMutableArray<BaseWindowController *>* openWindows;
@property (nonatomic, strong, readonly) NativeMenu *mainMenu;

@property (nonatomic, readonly) BOOL isSystemReady;
@property (nonatomic, readonly) BOOL isAppExpired;
@property (nonatomic, readonly) BOOL isFirstTime;
@property (nonatomic, readonly) BOOL isUserAuthed;

+ (AppDelegate*)sharedDelegate;
@end

