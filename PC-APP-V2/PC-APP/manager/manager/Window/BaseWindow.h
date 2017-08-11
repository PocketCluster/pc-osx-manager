//
//  BaseWindow.h
//  manager
//
//  Created by Almighty Kim on 10/26/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+Window.h"

@class BaseWindow;

@protocol PCWindowDelegate <NSWindowDelegate>
@end

@interface BaseWindow : NSWindow <NSWindowDelegate>
@property (weak, nonatomic) id<PCWindowDelegate> actionDelegate;
-(instancetype)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil;
- (void)finishConstruction;
@end
