//
//  NSString+Contain.h
//  ollehschoolhigh_ui
//
//  Created by Almighty Kim on 1/8/12.
//  Copyright (c) 2012 __MyCompanyName__. All rights reserved.
//

#import <Foundation/Foundation.h>
@interface NSString(ContainSubstring)
- (BOOL) containsString:(NSString *)inString
                options:(NSStringCompareOptions)inOptions;
- (BOOL) containsString:(NSString *)inString;
@end
