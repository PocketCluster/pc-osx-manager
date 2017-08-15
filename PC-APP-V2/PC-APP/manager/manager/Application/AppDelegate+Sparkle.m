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
    Log(@"%s",__PRETTY_FUNCTION__);    
}

// Called when a valid update is found by the update driver.
- (void)updater:(SUUpdater *)updater didFindValidUpdate:(SUAppcastItem *)item {
    Log(@"%s %@",__PRETTY_FUNCTION__, item.fileURL.description);
    //TODO : make sure when didFindValudUpdate or updaterDidNotFindUpdate gets called before execute engine
}

// Called when a valid update is not found.
- (void)updaterDidNotFindUpdate:(SUUpdater *)updater {
    Log(@"%s",__PRETTY_FUNCTION__);
    //TODO : make sure when didFindValudUpdate or updaterDidNotFindUpdate gets called before execute engine
}

//  Called immediately before downloading the specified update.
- (void)updater:(SUUpdater *)updater willDownloadUpdate:(SUAppcastItem *)item withRequest:(NSMutableURLRequest *)request {
    Log(@"%s",__PRETTY_FUNCTION__);
}

// Called after the specified update failed to download.
- (void)updater:(SUUpdater *)updater failedToDownloadUpdate:(SUAppcastItem *)item error:(NSError *)error {
    Log(@"%s",__PRETTY_FUNCTION__);
}

//  Called when the user clicks the cancel button while and update is being downloaded.
- (void)userDidCancelDownload:(SUUpdater *)updater {
    Log(@"%s",__PRETTY_FUNCTION__);
}

//  Called immediately before installing the specified update.
- (void)updater:(SUUpdater *)updater willInstallUpdate:(SUAppcastItem *)item {
    Log(@"%s",__PRETTY_FUNCTION__);
}

// Returns whether the relaunch should be delayed in order to perform other tasks.
//
// This is not called if the user didn't relaunch on the previous update,
// in that case it will immediately restart.

// - (BOOL)updater:(SUUpdater *)updater shouldPostponeRelaunchForUpdate:(SUAppcastItem *)item untilInvoking:(NSInvocation *)invocation {
// }

//  Called before an updater shows a modal alert window,
- (void)updaterWillShowModalAlert:(SUUpdater *)updater {
    Log(@"%s",__PRETTY_FUNCTION__);
}

//  Called after an updater shows a modal alert window,
- (void)updaterDidShowModalAlert:(SUUpdater *)updater {
    Log(@"%s",__PRETTY_FUNCTION__);
}

// Called when an update is scheduled to be silently installed on quit.
- (void)updater:(SUUpdater *)updater willInstallUpdateOnQuit:(SUAppcastItem *)item immediateInstallationInvocation:(NSInvocation *)invocation {
    Log(@"%s",__PRETTY_FUNCTION__);    
}


#pragma mark - SPARKLE UPDATER DELEGATE

// from v0.1.3
#if 0
- (NSArray*)feedParametersForUpdater:(SUUpdater *)updater sendingSystemProfile:(BOOL)sendingProfile {
    NSMutableArray *data = [[NSMutableArray alloc] init];
    [data addObject:@{@"key": @"machineid", @"value": [Util getMachineId]}];
    [data addObject:@{@"key": @"appversion", @"value": [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleShortVersionString"]}];
    if(sendingProfile) {
        [data addObject:@{@"key": @"profile", @"value": @"1"}];
    }
    
    return data;
}

- (void)updater:(SUUpdater *)updater didFindValidUpdate:(SUAppcastItem *)update {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kPOCKET_CLUSTER_UPDATE_AVAILABLE
     object:nil
     userInfo:@{kPOCKET_CLUSTER_UPDATE_VALUE: @(YES)}];
}

- (void)updaterDidNotFindUpdate:(SUUpdater *)update {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kPOCKET_CLUSTER_UPDATE_AVAILABLE
     object:nil
     userInfo:@{kPOCKET_CLUSTER_UPDATE_VALUE: @(NO)}];
}

- (id<SUVersionComparison>)versionComparatorForUpdater:(SUUpdater *)updater {
    return [[VersionComparison alloc] init];
}

- (SUAppcastItem *)bestValidUpdateInAppcast:(SUAppcast *)appcast forUpdater:(SUUpdater *)bundle {
    SUAppcastItem *bestItem = nil;
    
    NSString *appVersion = [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleShortVersionString"];
    
    for(SUAppcastItem *item in [appcast items]) {
        if([appVersion compare:item.versionString options:NSNumericSearch] == NSOrderedAscending) {
            if(!bestItem || [bestItem.versionString compare:item.versionString options:NSNumericSearch] == NSOrderedAscending) {
                bestItem = item;
            }
        }
    }
    
    return bestItem;
}
#endif

@end
