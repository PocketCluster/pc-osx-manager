//
//  BookmarkManager.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "Raspberry.h"

@interface RaspberryManager : NSObject {
    NSMutableArray *_raspberries;
}

+ (RaspberryManager *)sharedManager;

- (void)loadRaspberries;
- (void)saveRaspberries;
- (void)clearRaspberries;

- (NSMutableArray<Raspberry *> *)getRaspberries;
- (Raspberry *) addRaspberry:(Raspberry *)aRaspberry;


@end
