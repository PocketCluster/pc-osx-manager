//
//  Node.h
//

/*
 * This is a faithful representation of "NodeStat" in "github.com/stkim1/pc-core/service/health"
 */

@interface Node : NSObject
@property (strong, nonatomic, readonly) NSString *Name;
@property (strong, nonatomic, readonly) NSString *MacAddr;
@property (nonatomic, readonly) BOOL Registered;
@property (nonatomic, readonly) BOOL Bounded;
@property (nonatomic, readonly) BOOL PcsshOn;
@property (nonatomic, readonly) BOOL OrchstOn;

- (instancetype) initWithDictionary:(NSDictionary *)aNodeDict;
- (BOOL) isReady;
@end
