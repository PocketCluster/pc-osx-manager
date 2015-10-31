//
//  PCPrefWC.h
//  manager
//
//  Created by Almighty Kim on 10/30/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "BaseWindowController.h"

extern NSString * const kPCPrefDefaultTerm;

@interface PCPrefWC : BaseWindowController
@property (nonatomic, weak) IBOutlet NSPopUpButton *terminalPreferencePopUpButton;
@property (nonatomic, weak) IBOutlet NSButton *launchAtLoginCheckBox;

- (IBAction)launchAtLoginCheckBoxClicked:(id)sender;
- (IBAction)terminalPreferencePopUpButtonClicked:(id)sender;
@end
