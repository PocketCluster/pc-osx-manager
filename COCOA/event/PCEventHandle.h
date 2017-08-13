//
//  PCEventHandle.h
//  static-core
//
//  Created by Almighty Kim on 3/24/17.
//  Copyright © 2017 PocketCluster. All rights reserved.
//

#ifndef __PCEVENTHANDLE_H__
#define __PCEVENTHANDLE_H__

extern void
PCEventFeedGet(char* path);

extern void
PCEventFeedPost(char* path, char* payload);

extern void
PCEventFeedPut(char* path, char* payload);

extern void
PCEventFeedDelete(char* path);

#endif // __PCEVENTHANDLE_H__
