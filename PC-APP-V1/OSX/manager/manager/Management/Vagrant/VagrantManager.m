//
//  VagrantInstanceCollection.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "VagrantGlobalStatusScanner.h"
#import "SynthesizeSingleton.h"
#import "VagrantManager.h"

#import "NullStringChecker.h"
#import "TaskOutputWindow.h"
#import "Util.h"

@interface VagrantManager()
@property (strong, nonatomic) NSTimer *refreshTimer;
- (void)instanceAdded:(VagrantInstance*)instance;
- (void)instanceRemoved:(VagrantInstance*)instance;
- (void)instanceUpdated:(VagrantInstance*)oldInstance withInstance:(VagrantInstance*)newInstance;
@end

@implementation VagrantManager {
    //all known vagrant instances
    NSMutableArray *_instances;
    
    //map provider identifiers to providers
    NSMutableDictionary *_providers;
    
    BOOL isRefreshingVagrantMachines;
    int queuedRefreshes;
}
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(VagrantManager, sharedManager);

- (id)init {
    self = [super init];
    
    if(self) {
        _instances = [[NSMutableArray alloc] init];
        _providers = [[NSMutableDictionary alloc] init];
        isRefreshingVagrantMachines = NO;
        queuedRefreshes = 0;
        
    }
    
    return self;
}

//get all instances
- (NSArray*)getInstances {
    return [NSArray arrayWithArray:_instances];
}

//get count of machines in running state
- (int)getRunningVmCount {
    int count = 0;
    
    @synchronized(_instances) {
        for(VagrantInstance *instance in _instances) {
            for(VagrantMachine *machine in instance.machines) {
                if(machine.state == RunningState) {
                    ++count;
                }
            }
        }
    }
    
    return count;
}

//get count of machines in a particular state
- (NSArray*)getMachinesWithState:(VagrantMachineState)state {
    NSMutableArray *machines = [[NSMutableArray alloc] init];
    for(VagrantInstance *instance in _instances) {
        for(VagrantMachine *machine in instance.machines) {
            if(machine.state == state) {
                [machines addObject:machine];
            }
        }
    }
    
    return machines;
}

//register a new service provider
- (void)registerServiceProvider:(id<VirtualMachineServiceProvider>)provider {
    [_providers setObject:provider forKey:[provider getProviderIdentifier]];
}

//get instance at a particular path
- (VagrantInstance*)getInstanceForPath:(NSString*)path {
    path = [Util trimTrailingSlash:path];
    
    for(VagrantInstance *instance in _instances) {
        if([instance.path isEqualToString:path]) {
            return instance;
        }
    }
    
    return nil;
}

