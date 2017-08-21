//
//  BaseBrandView.m
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "BaseBrandView.h"

@interface BaseBrandView()
-(void) finishConstruction;
@end

@implementation BaseBrandView 
@synthesize imageView = _imageView;
@synthesize contentBox = _contentBox;

- (instancetype)initWithFrame:(NSRect)frameRect {
    self = [super initWithFrame:DEFAULT_BRANDVIEW_FRAME];
    if (self != nil) {
        [self finishConstruction];
    }
    return self;
}

- (nullable instancetype)initWithCoder:(NSCoder *)coder {
    self = [super initWithCoder:coder];
    if (self != nil) {
        [self finishConstruction];
    }
    return self;
}

-(void) finishConstruction {
    [self setFrame:DEFAULT_BRANDVIEW_FRAME];
    
    self.imageView = [[NSImageView alloc] initWithFrame:NSMakeRect(-40, 60, 320, 320)];
    [_imageView setAlphaValue:0.3];
    [_imageView setImageScaling:NSImageScaleProportionallyUpOrDown];
    
    // 'applicationIconImage' doesn't bring the right one. need to fix
    //[_imageView setImage:[[NSApplication sharedApplication] applicationIconImage]];
    [_imageView setImage:[NSImage imageNamed:@"AppIcon"]];
    [self addSubview:_imageView];

    self.contentBox = [[NSBox alloc] initWithFrame:DEFAULT_BRAND_BOX_FRAME];
    [_contentBox setTitlePosition:(NSNoTitle)];
    [self addSubview:_contentBox];
}

-(void)dealloc {
    if (_imageView != nil) {
        [self.imageView removeFromSuperview];
        self.imageView = nil;
    }

    if (_contentBox != nil) {
        [self.contentBox removeFromSuperview];
        self.contentBox = nil;
    }
}

@end
