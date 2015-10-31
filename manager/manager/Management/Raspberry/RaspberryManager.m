//
//  RaspberryManager.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//
#import "VagrantManager.h"
#import "SynthesizeSingleton.h"
#import "RaspberryManager.h"

NSString * const kRaspberryCollection = @"raspberries";

@interface RaspberryManager()

- (void)removeRaspberryWithName:(NSString*)aName;
- (Raspberry*)getRaspberryWithName:(NSString*)aName;
- (int)getIndexOfRaspberryWithName:(NSString*)aName;
@end



@implementation RaspberryManager
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(RaspberryManager, sharedManager);

- (id)init {
    self = [super init];
    
    if(self) {
        _raspberries = [[NSMutableArray alloc] init];
    }

    return self;
}

//load bookmarks from shared preferences
- (void)loadRaspberries {
    @synchronized(_raspberries) {
        [_raspberries removeAllObjects];
        id data = [[NSUserDefaults standardUserDefaults] arrayForKey:kRaspberryCollection];
        if(data) {
            NSArray *saved = (NSArray *)[NSKeyedUnarchiver unarchiveObjectWithData:data];
            [_raspberries addObjectsFromArray:saved];
        }
    }
}

//save bookmarks to shared preferences
- (void)saveRaspberries {
    @synchronized(_raspberries) {
        NSMutableArray *rpis = [self getRaspberries];
        if(rpis && [rpis count]) {
            NSData *data = [NSKeyedArchiver archivedDataWithRootObject:rpis];
            if (data){
                [[NSUserDefaults standardUserDefaults] setObject:data forKey:kRaspberryCollection];
                [[NSUserDefaults standardUserDefaults] synchronize];
            }
        }
    }
}

- (void)clearRaspberries {
    @synchronized(_raspberries) {
        [_raspberries removeAllObjects];
    }
}

- (Raspberry*)addRaspberry:(Raspberry*)aRaspberry {
    Raspberry *existing = [self getRaspberryWithName:aRaspberry.slaveNodeName];
    
    if(existing) {
        return existing;
    }
    
    @synchronized(_raspberries) {
        [_raspberries addObject:aRaspberry];
    }
    
    return aRaspberry;
}

- (NSMutableArray*)getRaspberries {
    NSMutableArray *bookmarks;
    @synchronized(_raspberries) {
        bookmarks = [NSMutableArray arrayWithArray:_raspberries];
    }
    return bookmarks;
}

- (void)removeRaspberryWithName:(NSString*)aName {
    Raspberry *bookmark = [self getRaspberryWithName:aName];
    if(bookmark) {
        @synchronized(_raspberries) {
            [_raspberries removeObject:bookmark];
        }
    }
}

- (Raspberry*)getRaspberryWithName:(NSString*)aName {
    @synchronized(_raspberries) {
        for(Raspberry *rpi in _raspberries) {
            if([rpi.slaveNodeName isEqualToString:aName]) {
                return rpi;
            }
        }
    }
    
    return nil;
}

- (int)getIndexOfRaspberryWithName:(NSString*)aName {
    for(int i=0; i<_raspberries.count; ++i) {
        Raspberry *rpi = [_raspberries objectAtIndex:i];
        if([rpi.slaveNodeName isEqualToString:aName]) {
            return i;
        }
    }

    return -1;
}

@end
