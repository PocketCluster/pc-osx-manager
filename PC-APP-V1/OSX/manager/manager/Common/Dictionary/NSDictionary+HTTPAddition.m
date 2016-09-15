//
//  NSDictionary(HTTPAddition).m
//  Invitation
//
//  Created by Almighty Kim on 3/29/13.
//  Copyright (c) 2013 Colorful Glue. All rights reserved.
//

#import "NSDictionary+HTTPAddition.h"

NSString*
StringByAddingPercentEscapesForURLArgument(NSString *string)
{
	NSString *escapedString = \
	(__bridge NSString *)CFURLCreateStringByAddingPercentEscapes(kCFAllocatorDefault,
																 (CFStringRef)string,
																 NULL,
																 //(CFStringRef)@"!*'();:@&=+$,/?%#[]",
																 CFSTR(":/?#[]@!$ &'()*+,;=\"<>%{}|\\^~`"),
																 kCFStringEncodingUTF8);
	return escapedString;
}

@implementation NSDictionary(HTTPAddition)
-(NSString*)HTTPArgumentsString
{
	NSMutableArray *arguments = [NSMutableArray array];
    
	for (NSString *key in self)
	{
		NSString *parameter = \
		[NSString
		 stringWithFormat:@"%@=%@"
		 ,key
		 ,StringByAddingPercentEscapesForURLArgument([self objectForKey:key])];
		[arguments addObject:parameter];
	}
	
	return [arguments componentsJoinedByString:@"&"];
}
@end
