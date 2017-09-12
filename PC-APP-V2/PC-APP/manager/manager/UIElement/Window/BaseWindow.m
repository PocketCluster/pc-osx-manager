//
//  BaseWindow.m
//  manager
//
//  Created by Almighty Kim on 10/26/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "BaseWindow.h"

@implementation BaseWindow
-(instancetype)initWithNibName:(NSString *)aNibName bundle:(NSBundle *)aNibBundle {
    @autoreleasepool {
        __autoreleasing NSArray *nibContent;

        NSString *nibName = (aNibName)?aNibName:[self className];
        NSBundle *bundle = (aNibBundle)?aNibBundle:[NSBundle mainBundle];

        [bundle
         loadNibNamed:nibName
         owner:self
         topLevelObjects:&nibContent];
        
        for(id obj in nibContent){
            if ([obj isKindOfClass:[self class]]){
                self = nil;
                self = obj;
                self.delegate = self;
                break;
            }
        }
    }

    if(self != nil) {
        [self finishConstruction];
    }
    return self;
}

- (void)finishConstruction {
}
@end