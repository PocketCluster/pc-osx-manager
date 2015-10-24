/*
 Copyright (c) 2015 funkensturm. https://github.com/halo/LinkLiar
 
 Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the
 "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish,
 distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to
 the following conditions:
 
 The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 
 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
 LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

#import <Foundation/Foundation.h>

@interface LinkInterface : NSObject

@property (nonatomic, strong) NSString *BSDName;
@property (nonatomic, strong) NSString *displayName;
@property (nonatomic, strong) NSString *hardMAC;
@property (nonatomic, strong) NSString *kind;
@property (nonatomic, strong) NSString *ip4Address;
@property (nonatomic, strong) NSString *ip6Address;

@property (readonly) NSString *softVendorName;
@property (readonly) NSString *softMAC;
@property (readonly) NSString *softMACLive;
@property (readonly) NSString *displayNameAndBSDName;
@property (readonly) NSInteger BSDNumber;

@property (readonly) BOOL hasOriginalMAC;

@end
