//
//  storage_controller.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __STORAGE_CONTROLLER_H__
#define __STORAGE_CONTROLLER_H__

HRESULT VboxGetStorageControllerName(IStorageController* ccontroller, char** cname);

HRESULT VboxGetStorageControllerBus(IStorageController* ccontroller, PRUint32* cbus);

HRESULT VboxGetStorageControllerType(IStorageController* ccontroller, PRUint32* ctype);

HRESULT VboxSetStorageControllerType(IStorageController* ccontroller, PRUint32 ctype);

HRESULT VboxIStorageControllerRelease(IStorageController* ccontroller);

HRESULT VboxMachineAddStorageController(IMachine* cmachine, const char* cname, PRUint32 connectionType, IStorageController** ccontroller);

#endif /* storage_controller_h */
