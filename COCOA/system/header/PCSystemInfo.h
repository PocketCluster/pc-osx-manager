//
//  PCSystemInfo.h
//  SysUtil
//
//  Created by Almighty Kim on 10/27/16.
//  Copyright © 2016 PocketCluster. All rights reserved.
//

#ifndef __PCSYSTEMINFO_H__
#define __PCSYSTEMINFO_H__

extern unsigned long
PCSystemProcessorCount(void);

extern unsigned long
PCSystemActiveProcessorCount(void);

// memory size in Byte
extern unsigned long long
PCSystemPhysicalMemorySize(void);

extern unsigned long
PCSystemPhysicalCoreCount(void);

#endif /* __PCSYSTEMINFO_H__ */
