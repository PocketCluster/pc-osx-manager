//
//  NSString+Random.m
//  Simulator
//
//  Created by Sung Taek Kim on 1/20/15.
//  Copyright (c) 2015 KETI. All rights reserved.
//

#import "NSString+Random.h"

@implementation NSString(Random)
+(instancetype)randomASCIIStringOfLengh:(NSUInteger)aLength
{
/*
    NSMutableString *str = [NSMutableString stringWithCapacity:aLength];
    for (NSUInteger i = 0; i < aLength; i++)
    {
        [str appendFormat:@"%c", 65 + (arc4random() % 27)];
        //[payload appendFormat:@"%c", (arc4random() % 256)];
    }
*/
    
    char *cstr = calloc(aLength,sizeof(char));
    for (NSUInteger i = 0; i < aLength; i++)
    {
        cstr[i] = 32 + (arc4random() % 94);
    }

    NSString *str = [NSString stringWithUTF8String:(const char*)cstr];
    free(cstr);
    return str;
}

+(instancetype)randomASCIICapStringOfLengh:(NSUInteger)aLength
{
    char *cstr = calloc(aLength,sizeof(char));
    for (NSUInteger i = 0; i < aLength; i++)
    {
        cstr[i] = 65 + (arc4random() % 26);
    }
    
    NSString *str = [NSString stringWithUTF8String:(const char*)cstr];
    free(cstr);
    return str;
}


@end
