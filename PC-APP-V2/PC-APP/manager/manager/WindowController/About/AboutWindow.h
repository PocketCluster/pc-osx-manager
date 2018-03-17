//
//  AboutWindowController.h
//  PocketCluster
//
//  Copyright (c) 2015,2017 PocketCluster. All rights reserved.
//

#import "BaseWindowController.h"

@interface AboutWindow : BaseWindowController
@property (nonatomic, weak) IBOutlet NSTextField *copyright;

-(IBAction)homepage:(id)sender;
@end
