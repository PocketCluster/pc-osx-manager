//
//  BaseWindow.h
//  manager
//
//  Created by Almighty Kim on 10/26/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+Window.h"

@class BaseWindow;

@interface BaseWindow : NSWindow <NSWindowDelegate>
-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil;
- (void)finishConstruction;
@end