//refresh list of instances by querying bookmarks, service providers, and NFS
- (void)refreshInstances {
    
    NSMutableArray<VagrantInstance*> *instances = [[NSMutableArray alloc] init];
    NSMutableArray *allPaths = [[NSMutableArray alloc] init];

    //TODO: display name should be editable. beside, Vagrant instance must have an unique id
    
    //scan vagrant global-status output
    VagrantGlobalStatusScanner *globalStatusScanner = [[VagrantGlobalStatusScanner alloc] init];
    for(NSString *path in [globalStatusScanner getInstancePaths]) {
        if([path isEqualToString:@"/pocket/boxes"] && ![allPaths containsObject:path]) {
            [allPaths addObject:path];
            [instances addObject:[[VagrantInstance alloc] initWithPath:path displayName:@"Cluster 1" providerIdentifier:nil]];
        }
    }

    //create instance for each detected path
    NSDictionary *detectedPaths = [self detectInstancePaths];
    for(NSString *providerIdentifier in [detectedPaths allKeys]) {
        NSArray *paths = [detectedPaths objectForKey:providerIdentifier];
        for(NSString *path in paths) {
            //make sure it has not already been detected
            if([path isEqualToString:@"/pocket/boxes"] && ![allPaths containsObject:path]) {
                [allPaths addObject:path];
                [instances addObject:[[VagrantInstance alloc] initWithPath:path displayName:@"Cluster 1" providerIdentifier:providerIdentifier]];
            }
        }
    }

    NSMutableArray *validPaths = [[NSMutableArray alloc] init];    
    //query all known instances for machines, process in parallel
    dispatch_group_t queryMachinesGroup = dispatch_group_create();
    dispatch_queue_t queryMachinesQueue = dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0);
    for(VagrantInstance *instance in instances) {
        dispatch_group_async(queryMachinesGroup, queryMachinesQueue, ^{
            //query instance machines
            [instance queryMachines];
            
            @synchronized(_instances) {
                VagrantInstance *existingInstance = [self getInstanceForPath:instance.path];
                if(existingInstance) {
                    
                    //instance already exists, check for changes
                    int idx = (int)[_instances indexOfObject:existingInstance];
                    if(instance.machines.count != existingInstance.machines.count || ![existingInstance.displayName isEqualToString:instance.displayName] || ![existingInstance.providerIdentifier isEqualToString:instance.providerIdentifier]) {
                        //instance has updated
                        [_instances replaceObjectAtIndex:idx withObject:instance];
                        [self instanceUpdated:existingInstance withInstance:instance];
                        
                    } else {
                        
                        for(VagrantMachine *machine in instance.machines) {
                            VagrantMachine *existingMachine = [existingInstance getMachineWithName:machine.name];
                            
                            if(!existingMachine || ![existingMachine.stateString isEqualToString:machine.stateString]) {
                                //machine did not exist, or state has changed
                                [_instances replaceObjectAtIndex:idx withObject:instance];
                                [self instanceUpdated:existingInstance withInstance:instance];
                            }
                        }
                    }
                } else {
                    //new instance
                    [_instances addObject:instance];
                    [self instanceAdded:instance];
                }
                
                //add path to list for pruning stale instances
                [validPaths addObject:instance.path];
            }
        });
    }
    
    //wait for the machine queries to finish
    dispatch_group_wait(queryMachinesGroup, DISPATCH_TIME_FOREVER);
    
    for(int i=(int)_instances.count-1; i>=0; --i) {
        VagrantInstance *instance = [_instances objectAtIndex:i];
        if(![validPaths containsObject:instance.path]) {
            [_instances removeObjectAtIndex:i];
            [self instanceRemoved:instance];
            
            //TODO: "last seen" functionality may have to be implemented here as well so that this instance doesn't disappear from the list during this pass
        }
    }
}

//query all service providers for instances
- (NSDictionary*)detectInstancePaths {
    NSMutableArray *allPaths = [[NSMutableArray alloc] init];
    NSMutableDictionary *keyedPaths = [[NSMutableDictionary alloc] init];
    
    //find Vagrant instances for each registered provider
    for(id<VirtualMachineServiceProvider> provider in [_providers allValues]) {
        NSArray *paths = [provider getVagrantInstancePaths];
        NSMutableArray *uniquePaths = [[NSMutableArray alloc] init];
        //make sure we haven't already detected this path
        for(NSString *path in paths) {
            NSString *p = [Util trimTrailingSlash:path];
            if(![allPaths containsObject:p]) {
                [allPaths addObject:p];
                [uniquePaths addObject:p];
            }
        }
        [keyedPaths setObject:uniquePaths forKey:[provider getProviderIdentifier]];
    }
    
    return [NSDictionary dictionaryWithDictionary:keyedPaths];
}

//try to determine the vagrant provider for an instance
- (NSString*)detectVagrantProvider:(NSString*)path {
    NSFileManager *fileManager = [NSFileManager defaultManager];
    NSError *error = nil;
    NSArray *machinePaths = [[NSFileManager defaultManager] contentsOfDirectoryAtPath:[NSString stringWithFormat:@"%@/.vagrant/machines", path] error:&error];
    
    if(!error && machinePaths) {
        for(NSString *machinePath in machinePaths) {
            for(NSString *providerIdentifier in [self getProviderIdentifiers]) {
                if([fileManager fileExistsAtPath:[NSString stringWithFormat:@"%@/.vagrant/machines/%@/%@", path, machinePath, providerIdentifier]]) {
                    return providerIdentifier;
                }
            }
        }
    }
    
    return @"virtualbox";
}

- (NSArray*)getProviderIdentifiers {
    NSMutableArray *providerIdentifiers = [NSMutableArray arrayWithArray:[_providers allKeys]];
    [providerIdentifiers addObject:@"vmware_workstation"];
    [providerIdentifiers addObject:@"vmware_fusion"];
    [providerIdentifiers addObject:@"docker"];
    return providerIdentifiers;
}

