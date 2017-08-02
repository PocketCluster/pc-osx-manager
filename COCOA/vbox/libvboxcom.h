//
//  libvboxcom.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/12/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//


#ifndef __LIBVBOXCOM_H__
#define __LIBVBOXCOM_H__

typedef void** VOID_DPTR;

typedef enum VBRESULT {
    GOOD = 0,
    INFO,
    FATAL
} VBRESULT;

typedef void VBoxGlue;

#pragma mark init & close
VBRESULT
NewVBoxGlue(VBoxGlue**);

VBRESULT
CloseVBoxGlue(VBoxGlue* glue);

#pragma mark app & api version
extern unsigned int VBoxAppVersion(void);
extern unsigned int VBoxApiVersion(void);

#endif /* __LIBVBOXCOM_H__ */
