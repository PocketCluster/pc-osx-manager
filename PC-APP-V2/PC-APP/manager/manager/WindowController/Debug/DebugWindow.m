//
//  DebugWindow.m
//  manager
//
//  Created by Almighty Kim on 4/10/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "DebugWindow.h"
#import "pc-core.h"

@interface DebugWindow ()
@end

@implementation DebugWindow

- (void)windowDidLoad {
    [super windowDidLoad];
    
    // Implement this method to handle any initialization after your window controller's window has been loaded from its nib file.
}


- (IBAction)opsCmdTeleportStart:(id)sender {
    OpsCmdTeleportStart();
}

- (IBAction)opsCmdTeleportStop:(id)sender {
    OpsCmdTeleportStop();
}

- (IBAction)opsCmdRegistryStart:(id)sender {
    OpsCmdRegistryStart();
}

- (IBAction)opsCmdRegistryStop:(id)sender {
    OpsCmdRegistryStop();
}

- (IBAction)opsCmdCntrOrchStart:(id)sender {
    OpsCmdCntrOrchStart();
}

- (IBAction)opsCmdCntrOrchStop:(id)sender {
    OpsCmdCntrOrchStop();
}

- (IBAction)opsCmdStorageStart:(id)sender {
    OpsCmdStorageStart();
}

- (IBAction)opsCmdStorageStop:(id)sender {
    OpsCmdStorageStop();
}

@end
