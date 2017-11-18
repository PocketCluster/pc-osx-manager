//
//  NSString+EmailForm.m
//  manager
//
//  Created by Almighty Kim on 11/14/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "NSString+EmailForm.h"

@implementation NSString(EmailForm)
- (BOOL) isValidEmailForm {
    // https://stackoverflow.com/questions/11760787/regex-for-email-address
    static NSString * const emailRegex = @"[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,63}";

    return [[NSPredicate
            predicateWithFormat:@"SELF MATCHES %@", emailRegex]
            evaluateWithObject:self];
}
@end
