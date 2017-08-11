//
//  PCPrefWC.m
//  manager
//
//  Created by Almighty Kim on 10/30/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCPrefWC.h"
#import "PCConstants.h"

@interface PCPrefWC()
- (void)setLaunchOnLogin:(BOOL)launchOnLogin;
- (BOOL)willStartAtLogin;
@end

@implementation PCPrefWC

+ (NSDictionary *)defaultPreferences {
    return @{kPCPrefDefaultTerm:@101};
}

+ (void)load {
    [[NSUserDefaults standardUserDefaults] registerDefaults:[self defaultPreferences]];
}

- (void)windowDidLoad {
    [super windowDidLoad];

    [self.launchAtLoginCheckBox setState:[self willStartAtLogin] ? NSOnState : NSOffState];

    NSNumber *terminalPreference = (NSNumber *)[[NSUserDefaults standardUserDefaults] stringForKey:kPCPrefDefaultTerm];
    if ([terminalPreference integerValue] == 101) {
        [self.terminalPreferencePopUpButton selectItemWithTag:101];
    } else {
        [self.terminalPreferencePopUpButton selectItemWithTag:100];
    }
}

- (IBAction)terminalPreferencePopUpButtonClicked:(id)sender {

    NSNumber *term = [NSNumber numberWithInteger:self.terminalPreferencePopUpButton.selectedItem.tag];
    [[NSUserDefaults standardUserDefaults] setValue:term forKey:kPCPrefDefaultTerm];
    [[NSUserDefaults standardUserDefaults] synchronize];
}

- (IBAction)launchAtLoginCheckBoxClicked:(id)sender {
    [self setLaunchOnLogin:(self.launchAtLoginCheckBox.state == NSOnState)];
}

#pragma mark - LAUNCH ON LOGIN CHECK
- (void)setLaunchOnLogin:(BOOL)launchOnLogin {
    NSURL *bundleURL = [NSURL fileURLWithPath:[[NSBundle mainBundle] bundlePath]];
    
    LSSharedFileListItemRef existingItem = NULL;
    
    LSSharedFileListRef loginItems = LSSharedFileListCreate(NULL, kLSSharedFileListSessionLoginItems, NULL);
    if (loginItems) {
        UInt32 seed = 0U;
        NSArray *currentLoginItems = (__bridge NSArray *)(LSSharedFileListCopySnapshot(loginItems, &seed));
        for (id itemObject in currentLoginItems) {
            LSSharedFileListItemRef item = (__bridge LSSharedFileListItemRef)itemObject;
            
            UInt32 resolutionFlags = kLSSharedFileListNoUserInteraction | kLSSharedFileListDoNotMountVolumes;
            CFURLRef URL = NULL;
            OSStatus err = LSSharedFileListItemResolve(item, resolutionFlags, &URL, NULL);
            if (err == noErr) {
                Boolean foundIt = CFEqual(URL, (__bridge CFTypeRef)(bundleURL));
                CFRelease(URL);
                
                if (foundIt) {
                    existingItem = item;
                    break;
                }
            }
        }
        
        if (launchOnLogin && (existingItem == NULL)) {
            LSSharedFileListInsertItemURL(loginItems, kLSSharedFileListItemBeforeFirst, NULL, NULL, (__bridge CFURLRef)bundleURL, NULL, NULL);
            
        } else if (!launchOnLogin && (existingItem != NULL)) {
            LSSharedFileListItemRemove(loginItems, existingItem);
        }
        
        CFRelease(loginItems);
    }
}

- (BOOL)willStartAtLogin {
    NSURL *bundleURL = [NSURL fileURLWithPath:[[NSBundle mainBundle] bundlePath]];
    BOOL foundIt = NO;
    
    LSSharedFileListRef loginItems = LSSharedFileListCreate(NULL, kLSSharedFileListSessionLoginItems, NULL);
    if (loginItems) {
        UInt32 seed = 0U;
        NSArray *currentLoginItems = (__bridge NSArray *)(LSSharedFileListCopySnapshot(loginItems, &seed));
        for (id itemObject in currentLoginItems) {
            LSSharedFileListItemRef item = (__bridge LSSharedFileListItemRef)itemObject;
            
            UInt32 resolutionFlags = kLSSharedFileListNoUserInteraction | kLSSharedFileListDoNotMountVolumes;
            CFURLRef URL = NULL;
            OSStatus err = LSSharedFileListItemResolve(item, resolutionFlags, &URL, NULL);
            if (err == noErr) {
                foundIt = (BOOL)CFEqual(URL, (__bridge CFTypeRef)(bundleURL));
                CFRelease(URL);
                
                if (foundIt)
                    break;
            }
        }
        CFRelease(loginItems);
    }
    return foundIt;
}

@end
