//
//  AgreementVC.m
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AgreementVC.h"

@interface AgreementVC ()

@end

@implementation AgreementVC

- (void)viewDidLoad {
    [super viewDidLoad];
    
    [self.agreement setVerticallyResizable:YES];
    [self.agreement setHorizontallyResizable:NO];
    [self.agreement setEditable:NO];
    
    [self.agreement.textStorage
     setAttributedString:[[NSAttributedString alloc]
                          initWithPath:[[NSBundle mainBundle] pathForResource:@"EULA" ofType:@"rtf"]
                          documentAttributes:nil]];
}

-(IBAction)agreed:(id)sender {
    
}

-(IBAction)declined:(id)sender {
    
}
@end
