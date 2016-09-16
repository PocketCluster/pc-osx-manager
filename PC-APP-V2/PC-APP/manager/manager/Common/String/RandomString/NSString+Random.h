//
//  NSString+Random.h
//  Simulator
//
//  Created by Sung Taek Kim on 1/20/15.
//  Copyright (c) 2015 KETI. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface NSString(Random)
+(instancetype)randomASCIIStringOfLengh:(NSUInteger)aLength;
+(instancetype)randomASCIICapStringOfLengh:(NSUInteger)aLength;
@end
