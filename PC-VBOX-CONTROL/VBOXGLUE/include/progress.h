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

HRESULT VboxGetProgressPercent(IProgress* cprogress, PRUint32* cpercent);

HRESULT VboxGetProgressResultCode(IProgress* cprogress, PRInt32* code);

HRESULT VboxIProgressRelease(IProgress* cprogress);

#endif /* progress_h */
