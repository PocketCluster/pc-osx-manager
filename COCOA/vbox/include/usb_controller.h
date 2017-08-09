//
//  usb_controller.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __USB_CONTROLLER_H__
#define __USB_CONTROLLER_H__

HRESULT VboxGetUsbControllerName(IUSBController* ccontroller, char** cname);

HRESULT VboxGetUsbControllerStandard(IUSBController* ccontroller, PRUint16* cstandard);

HRESULT VboxGetUsbControllerType(IUSBController* ccontroller, PRUint32* ctype);

HRESULT VboxIUSBControllerRelease(IUSBController* ccontroller);

HRESULT VboxMachineAddUsbController(IMachine* cmachine, char* cname, PRUint32 ccontrollerType, IUSBController** ccontroller);

#endif /* usb_controller_h */
