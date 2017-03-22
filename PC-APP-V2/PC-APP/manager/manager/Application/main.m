//
//  main.m
//  manager
//
//  Created by Almighty Kim on 10/15/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

int
PCNativeMainStart(int argc, const char * argv[]) {
    return NSApplicationMain(argc, argv);
}

void
PCNativeMainStop() {
    [NSApp terminate:nil];
}
