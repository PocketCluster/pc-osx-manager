//
//  CustomCommandManager.m
//  Vagrant Manager
//
//  Copyright (c) 2015 Lanayo. All rights reserved.
//

@interface CustomCommand : NSObject

@property (strong, nonatomic) NSString *displayName;
@property (strong, nonatomic) NSString *command;
@property BOOL runInTerminal;

@end
