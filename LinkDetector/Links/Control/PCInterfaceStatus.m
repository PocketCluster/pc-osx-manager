//
//  PCInterfaceList.m
//  NETUTIL
//
//  Created by Almighty Kim on 10/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCInterfaceStatus.h"
#import "LinkObserver.h"
#import "PCInterfaceTypes.h"

static __weak PCInterfaceStatus* pc_interface_status_ptr;

void interface_status(pc_interface_callback callback) {
    

    
    
    
    
}

@interface PCInterfaceStatus()
@property (readonly) LinkObserver *linkObserver;
-(CFArrayRef)_refreshInterfaceStatus;
@end

@implementation PCInterfaceStatus
@synthesize linkObserver;

-(instancetype)init {
    
    self = [super init];
    if (self) {
        pc_interface_status_ptr = self;
    }
    return self;
}

#pragma mark - PROPERTIES
- (LinkObserver*) linkObserver {
    if (linkObserver) return linkObserver;
    linkObserver = [LinkObserver new];
    return linkObserver;
}

#pragma mark - METHODS
- (void) interfacesDidChange:(NSNotification*)notififcation {
    NSLog(@"Interface change detected...");
}

- (void) startMonitoring {
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(interfacesDidChange:) name:@"State:/Network/Interface" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(interfacesDidChange:) name:@"State:/Network/Interface/en0/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(interfacesDidChange:) name:@"State:/Network/Interface/en1/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(interfacesDidChange:) name:@"State:/Network/Interface/en2/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(interfacesDidChange:) name:@"State:/Network/Interface/en3/AirPort" object:self.linkObserver];
}

- (void) stopMonitoring {
    [[NSNotificationCenter defaultCenter] removeObserver:self name:@"State:/Network/Interface" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] removeObserver:self name:@"State:/Network/Interface/en0/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] removeObserver:self name:@"State:/Network/Interface/en1/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] removeObserver:self name:@"State:/Network/Interface/en2/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] removeObserver:self name:@"State:/Network/Interface/en3/AirPort" object:self.linkObserver];
}

-(CFArrayRef)_refreshInterfaceStatus {
    
    return nil;
}

@end