- (NSString *)vboxInterface {
    return (NSString *)[[NSUserDefaults standardUserDefaults] objectForKey:kPCVagrantNetInterface];
}

- (void)setVboxInterface:(NSString *)aVboxIface {
    if(!ISNULL_STRING(aVboxIface)){
        [[NSUserDefaults standardUserDefaults] setObject:aVboxIface forKey:kPCVagrantNetInterface];
        [[NSUserDefaults standardUserDefaults] synchronize];
    }
}

- (void)refreshInstanceRelatedPackages {
    [self.instances makeObjectsPerformSelector:@selector(checkRelatedPackage)];
}


#pragma mark - VAGRANT MACHINE CONTROL
- (void)haltRefreshTimer {
    if (self.refreshTimer) {
        [self.refreshTimer invalidate];
        self.refreshTimer = nil;
    }
}

- (void)refreshTimerState {
    
    [self haltRefreshTimer];
    
    //if ([[NSUserDefaults standardUserDefaults] boolForKey:@"refreshEvery"])
    {
        self.refreshTimer =
        [NSTimer
         scheduledTimerWithTimeInterval:30//[[NSUserDefaults standardUserDefaults] integerForKey:@"refreshEveryInterval"]
         target:self
         selector:@selector(refreshVagrantMachines)
         userInfo:nil
         repeats:YES];
    }
}

- (void)updateRunningVmCount {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kVAGRANT_MANAGER_UPDATE_RUNNING_VM_COUNT
     object:nil
     userInfo:@{kPOCKET_CLUSTER_LIVE_NODE_COUNT: [NSNumber numberWithUnsignedInteger:[self getRunningVmCount]]}];
}

- (void)updateInstancesCount {
    [[NSNotificationCenter defaultCenter]
     postNotificationName:kVAGRANT_MANAGER_UPDATE_INSTANCES_COUNT
     object:nil
     userInfo:@{kPOCKET_CLUSTER_NODE_COUNT: [NSNumber numberWithUnsignedInteger:[[self getInstances] count]]}];
}

- (void)refreshVagrantMachines {
    //only run if not already refreshing
    if(!isRefreshingVagrantMachines) {
        isRefreshingVagrantMachines = YES;
        
        WEAK_SELF(self);
        
        //tell popup controller refreshing has started
        [[NSNotificationCenter defaultCenter] postNotificationName:kVAGRANT_MANAGER_REFRESHING_STARTED object:nil];
        dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
            //tell manager to refresh all instances
            [belf refreshInstances];
            
            //TODO: refactor!
            [belf refreshInstanceRelatedPackages];
            
            dispatch_async(dispatch_get_main_queue(), ^{
                //tell popup controller refreshing has ended
                isRefreshingVagrantMachines = NO;
                [[NSNotificationCenter defaultCenter] postNotificationName:kVAGRANT_MANAGER_REFRESHING_ENDED object:nil];
                [belf updateInstancesCount];
                [belf updateRunningVmCount];
                
                if(queuedRefreshes > 0) {
                    --queuedRefreshes;
                    [belf refreshVagrantMachines];
                }
            });
        });
    } else {
        ++queuedRefreshes;
    }
}


#pragma mark - VAGRANT INSTANCE NOTIFIACTION
- (void)instanceAdded:(VagrantInstance *)instance {
    dispatch_async(dispatch_get_main_queue(), ^{
        [[NSNotificationCenter defaultCenter]
         postNotificationName:kVAGRANT_MANAGER_INSTANCE_ADDED
         object:nil
         userInfo:@{kVAGRANT_MANAGER_INSTANCE: instance}];
    });
}

- (void)instanceRemoved:(VagrantInstance *)instance {
    dispatch_async(dispatch_get_main_queue(), ^{
        [[NSNotificationCenter defaultCenter]
         postNotificationName:kVAGRANT_MANAGER_INSTANCE_REMOVED
         object:nil
         userInfo:@{kVAGRANT_MANAGER_INSTANCE: instance}];
    });
}

- (void)instanceUpdated:(VagrantInstance *)oldInstance withInstance:(VagrantInstance *)newInstance {
    dispatch_async(dispatch_get_main_queue(), ^{
        [[NSNotificationCenter defaultCenter]
         postNotificationName:kVAGRANT_MANAGER_INSTANCE_UPDATED
         object:nil
         userInfo:@{kVAGRANT_MANAGER_INSTANCE_OLD:oldInstance,
                    kVAGRANT_MANAGER_INSTANCE_NEW:newInstance}];
    });
}

