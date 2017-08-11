//
//  Node.h
//

@interface Node : NSObject

@property (strong, nonatomic) NSString *SlaveID;
@property (strong, nonatomic) NSString *Hardware;
@property (strong, nonatomic) NSString *NodeName;
@property (strong, nonatomic) NSString *IP4Address;
@property (strong, nonatomic) NSString *IP4Gateway;
@property (strong, nonatomic) NSString *USerMadeName;
@property (nonatomic, readonly) NSDate *LastAlive;

-(instancetype)initWithDictionary:(NSDictionary *)aDict;

@end
