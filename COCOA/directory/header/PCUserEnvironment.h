//
//  PCUserEnvironment.h
//  SysUtil
//
//  Created by Almighty Kim on 10/24/16.
//  Copyright © 2016 PocketCluster. All rights reserved.
//

#ifndef __PCUSERENVIRONMENT_H__
#define __PCUSERENVIRONMENT_H__

extern const char*
PCEnvironmentCocoaHomeDirectory(void);

extern const char*
PCEnvironmentPosixHomeDirectory(void);

extern const char*
PCEnvironmentFullUserName(void);

extern const char*
PCEnvironmentUserTemporaryDirectory(void);

extern const char*
PCEnvironmentLoginUserName(void);

const char*
PCEnvironmentCurrentLanguageCode(void);

const char*
PCEnvironmentCurrentCountryCode(void);
#endif /* __PCUSERENVIRONMENT_H__ */
