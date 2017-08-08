//
//  progress.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __PROGRESS_H__
#define __PROGRESS_H__

HRESULT VboxProgressWaitForCompletion(IProgress* cprogress, int timeout);

HRESULT VboxProgressGetPercent(IProgress* cprogress, PRUint32* cpercent);

HRESULT VboxProgressGetResultCode(IProgress* cprogress, PRInt32* code);

HRESULT VboxProgressGetResultInfo(IProgress* cprogress, char** cErrorMessage);

HRESULT VboxProgressRelease(IProgress* cprogress);

#endif /* progress_h */
