//
//  AboutWindowController.m
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#import "AboutWindow.h"

@implementation AboutWindow

- (id)initWithWindow:(NSWindow *)window {
    self = [super initWithWindow:window];
    return self;
}

- (void)windowDidLoad {
    [super windowDidLoad];

    [self.window setBackgroundColor:[NSColor whiteColor]];
    [self.window setTitleVisibility:NSWindowTitleHidden];
    [self.window setTitlebarAppearsTransparent:YES];
    self.window.styleMask |= NSFullSizeContentViewWindowMask;

    NSString *str = @"Version {VERSION} - Early Evaluation\n\nCopyright ©2015,{YEAR} PocketCluster";
    NSString *dateString = [NSString stringWithCString:__DATE__ encoding:NSASCIIStringEncoding];
    NSString *yearString = [dateString substringWithRange:NSMakeRange([dateString length] - 4, 4)];

    str = [str stringByReplacingOccurrencesOfString:@"{YEAR}" withString:yearString];
    str = [str stringByReplacingOccurrencesOfString:@"{VERSION}" withString:[[[NSBundle mainBundle] infoDictionary] valueForKey:@"CFBundleShortVersionString"]];
    [self.copyright setStringValue:str];
}

-(IBAction)homepage:(id)sender {
     [[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:@"https://github.com/pocketcluster/pocketcluster"]];
}
@end
