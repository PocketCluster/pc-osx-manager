//
//  NSBundle+DirectoryPath.m
//  Album0
//
//  Created by Almighty Kim on 2/3/12.
//  Copyright (c) 2012 __MyCompanyName__. All rights reserved.
//

#import "NSBundle+DirectoryPath.h"

@implementation NSBundle(DirectoryPath)
-(NSString*)resourcePath:(NSString *)inStrFile
{
	NSArray* path = [inStrFile componentsSeparatedByString:@"."];

	return [[NSBundle mainBundle] 
			pathForResource: [path objectAtIndex: 0] 
			ofType:[path objectAtIndex: 1]
			];
}

@end
