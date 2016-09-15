//
//  NSDictionary(HTTPAddition).h
//  Invitation
//
//  Created by Almighty Kim on 3/29/13.
//  Copyright (c) 2013 Colorful Glue. All rights reserved.
//

#import <Foundation/Foundation.h>

NSString*
StringByAddingPercentEscapesForURLArgument(NSString *string);

@interface NSDictionary(HTTPAddition)
-(NSString*)HTTPArgumentsString;
@end
