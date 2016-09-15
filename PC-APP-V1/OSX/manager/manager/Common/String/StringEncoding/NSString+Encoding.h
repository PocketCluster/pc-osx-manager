//
//  NSString+Encoding.h
//  ollehschoolhigh_ui
//
//  Created by Almighty Kim on 1/31/12.
//  Copyright (c) 2012 __MyCompanyName__. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface NSString(Encoding)
- (NSString *)encodeString:(NSStringEncoding)encoding;
- (NSString *)urlencode;
@end