#pragma mark - VAGRANT ACTION METHODS
- (void)runVagrantCustomCommand:(NSString*)command withMachine:(VagrantMachine*)machine {
    
    Assert([NSThread isMainThread], @"runVagrantCustomCommand:withMachine: should run in Main Thread");
    
    NSTask *task = [[NSTask alloc] init];
    [task setLaunchPath:@"/bin/bash"];
    
    NSString *taskCommand = [NSString stringWithFormat:@"cd %@; vagrant ssh %@ -c %@", [Util escapeShellArg:machine.instance.path], [Util escapeShellArg:machine.name], [Util escapeShellArg:command]];
    
    [task setArguments:@[@"-l", @"-c", taskCommand]];
    
    TaskOutputWindow *outputWindow = [[TaskOutputWindow alloc] initWithWindowNibName:@"TaskOutputWindow"];
    outputWindow.task = task;
    outputWindow.taskCommand = taskCommand;
    outputWindow.target = machine;
    outputWindow.taskAction = command;
    
    [NSApp activateIgnoringOtherApps:YES];
    [outputWindow showWindow:[Util getApp]];
    
    [[Util getApp] addOpenWindow:outputWindow];
}

- (void)runVagrantAction:(NSString*)action withMachine:(VagrantMachine*)machine {
    
    Assert([NSThread isMainThread], @"runVagrantAction:withMachine: should run in Main Thread");
    
    NSMutableArray *commandParts = [[NSMutableArray alloc] init];
    
    if([action isEqualToString:@"up"]) {
        [commandParts addObject:@"vagrant up"];
        if(machine.instance.providerIdentifier) {
            [commandParts addObject:[NSString stringWithFormat:@"--provider=%@", machine.instance.providerIdentifier]];
        }
    } else if([action isEqualToString:@"up-provision"]) {
        [commandParts addObject:@"vagrant up --provision"];
        if(machine.instance.providerIdentifier) {
            [commandParts addObject:[NSString stringWithFormat:@"--provider=%@", machine.instance.providerIdentifier]];
        }
    } else if([action isEqualToString:@"reload"]) {
        [commandParts addObject:@"vagrant reload"];
    } else if([action isEqualToString:@"suspend"]) {
        [commandParts addObject:@"vagrant suspend"];
    } else if([action isEqualToString:@"halt"]) {
        [commandParts addObject:@"vagrant halt"];
    } else if([action isEqualToString:@"provision"]) {
        [commandParts addObject:@"vagrant provision"];
    } else if([action isEqualToString:@"destroy"]) {
        [commandParts addObject:@"vagrant destroy -f"];
    } else if([action isEqualToString:@"rdp"]) {
        [commandParts addObject:@"vagrant rdp"];
    } else {
        return;
    }
    
    [commandParts addObject:@"--no-color"];
    
    NSString *command = [commandParts componentsJoinedByString:@" "];
    
    NSTask *task = [[NSTask alloc] init];
    [task setLaunchPath:@"/bin/bash"];
    
    NSString *taskCommand = [NSString stringWithFormat:@"cd %@; %@ %@", [Util escapeShellArg:machine.instance.path], command, [Util escapeShellArg:machine.name]];
    
    [task setArguments:@[@"-l", @"-c", taskCommand]];
    
    TaskOutputWindow *outputWindow = [[TaskOutputWindow alloc] initWithWindowNibName:@"TaskOutputWindow"];
    outputWindow.task = task;
    outputWindow.taskCommand = taskCommand;
    outputWindow.target = machine;
    outputWindow.taskAction = command;
    
    [NSApp activateIgnoringOtherApps:YES];
    [outputWindow showWindow:[Util getApp]];
    
    [[Util getApp] addOpenWindow:outputWindow];
}

