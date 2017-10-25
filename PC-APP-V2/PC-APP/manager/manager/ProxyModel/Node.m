//
//  Node.m
//

#import "Node.h"

@interface Node()
@property (strong, nonatomic, readwrite) NSString *Name;
@property (strong, nonatomic, readwrite) NSString *MacAddr;
@property (nonatomic, readwrite) BOOL Registered;
@property (nonatomic, readwrite) BOOL Bounded;
@property (nonatomic, readwrite) BOOL PcsshOn;
@property (nonatomic, readwrite) BOOL OrchstOn;
@end

@implementation Node
@synthesize Name;
@synthesize MacAddr;
@synthesize Registered;
@synthesize Bounded;
@synthesize PcsshOn;
@synthesize OrchstOn;

- (instancetype)initWithDictionary:(NSDictionary *)aNodeDict {
    self = [super init];
    if(self){
        self.Name       = [aNodeDict valueForKey:@"name"];
        self.MacAddr    = [aNodeDict valueForKey:@"mac"];
        self.Registered = [[aNodeDict valueForKey:@"rgstd"] boolValue];
        self.Bounded    = [[aNodeDict valueForKey:@"bound"] boolValue];
        self.PcsshOn    = [[aNodeDict valueForKey:@"pcssh"] boolValue];
        self.OrchstOn   = [[aNodeDict valueForKey:@"orchst"] boolValue];
    }
    return self;
}

- (BOOL) isReady {
    return (self.Registered && self.Bounded && self.PcsshOn && self.OrchstOn);
}


@end
