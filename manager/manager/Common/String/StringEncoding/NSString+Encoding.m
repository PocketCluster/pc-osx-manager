//
//  NSString+Encoding.m
//  ollehschoolhigh_ui
//
//  Created by Almighty Kim on 1/31/12.
//  Copyright (c) 2012 __MyCompanyName__. All rights reserved.
//

#import "NSString+Encoding.h"
#import <string.h>

static NSString * const URL_ESC_CHARS = @"!*'();:@&=+$,/?%#[]";

@implementation NSString(Encoding)
- (NSString *)encodeString:(NSStringEncoding)encoding
{
	// Encode all the reserved characters, per RFC 3986
	// (<http://www.ietf.org/rfc/rfc3986.txt>)
	CFStringRef cfStrEscaped = \
		CFURLCreateStringByAddingPercentEscapes
			(kCFAllocatorDefault,
             (CFStringRef)self,
			 NULL,
			 (CFStringRef)URL_ESC_CHARS,
			 CFStringConvertNSStringEncodingToEncoding(encoding));
	
	NSString *nsStrEscaped = \
		[NSString stringWithFormat:@"%@",(__bridge NSString *)cfStrEscaped];

	CFRelease(cfStrEscaped);

	return nsStrEscaped;
}

- (NSString *)urlencode
{
    NSMutableString *output = [NSMutableString string];
    const unsigned char *source = (const unsigned char *)[self UTF8String];
    int sourceLen = (int)strlen((const char *)source);
    for (int i = 0; i < sourceLen; ++i) {
        const unsigned char thisChar = source[i];
        if (thisChar == ' '){
            [output appendString:@"+"];
        } else if (thisChar == '.' || thisChar == '-' || thisChar == '_' || thisChar == '~' ||
                   (thisChar >= 'a' && thisChar <= 'z') ||
                   (thisChar >= 'A' && thisChar <= 'Z') ||
                   (thisChar >= '0' && thisChar <= '9')) {
            [output appendFormat:@"%c", thisChar];
        } else {
            [output appendFormat:@"%%%02X", thisChar];
        }
    }
    return output;
}

@end
