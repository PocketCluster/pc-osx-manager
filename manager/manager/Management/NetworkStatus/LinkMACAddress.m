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

#import "LinkMACAddress.h"

@implementation LinkMACAddress

@synthesize string;

/* This method takes some junk string as input and tries to format it as MAC address
 * separated by colons. Note that it won't check for valid characters, it will only
 * deal with the colons. The character validity is guaranteed via the MACAddressFormatter.
 */

- (LinkMACAddress*) sanitize {
  // Stripping all existing colons
  string = [[self string] stringByReplacingOccurrencesOfString:@":" withString:@""];
  // Adding fresh colons
  NSMutableString* formatted = [string mutableCopy];
  if ([formatted length] > 10) [formatted insertString:@":" atIndex:10];
  if ([formatted length] > 8) [formatted insertString:@":" atIndex:8];
  if ([formatted length] > 6) [formatted insertString:@":" atIndex:6];
  if ([formatted length] > 4) [formatted insertString:@":" atIndex:4];
  if ([formatted length] > 2) [formatted insertString:@":" atIndex:2];
  self.string = formatted;
  return self;
}

- (BOOL) valid {
  [self sanitize];
  return [string length] == 17;
}

@end
