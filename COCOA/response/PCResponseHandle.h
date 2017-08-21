//
//  PCResponseHandle.h
//  static-core
//
//  Created by Almighty Kim on 3/24/17.
//  Copyright © 2017 PocketCluster. All rights reserved.
//

#ifndef __PC_RESPONSE_HANDLE_H__
#define __PC_RESPONSE_HANDLE_H__

extern void
PCFeedResponseForGet(char* path, char* response);

extern void
PCFeedResponseForPost(char* path, char* response);

extern void
PCFeedResponseForPut(char* path, char* response);

extern void
PCFeedResponseForDelete(char* path, char* response);

#endif // __PC_RESPONSE_HANDLE_H__
