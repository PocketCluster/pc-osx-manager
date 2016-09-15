//
//  SingletonTaskHandler.h
//  Vivee
//
//  Created by Almighty Kim on 12/25/11.
//  Copyright (c) 2011 __MyCompanyName__. All rights reserved.
//

#import <Foundation/Foundation.h>
//#import "SynthesizeSingleton.h"
@protocol SingletonTaskHandler <NSObject>
@optional
-(void)applicationAwaken;
-(void)applicationStarted;

-(void)applicationResignActive;
-(void)applicationEnterBackground;

-(void)applicationEnterForeground;
-(void)applicationBecomeActive;

-(void)applicationWillTerminate;
@end
