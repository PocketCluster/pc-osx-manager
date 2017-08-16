//
//  BaseBrandView.h
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import <Cocoa/Cocoa.h>

#define DEFAULT_BRANDVIEW_FRAME ((NSRect){{0,0},{580,420}})

@interface BaseBrandView : NSView {
    NSImageView *_imageView;
    NSBox *_contentBox;
}
@property (nonatomic, strong) NSImageView *imageView;
@property (nonatomic, strong) NSBox *contentBox;
@end