- (void)runVagrantAction:(NSString*)action withInstance:(VagrantInstance*)instance {
    
    Assert([NSThread isMainThread], @"runVagrantAction:withInstance: should run in Main Thread");
    
    NSMutableArray *commandParts = [[NSMutableArray alloc] init];
    
    if([action isEqualToString:@"up"]) {
        [commandParts addObject:@"vagrant up"];
        if(instance.providerIdentifier) {
            [commandParts addObject:[NSString stringWithFormat:@"--provider=%@", instance.providerIdentifier]];
        }
    } else if([action isEqualToString:@"up-provision"]) {
        [commandParts addObject:@"vagrant up --provision"];
        if(instance.providerIdentifier) {
            [commandParts addObject:[NSString stringWithFormat:@"--provider=%@", instance.providerIdentifier]];
        }
    } else if([action isEqualToString:@"reload"]) {
        [commandParts addObject:@"vagrant reload"];
    } else if([action isEqualToString:@"suspend"]) {
        [commandParts addObject:@"vagrant suspend"];
    } else if([action isEqualToString:@"halt"]) {
        [commandParts addObject:@"vagrant halt"];
    } else if([action isEqualToString:@"provision"]) {
        [commandParts addObject:@"vagrant provision"];
    } else if([action isEqualToString:@"destroy"]) {
        [commandParts addObject:@"vagrant destroy -f"];
    } else if([action isEqualToString:@"rdp"]) {
        [commandParts addObject:@"vagrant rdp"];
    } else {
        return;
    }
    
    [commandParts addObject:@"--no-color"];
    
    NSString *command = [commandParts componentsJoinedByString:@" "];
    
    NSTask *task = [[NSTask alloc] init];
    [task setLaunchPath:@"/bin/bash"];
    
    NSString *taskCommand = [NSString stringWithFormat:@"cd %@; %@", [Util escapeShellArg:instance.path], command];
    
    [task setArguments:@[@"-c", @"-l", taskCommand]];
    
    TaskOutputWindow *outputWindow = [[TaskOutputWindow alloc] initWithWindowNibName:@"TaskOutputWindow"];
    outputWindow.task = task;
    outputWindow.taskCommand = taskCommand;
    outputWindow.target = instance;
    outputWindow.taskAction = command;
    
    [NSApp activateIgnoringOtherApps:YES];
    [outputWindow showWindow:[Util getApp]];
    
    [[Util getApp] addOpenWindow:outputWindow];
}

- (void)openInstanceInFinder:(VagrantInstance *)instance {
    NSString *path = instance.path;
    
    BOOL isDir = NO;
    if([[NSFileManager defaultManager] fileExistsAtPath:path isDirectory:&isDir] && isDir) {
        NSURL *fileURL = [NSURL fileURLWithPath:path];
        [[NSWorkspace sharedWorkspace] openURL:fileURL];
    } else {
        [[NSAlert alertWithMessageText:[NSString stringWithFormat:@"Path not found: %@", path] defaultButton:@"OK" alternateButton:nil otherButton:nil informativeTextWithFormat:@""] runModal];
    }
}

- (void)openInstanceInTerminal:(VagrantInstance *)instance {
    NSString *path = instance.path;
    
    BOOL isDir = NO;
    if([[NSFileManager defaultManager] fileExistsAtPath:path isDirectory:&isDir] && isDir) {
        [Util runTerminalCommand:[NSString stringWithFormat:@"cd %@", [Util escapeShellArg:path]]];
    } else {
        [[NSAlert alertWithMessageText:[NSString stringWithFormat:@"Path not found: %@", path] defaultButton:@"OK" alternateButton:nil otherButton:nil informativeTextWithFormat:@""] runModal];
    }
}

#pragma mark - NativeMenu Delegate
- (void)performVagrantAction:(NSString *)action withInstance:(VagrantInstance *)instance {
    if([action isEqualToString:@"ssh"]) {
        NSString *action = [NSString stringWithFormat:@"cd %@; vagrant ssh", [Util escapeShellArg:instance.path]];
        [Util runTerminalCommand:action];
    } else {
        [self runVagrantAction:action withInstance:instance];
    }
}

- (void)performVagrantAction:(NSString *)action withMachine:(VagrantMachine *)machine {
    if([action isEqualToString:@"ssh"]) {
        NSString *action = [NSString stringWithFormat:@"cd %@; vagrant ssh %@", [Util escapeShellArg:machine.instance.path], machine.name];
        [Util runTerminalCommand:action];
    } else {
        [self runVagrantAction:action withMachine:machine];
    }
}
@end
