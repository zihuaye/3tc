..
..
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
..
..     http://www.apache.org/licenses/LICENSE-2.0
..
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
..

.. _to-api-deliveryservices-id-unassigned_servers:

**********************************************
``deliveryservices/{{ID}}/unassigned_servers``
**********************************************

.. caution:: This route does not appear to work properly, and its use is strongly discouraged! Also note that the documentation here is not being updated as a result of this, and may contain out-of-date and/or erroneous information.

``GET``
=======
Retrieves properties of Edge-tier servers not assigned to a :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [1]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------+
	| Name | Description                                       |
	+======+===================================================+
	| ID   | Delivery service ID.                              |
	+------+---------------------------------------------------+

Response Structure
------------------
:cachegroup:     The cache group name
:cachegroupId:   The cache group id
:cdnId:          Id of the CDN to which the server belongs to
:cdnName:        Name of the CDN to which the server belongs to
:domainName:     The domain name part of the FQDN of the cache
:guid:           An identifier used to uniquely identify the server
:hostName:       The host name part of the cache
:httpsPort:      The HTTPS port on which the main application listens (443 in most cases)
:id:             The server id (database row number
:iloIpAddress:   The IPv4 address of the lights-out-management port
:iloIpGateway:   The IPv4 gateway address of the lights-out-management port
:iloIpNetmask:   The IPv4 netmask of the lights-out-management port
:iloPassword:    The password of the of the lights-out-management user (displays as ****** unless you are an 'admin' user)
:iloUsername:    The user name for lights-out-management
:interfaceMtu:   The Maximum Transmission Unit (MTU) to configure for ``interfaceName``
:interfaceName:  The network interface name used for serving traffic
:ip6Address:     The IPv6 address/netmask for ``interfaceName``
:ip6Gateway:     The IPv6 gateway for ``interfaceName``
:ipAddress:      The IPv4 address for ``interfaceName``
:ipGateway:      The IPv4 gateway for ``interfaceName``
:ipNetmask:      The IPv4 netmask for ``interfaceName``
:lastUpdated:    The Time and Date for the last update for this server
:mgmtIpAddress:  The IPv4 address of the management port (optional
:mgmtIpGateway:  The IPv4 gateway of the management port (optional
:mgmtIpNetmask:  The IPv4 netmask of the management port (optional
:offlineReason:  A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status
:physLocation:   The physical location name
:physLocationId: The physical location id
:profile:        The assigned profile name
:profileDesc:    The assigned profile description
:profileId:      The assigned profile Id
:rack:           A string indicating rack location
:routerHostName: The human readable name of the router
:routerPortName: The human readable name of the router port
:status:         The Status string
:statusId:       The Status id
:tcpPort:        The default TCP port on which the main application listens (80 for a cache in most cases
:type:           The name of the type of this server
:typeId:         The id of the type of this server
:updPending:     bool

.. code-block:: json
	:caption: Response Example

	 {
			"response": [
					{
							"cachegroup": "us-il-chicago",
							"cachegroupId": "3",
							"cdnId": "3",
							"cdnName": "CDN-1",
							"domainName": "chi.kabletown.net",
							"guid": null,
							"hostName": "atsec-chi-00",
							"id": "19",
							"iloIpAddress": "172.16.2.6",
							"iloIpGateway": "172.16.2.1",
							"iloIpNetmask": "255.255.255.0",
							"iloPassword": "********",
							"iloUsername": "",
							"interfaceMtu": "9000",
							"interfaceName": "bond0",
							"ip6Address": "2033:D0D0:3300::2:2/64",
							"ip6Gateway": "2033:D0D0:3300::2:1",
							"ipAddress": "10.10.2.2",
							"ipGateway": "10.10.2.1",
							"ipNetmask": "255.255.255.0",
							"lastUpdated": "2015-03-08 15:57:32",
							"mgmtIpAddress": "",
							"mgmtIpGateway": "",
							"mgmtIpNetmask": "",
							"offlineReason": "N/A",
							"physLocation": "plocation-chi-1",
							"physLocationId": "9",
							"profile": "EDGE1_CDN1_421_SSL",
							"profileDesc": "EDGE1_CDN1_421_SSL profile",
							"profileId": "12",
							"rack": "RR 119.02",
							"routerHostName": "rtr-chi.kabletown.net",
							"routerPortName": "2",
							"status": "ONLINE",
							"statusId": "6",
							"tcpPort": "80",
							"httpsPort": "443",
							"type": "EDGE",
							"typeId": "3",
							"updPending": false
					},
				]
		}

.. [1] Users with the roles "admin" and/or "operations" will be able to see servers not assigned to *any* given :term:`Delivery Service`, whereas any other user will only be able to see the servers not assigned to :term:`Delivery Service`\ s their Tenant is allowed to see.

