//
//  medium_format.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/8/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//

#ifndef __MEDIUM_FORMAT_H__
#define __MEDIUM_FORMAT_H__

HRESULT VboxGetMediumFormats(ISystemProperties* cprops, IMediumFormat*** cformats, ULONG* formatCount);

HRESULT VboxGetMediumFormatId(IMediumFormat* cformat, char** cid);

HRESULT VboxIMediumFormatRelease(IMediumFormat* cformat);

#endif /* medium_format_h */
