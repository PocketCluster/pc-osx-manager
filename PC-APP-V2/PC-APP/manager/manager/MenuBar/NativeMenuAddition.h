//
//  NativeMenuAddition.h
//  manager
//
//  Created by Almighty Kim on 10/20/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#ifndef __NATIVEMENUADDITION_H__
#define __NATIVEMENUADDITION_H__

enum {
    MENUITEM_TOP_STATUS = 1,
    MENUITEM_PKG_DIV,
    MENUITEM_PREF,
    MENUITEM_UPDATE,
    MENUITEM_SLACK,
    MENUITEM_ABOUT,
    MENUITEM_DEBUG,
    MENUITEM_QUIT,
    MENUITEM_COUNT = MENUITEM_QUIT,
};

#define PKG_TAG_BUMPER     1000
#define PKG_TAG_BUILDER(x) (x + PKG_TAG_BUMPER)

// this tags matches 1:1 with 'ExecState' so that one can easily sort out tags
enum {
    EXEC_IDLE = 0,
    EXEC_STARTING,
    EXEC_STARTED,
    EXEC_RUN,
    EXEC_STOPPING
};

#endif /* NativeMenuAddition_h */
