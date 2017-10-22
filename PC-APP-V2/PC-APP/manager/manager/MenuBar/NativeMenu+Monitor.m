//
//  NativeMenu+Monitor.m
//  manager
//
//  Created by Almighty Kim on 10/22/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "NativeMenuAddition.h"
#import "NativeMenu+Monitor.h"

@implementation NativeMenu(Monitor)

- (void) onUpdatedWith:(StatusCache *)aCache forPackageListAvailable:(BOOL)isSuccess {
}

- (void) onUpdatedWith:(StatusCache *)aCache forPackageListInstalled:(BOOL)isSuccess {
    if (!isSuccess) {
        return;
    }

    NSInteger indexBegin = ([self.statusItem.menu
                             indexOfItem:[self.statusItem.menu
                                          itemWithTag:MENUITEM_PKG_DIV]] + 1);

    // remove all old package menues
    for (NSMenuItem *item in [self.statusItem.menu itemArray]) {
        if ([item tag] < PKG_TAG_BUMPER) {
            continue;
        }
        [self.statusItem.menu removeItem:item];
    }

    // all the package list
    NSArray<Package *>* plst = [aCache packageList];
    NSInteger pndx = 0;

    // add packages according to the list
    for (Package *pkg in plst) {
        if (![pkg installed]) {
            continue;
        }
        
        NSMenuItem *penu = [[NSMenuItem alloc] initWithTitle:pkg.packageDescription action:nil keyEquivalent:@""];
        [penu setTag:PKG_TAG_BUILDER(pndx)];
        [penu setSubmenu:[NSMenu new]];

        // add submenu - start
        NSMenuItem *smStart = [[NSMenuItem alloc] initWithTitle:@"Start" action:@selector(startPackage) keyEquivalent:@""];
        [smStart setTarget:self];
        [penu.submenu addItem:smStart];

        // add submneu - stop
        NSMenuItem *smStop = [[NSMenuItem alloc] initWithTitle:@"Stop" action:@selector(stopPackage) keyEquivalent:@""];
        [smStop setTarget:self];
        [penu.submenu addItem:smStop];

        // add open web port menu
        NSMenuItem *smWeb = [[NSMenuItem alloc] initWithTitle:@"Web Console" action:@selector(openWebConsole) keyEquivalent:@""];
        [smWeb setTarget:self];
        [penu.submenu addItem:smWeb];

        [self.statusItem.menu insertItem:penu atIndex:(indexBegin + pndx)];
    }
}

- (void) startPackage {

}

- (void) stopPackage {

}

- (void) openWebConsole {

}
@end
