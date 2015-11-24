// AFAppDotNetAPIClient.h
//
// Copyright (c) 2012 Mattt Thompson (http://mattt.me/)
// 
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

#import "PCFormulaClient.h"

static NSString * const PCGithubAPIBaseURLString = @"https://api.github.com/repos/pocketcluster/formulas/contents/";

static NSString * const PCGithubRawFileURLString = @"https://raw.githubusercontent.com/PocketCluster/formulas/master/";

static NSString * const PCWIPRawFileURLString = @"https://raw.githubusercontent.com/PocketCluster/WIP/master/";

@implementation PCFormulaClient
+ (instancetype)sharedGithubAPIClient {
    static PCFormulaClient *_sharedGithubAPIClient = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        _sharedGithubAPIClient = [[PCFormulaClient alloc] initWithBaseURL:[NSURL URLWithString:PCGithubAPIBaseURLString]];
        _sharedGithubAPIClient.securityPolicy = [AFSecurityPolicy policyWithPinningMode:AFSSLPinningModeNone];
        
    });
    return _sharedGithubAPIClient;
}

+ (instancetype)sharedGithubRawFileClient {
    static PCFormulaClient *_sharedGithubRawFileClient = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        _sharedGithubRawFileClient = [[PCFormulaClient alloc] initWithBaseURL:[NSURL URLWithString:PCGithubRawFileURLString]];
        _sharedGithubRawFileClient.securityPolicy = [AFSecurityPolicy policyWithPinningMode:AFSSLPinningModeNone];
        _sharedGithubRawFileClient.responseSerializer.acceptableContentTypes = [NSSet setWithArray:@[@"text/plain", @"text/plain",@"application/json"]];
    });
    return _sharedGithubRawFileClient;
}

+ (instancetype)sharedWIPRawFileClient {
    static PCFormulaClient *_sharedWIPRawFileClient = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        _sharedWIPRawFileClient = [[PCFormulaClient alloc] initWithBaseURL:[NSURL URLWithString:PCWIPRawFileURLString]];
        _sharedWIPRawFileClient.securityPolicy = [AFSecurityPolicy policyWithPinningMode:AFSSLPinningModeNone];
        _sharedWIPRawFileClient.responseSerializer.acceptableContentTypes = [NSSet setWithArray:@[@"text/plain", @"text/plain",@"application/json"]];
    });
    return _sharedWIPRawFileClient;
}

+ (AFURLSessionManager *)sharedDownloadManager {
    static AFURLSessionManager *_sharedDownloadManager = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        _sharedDownloadManager = [[AFURLSessionManager alloc] initWithSessionConfiguration:[NSURLSessionConfiguration defaultSessionConfiguration]];
    });
    return _sharedDownloadManager;
}

@end
