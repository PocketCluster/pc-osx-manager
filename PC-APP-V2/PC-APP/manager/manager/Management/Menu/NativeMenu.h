//
//  NativeMenu.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "Util.h"
#import "AboutWindow.h"

@interface NativeMenu : NSObject <NSMenuDelegate>{
@private
    AboutWindow      *_aboutWindow;
    NSStatusItem     *_statusItem;
}
@property (nonatomic, strong) AboutWindow *aboutWindow;
@property (nonatomic, strong) NSStatusItem *statusItem;
@end
