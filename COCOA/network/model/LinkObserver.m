/*
 Copyright (c) 2015 funkensturm. https://github.com/halo/LinkLiar
 
 Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the
 "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish,
 distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to
 the following conditions:
 
 The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 
 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
 LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

#import "LinkObserver.h"

static void observerNotificationProxy(SCDynamicStoreRef store, CFArrayRef triggeredKeys, void* info) {
  NSString *key;
  NSEnumerator *keys = [(__bridge NSArray *)triggeredKeys objectEnumerator];
  NSNotificationCenter* notificationCenter = [NSNotificationCenter defaultCenter];
  
  while (key = [keys nextObject]) {
    NSLog(@"%@", key);
    [notificationCenter postNotificationName:key object:(__bridge id)info userInfo:(__bridge NSDictionary*)SCDynamicStoreCopyValue(store, (__bridge CFStringRef) key)];
  }
}

static void _appendNotificationPattenKey(CFMutableArrayRef keysArray, CFMutableArrayRef patternArray, CFStringRef patternKey) {
    CFStringRef storeKey;
    storeKey = SCDynamicStoreKeyCreateNetworkGlobalEntity(NULL,
                                                          kSCDynamicStoreDomainState,
                                                          patternKey);
    CFArrayAppendValue(keysArray, storeKey);
    CFRelease(storeKey);
    
    storeKey = SCDynamicStoreKeyCreateNetworkServiceEntity(NULL,
                                                           kSCDynamicStoreDomainState,
                                                           kSCCompAnyRegex,
                                                           patternKey);
    CFArrayAppendValue(patternArray, storeKey);
    CFRelease(storeKey);
}


@implementation LinkObserver
@synthesize dynamicStore;
@synthesize runLoop;

- (instancetype) init {
    self = [super init];
    if (self) {
        //CFRunLoopAddSource([[NSRunLoop currentRunLoop] getCFRunLoop], self.runLoop, kCFRunLoopCommonModes);
        //SCDynamicStoreSetNotificationKeys(self.dynamicStore, NULL, (__bridge CFArrayRef)@[@".*"]);
        
        // Although this is not specifically designed to run on main runloop, it is desired to do so.
        // in that case, it is beneficial to be in common mode than default mode.
        // http://stackoverflow.com/questions/7222449/nsdefaultrunloopmode-vs-nsrunloopcommonmodes
        CFRunLoopAddSource(CFRunLoopGetCurrent(), self.runLoop, kCFRunLoopCommonModes);
    }
    return self;
}

- (void)dealloc {
    CFRunLoopRemoveSource(CFRunLoopGetCurrent(), self.runLoop, kCFRunLoopCommonModes);
    if (runLoop != NULL) {
        CFRelease(runLoop);
    }
    if (dynamicStore != NULL) {
        CFRelease(dynamicStore);
    }
}

- (SCDynamicStoreRef) dynamicStore {
    if (dynamicStore) {
        return dynamicStore;
    }
    SCDynamicStoreContext context = {0, (__bridge void *)(self), NULL, NULL, NULL};
    dynamicStore = SCDynamicStoreCreate(NULL,
                                        CFSTR("PocketClusterNetworkChangedCallback"),
                                        observerNotificationProxy,
                                        &context);
    
    CFMutableArrayRef keysArray = CFArrayCreateMutable(NULL, 0, &kCFTypeArrayCallBacks);
    CFMutableArrayRef regexArray = CFArrayCreateMutable(NULL, 0, &kCFTypeArrayCallBacks);
    
    _appendNotificationPattenKey(keysArray, regexArray, kSCEntNetLink);
    _appendNotificationPattenKey(keysArray, regexArray, kSCEntNetDNS);
    // TODO : IPv6 is not supported.
    _appendNotificationPattenKey(keysArray, regexArray, kSCEntNetIPv4);
    _appendNotificationPattenKey(keysArray, regexArray, kSCEntNetAirPort);
    _appendNotificationPattenKey(keysArray, regexArray, kSCEntNetDHCP);
    _appendNotificationPattenKey(keysArray, regexArray, kSCEntNetEthernet);
    _appendNotificationPattenKey(keysArray, regexArray, kSCEntNetFireWire);
    _appendNotificationPattenKey(keysArray, regexArray, kSCEntNetInterface);
    
    if (!SCDynamicStoreSetNotificationKeys(dynamicStore, keysArray, regexArray)) {
        CFRelease(dynamicStore);
        dynamicStore = NULL;
    }
    
    CFRelease(keysArray);
    CFRelease(regexArray);
    return dynamicStore;
}

- (CFRunLoopSourceRef) runLoop {
    if (runLoop) {
        return runLoop;
    }
    runLoop = SCDynamicStoreCreateRunLoopSource(NULL, self.dynamicStore, 0);
    return runLoop;
}

@end
