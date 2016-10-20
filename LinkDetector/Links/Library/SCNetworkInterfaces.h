//
//  netifaces.h
//  netifaces
//
//  Created by Sung-Taek, Kim on 10/2/16.
//  Copyright © 2016 PocketCluster. All rights reserved.
//

#ifndef __NETIFACES_H__
#define __NETIFACES_H__

#include <SystemConfiguration/SystemConfiguration.h>
#include <errno.h>
#include "SCNetworkTypes.h"

/*!
	@function SCNetworkInterfaceMediaStatus
	@discussion Returns if an interface is active.
	@result true if active, otherwise false.
 */
CF_EXPORT
bool SCNetworkInterfaceMediaStatus(SCNetworkInterfaceRef interface);

/*!
	@function SCNIMutableAddressArray
	@discussion Returns a new CF mutable array with callbacks for address.
	@result A mutable CF array;
 */

CF_EXPORT
CFMutableArrayRef SCNIMutableAddressArray(void);

/*!
	@function SCNetworkInterfaceAddresses
	@discussion Returns ipv4/v6 addresses (dotted format) linked to the interface.
	@param interface The network interface.
            A mutable CF array where addresses should be contained
	@result The list of ipv4/v6 addresses linked the interface;
            0 length CFArray if no ipv4/v6 addresses are supported or linked.
 */
CF_EXPORT
errno_t SCNetworkInterfaceAddresses(SCNetworkInterfaceRef interface, CFMutableArrayRef results);

/*!
	@function SCNetworkInterfaceAddressRelease
	@discussion Releases results addresses in a mutable array
	@param A mutable array containing all addresses
 */
CF_EXPORT
void SCNetworkInterfaceAddressRelease(CFMutableArrayRef results);

/*!
	@function SCNIMutableAddressArray
	@discussion Returns a new CF mutable array with callbacks for address.
	@result A mutable CF array;
 */

CF_EXPORT
CFMutableArrayRef SCNIMutableGatewayArray(void);

/*!
	@function SCNetworkInteraceGateways
	@discussion Returns all system gateways (dotted format).
	@result The list of ipv4/v6 in the system.
            0 length CFArray if no ipv4/ipv6 gateway exists.
 */
CF_EXPORT
errno_t SCNetworkGateways(CFMutableArrayRef results);

/*!
	@function SCNetworkInteraceGatewayRelease
	@discussion Releases results addresses in a mutable array
	@param A mutable array containing all addresses
 */
CF_EXPORT
void SCNetworkGatewayRelease(CFMutableArrayRef results);


#endif /* SCNetworkInterfaces */
