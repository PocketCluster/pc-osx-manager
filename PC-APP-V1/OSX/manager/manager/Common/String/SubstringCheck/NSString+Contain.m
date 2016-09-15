//
//  NSString+Contain.m
//  ollehschoolhigh_ui
//
//  Created by Almighty Kim on 1/8/12.
//  Copyright (c) 2012 __MyCompanyName__. All rights reserved.
//

#import "NSString+Contain.h"
#import "NullStringChecker.h"

@implementation NSString(ContainSubstring)
- (BOOL) containsString:(NSString *)inString
                options:(NSStringCompareOptions)inOptions 
{
	if(ISNULL_STRING(self))
		return NO;
	
	if(ISNULL_STRING(inString))
		return NO;
	
	NSRange retRange = [self rangeOfString:inString options:inOptions];
	return (retRange.location != NSNotFound);
}

- (BOOL) containsString:(NSString *)inString
{
	return [self containsString:inString options:NSCaseInsensitiveSearch];
}
@end
