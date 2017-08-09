//
//  keyboard.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __KEYBOARD_H__
#define __KEYBOARD_H__

HRESULT VboxKeyboardPutScancodes(IKeyboard* ckeyboard, PRUint32 scancodesCount, PRInt32* cscancodes, PRUint32* ccodesStored);

HRESULT VboxIKeyboardRelease(IKeyboard* ckeyboard);

#endif /* keyboard_h */
