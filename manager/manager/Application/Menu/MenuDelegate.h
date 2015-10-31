//
//  MenuDelegate.h
//  Vagrant Manager
//
//  Copyright (c) 2014 Lanayo. All rights reserved.
//

#import "VagrantInstance.h"
#import "VagrantMachine.h"

@protocol MenuDelegate <NSObject>
- (void)performVagrantAction:(NSString*)action withInstance:(VagrantInstance*)instance;
- (void)performVagrantAction:(NSString*)action withMachine:(VagrantMachine*)machine;
@end