//
//  AppDelegate+Sparkle.m
//  manager
//
//  Created by Almighty Kim on 4/3/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+Sparkle.h"

@implementation AppDelegate(Sparkle)

// Called after Sparkle has downloaded the appcast from the remote server.
- (void)updater:(SUUpdater *)updater didFinishLoadingAppcast:(SUAppcast *)appcast {
    NSLog(@"%s",__PRETTY_FUNCTION__);
}

// Called when a valid update is found by the update driver.
- (void)updater:(SUUpdater *)updater didFindValidUpdate:(SUAppcastItem *)item {
    NSLog(@"%s %@",__PRETTY_FUNCTION__, item.fileURL.description);

    [[NSOperationQueue mainQueue] addOperationWithBlock:^{
        [[SUUpdater sharedUpdater] checkForUpdates:nil];
    }];
}

// Called when a valid update is not found.
- (void)updaterDidNotFindUpdate:(SUUpdater *)updater {
    NSLog(@"%s",__PRETTY_FUNCTION__);
}

//  Called immediately before downloading the specified update.
- (void)updater:(SUUpdater *)updater willDownloadUpdate:(SUAppcastItem *)item withRequest:(NSMutableURLRequest *)request {
    NSLog(@"%s",__PRETTY_FUNCTION__);
}

// Called after the specified update failed to download.
- (void)updater:(SUUpdater *)updater failedToDownloadUpdate:(SUAppcastItem *)item error:(NSError *)error {
    NSLog(@"%s",__PRETTY_FUNCTION__);
}

//  Called when the user clicks the cancel button while and update is being downloaded.
- (void)userDidCancelDownload:(SUUpdater *)updater {
    NSLog(@"%s",__PRETTY_FUNCTION__);
}

//  Called immediately before installing the specified update.
- (void)updater:(SUUpdater *)updater willInstallUpdate:(SUAppcastItem *)item {
    NSLog(@"%s",__PRETTY_FUNCTION__);
}

// Returns whether the relaunch should be delayed in order to perform other tasks.
//
// This is not called if the user didn't relaunch on the previous update,
// in that case it will immediately restart.

// - (BOOL)updater:(SUUpdater *)updater shouldPostponeRelaunchForUpdate:(SUAppcastItem *)item untilInvoking:(NSInvocation *)invocation {
// }

//  Called before an updater shows a modal alert window,
- (void)updaterWillShowModalAlert:(SUUpdater *)updater {
    NSLog(@"%s",__PRETTY_FUNCTION__);
}

//  Called after an updater shows a modal alert window,
- (void)updaterDidShowModalAlert:(SUUpdater *)updater {
    NSLog(@"%s",__PRETTY_FUNCTION__);
}

// Called when an update is scheduled to be silently installed on quit.
- (void)updater:(SUUpdater *)updater willInstallUpdateOnQuit:(SUAppcastItem *)item immediateInstallationInvocation:(NSInvocation *)invocation {
    NSLog(@"%s",__PRETTY_FUNCTION__);    
}
@end
