//
//  RaspberryManager.m
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//
#import "VagrantManager.h"
#import "SynthesizeSingleton.h"
#import "RaspberryManager.h"
#import "PCConstants.h"


@interface RaspberryManager()
@property (nonatomic, strong) NSMutableArray *raspberries;
@property (nonatomic, strong) GCDAsyncUdpSocket *multSocket;
@property (nonatomic, strong) NSMutableArray<GCDAsyncUdpSocketDelegate> *multSockDelegates;

- (void)removeRaspberryWithName:(NSString*)aName;
- (Raspberry*)getRaspberryWithName:(NSString*)aName;
- (int)getIndexOfRaspberryWithName:(NSString*)aName;
@end

@implementation RaspberryManager
SYNTHESIZE_SINGLETON_FOR_CLASS_WITH_ACCESSOR(RaspberryManager, sharedManager);

- (id)init {
    self = [super init];
    
    if(self) {
        self.raspberries = [[NSMutableArray alloc] init];
        self.multSocket = [[GCDAsyncUdpSocket alloc] initWithDelegate:self delegateQueue:dispatch_get_main_queue()];
        [self.multSocket setIPv6Enabled:NO];
        self.multSockDelegates = [NSMutableArray<GCDAsyncUdpSocketDelegate> arrayWithCapacity:0];

    }

    return self;
}

//load bookmarks from shared preferences
- (void)loadRaspberries {
    @synchronized(_raspberries) {
        [_raspberries removeAllObjects];
        id data = [[NSUserDefaults standardUserDefaults] dataForKey:kRaspberryCollection];
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

#warning FIX to Actually live ones
- (NSUInteger)liveRaspberryCount {
    return [self.raspberries count];
}

- (NSUInteger)raspberryCount {
    return [self.raspberries count];
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



#pragma mark - GCDAsyncUdpSocket MANAGEMENT

- (void)addMultDelegateToQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates addObject:aDelegate];
    }
}

- (void)removeMultDelegateFromQueue:(id<GCDAsyncUdpSocketDelegate>)aDelegate {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates removeObject:aDelegate];
    }
}

-(void)startMulticastSocket
{
    // START udp echo server
    NSError *error = nil;
    if (![self.multSocket bindToPort:PAGENT_SEND_PORT error:&error])
    {
        Log(@"Error starting server (bind): %@", error);
        return;
    }
    
    [self.multSocket joinMulticastGroup:POCKETCAST_GROUP error:&error];
    
    if (![self.multSocket beginReceiving:&error])
    {
        [self.multSocket close];
        return;
    }
}

- (void)stopMulticastSocket
{
    [self.multSocket closeAfterSending];
}

- (void)multicastData:(NSData *)aData
{
    [self.multSocket
     sendData:aData
     toHost:POCKETCAST_GROUP
     port:PAGENT_RECV_PORT
     withTimeout:-1
     tag:0];
}


#pragma mark - GCDAsyncUdpSocket DELEGATE
/**
 * By design, UDP is a connectionless protocol, and connecting is not needed.
 * However, you may optionally choose to connect to a particular host for reasons
 * outlined in the documentation for the various connect methods listed above.
 *
 * This method is called if one of the connect methods are invoked, and the connection is successful.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didConnectToAddress:(NSData *)address {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didConnectToAddress:)]){
                [obj udpSocket:sock didConnectToAddress:address];
            }
        }];
    }
}

/**
 * By design, UDP is a connectionless protocol, and connecting is not needed.
 * However, you may optionally choose to connect to a particular host for reasons
 * outlined in the documentation for the various connect methods listed above.
 *
 * This method is called if one of the connect methods are invoked, and the connection fails.
 * This may happen, for example, if a domain name is given for the host and the domain name is unable to be resolved.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didNotConnect:(NSError *)error {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            [obj udpSocket:sock didNotConnect:error];
        }];
    }
}

/**
 * Called when the datagram with the given tag has been sent.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didSendDataWithTag:(long)tag {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didSendDataWithTag:)]){
                [obj udpSocket:sock didSendDataWithTag:tag];
            }
        }];
    }
}

/**
 * Called if an error occurs while trying to send a datagram.
 * This could be due to a timeout, or something more serious such as the data being too large to fit in a sigle packet.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didNotSendDataWithTag:(long)tag dueToError:(NSError *)error {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didNotSendDataWithTag:dueToError:)]){
                [obj udpSocket:sock didNotSendDataWithTag:tag dueToError:error];
            }
        }];
    }
}

/**
 * Called when the socket has received the requested datagram.
 **/
- (void)udpSocket:(GCDAsyncUdpSocket *)sock didReceiveData:(NSData *)data fromAddress:(NSData *)address withFilterContext:(id)filterContext {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocket:didReceiveData:fromAddress:withFilterContext:)]){
                [obj udpSocket:sock didReceiveData:data fromAddress:address withFilterContext:filterContext];
            }
        }];
    }
}

/**
 * Called when the socket is closed.
 **/
- (void)udpSocketDidClose:(GCDAsyncUdpSocket *)sock withError:(NSError *)error {
    @synchronized(self.multSockDelegates) {
        [self.multSockDelegates enumerateObjectsUsingBlock:^(id<GCDAsyncUdpSocketDelegate> _Nonnull obj, NSUInteger idx, BOOL * _Nonnull stop) {
            if ([obj respondsToSelector:@selector(udpSocketDidClose:withError:)]){
                [obj udpSocketDidClose:sock withError:error];
            }
        }];
    }
}

@end
