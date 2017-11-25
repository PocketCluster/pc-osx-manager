//
//  StepStage.h
//  manager
//
//  Created by Almighty Kim on 8/17/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

//typedef void (^ProgressCompletion)(NSDictionary *aResult);

// shouldControl -> didControl matching pair

@protocol StepControl;

@protocol StageStep <NSObject>
@property (nonatomic, weak) NSObject<StepControl> *stageControl;
- (void) control:(NSObject<StepControl> *)aControl askedProgressWithParam:(NSDictionary *)aParam;
- (void) didControl:(NSObject<StepControl> *)aControl progressedFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult;

- (void) control:(NSObject<StepControl> *)aControl askedRevertWithParam:(NSDictionary *)aParam;
- (void) didControl:(NSObject<StepControl> *)aControl revertedFrom:(NSObject<StageStep> *)aStep withResult:(NSDictionary *)aResult;
@end

@protocol StepControl <NSObject>
-(void)shouldControlProgressFrom:(NSObject<StageStep> *)aStep withParam:(NSDictionary *)aParam;
-(void)shouldControlRevertFrom:(NSObject<StageStep> *)aStep withParam:(NSDictionary *)aParam;
@end

