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

.. _to-api-cdns-name-snapshot:

**************************
``cdns/{{name}}/snapshot``
**************************
.. caution:: This page is a stub! Much of it may be missing or just downright wrong - it needs a lot of love from people with the domain knowledge required to update it.

``GET``
=======
Retrieves the *current* snapshot for a CDN, which represents the current *operating state* of the CDN, **not** the current *configuration* of the CDN. The contents of this snapshot are currently used by Traffic Monitor and Traffic Router.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------+
	| Name | Description                                                |
	+======+============================================================+
	| name | The name of the CDN for which a snapshot shall be returned |
	+------+------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/cdns/CDN-in-a-Box/snapshot HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:config: An object containing basic configurations on the actual CDN object

	:api.cache-control.max-age: A string containing an integer which specifies the value of ``max-age`` in the ``Cache-Control`` header of some HTTP responses, likely the Traffic Router API responses

		.. deprecated:: 1.1
			This field still exists for legacy compatibility reasons, but has no known use at the time of this writing

	:certificates.polling.interval: A string containing an integer which specifies the interval, in seconds, on which other Traffic Control components should check for updated SSL certificates
	:consistent.dns.routing:        A string containing a boolean which indicates whether DNS routing will use a consistent hashing method or "round-robin"

		"false"
			The "round-robin" method will be used to define DNS routing
		"true"
			A consistent hashing method will be used to define DNS routing

	:coveragezone.polling.interval:      A string containing an integer which specifies the interval, in seconds, on which Traffic Routers should check for a new Coverage Zone file
	:coveragezone.polling.url:           The URL where a Coverage Zone file may be requested by Traffic Routers
	:dnssec.dynamic.response.expiration: A string containing a number and unit suffix that specifies the length of time for which dynamic responses to DNSSEC lookup queries should remain valid
	:dnssec.enabled:                     A string that tells whether or not the CDN uses DNSSEC; one of:

		"false"
			DNSSEC is not used within this CDN
		"true"
			DNSSEC is used within this CDN

	:domain_name:                        The Top-Level Domain Name (TLD) served by the CDN
	:edge.dns.limit:                     This field is of unknown use, and may be remnants of a legacy system
	:edge.dns.routing:                   This field is of unknown use, and may be remnants of a legacy system
	:edge.http.limit:                    This field is of unknown use, and may be remnants of a legacy system
	:edge.http.routing:                  This field is of unknown use, and may be remnants of a legacy system
	:federationmapping.polling.interval: A string containing an integer which specifies the interval, in seconds, on which other Traffic Control components should check for new federation mappings
	:federationmapping.polling.url:      The URL where Traffic Control components can request federation mappings
	:geolocation.polling.interval:       A string containing an integer which specifies the interval, in seconds, on which other Traffic Control components should check for new IP-to-geographic-location mapping databases
	:geolocation.polling.url:            The URL where Traffic Control components can request IP-to-geographic-location mapping database files
	:keystore.maintenance.interval:      A string containing an integer which specifies the interval, in seconds, on which Traffic Routers should refresh their zone caches
	:neustar.polling.interval:           A string containing an integer which specifies the interval, in seconds, on which other Traffic Control components should check for new "Neustar" databases
	:neustar.polling.url:                The URL where Traffic Control components can request "Neustar" databases
	:soa:                                An object defining the Start of Authority (SOA) for the CDN's TLD (defined in ``domain_name``)

		:admin: The name of the administrator for this zone - i.e. the RNAME

			.. note:: This rarely represents a proper email address, unfortunately.

		:expire:  A string containing an integer that sets the number of seconds after which secondary name servers should stop answering requests for this zone if the master does not respond
		:minimum: A string containing an integer that sets the Time To Live (TTL) - in seconds - of the record for the purpose of negative caching
		:refresh: A string containing an integer that sets the number of seconds after which secondary name servers should query the master for the SOA record, to detect zone changes
		:retry:   A string containing an integer that sets the number of seconds after which secondary name servers should retry to request the serial number from the master if the master does not respond

			.. note:: :rfc:`1035` dictates that this should always be less than ``refresh``.

		.. seealso:: `The Wikipedia page on Start of Authority records <https://en.wikipedia.org/wiki/SOA_record>`_.

	:steeringmapping.polling.interval:       A string containing an integer which specifies the interval, in seconds, on which Traffic Control components should check for new steering mappings
	:ttls:                                   An object that contains keys which are types of DNS records that have values which are strings containing integers that specify the time for which a response to the specific type of record request should remain valid
	:zonemanager.cache.maintenance.interval: A configuration option for the ZoneManager Java class of Traffic Router
	:zonemanager.threadpool.scale:           A configuration option for the ZoneManager Java class of Traffic Router

:contentRouters: An object containing keys which are the (short) hostnames of the Traffic Routers that serve requests for :term:`Delivery Service`\ s in this CDN

	:api.port:  A string containing the port number on which the :ref:`tr-api` is served by this Traffic Router
	:fqdn:      This Traffic Router's Fully Qualified Domain Name (FQDN)
	:httpsPort: The port number on which this Traffic Router listens for incoming HTTPS requests
	:ip:        This Traffic Router's IPv4 address
	:ip6:       This Traffic Router's IPv6 address
	:location:  The name of the Cache Group to which this Traffic Router belongs
	:port:      The port number on which this Traffic Router listens for incoming HTTP requests
	:profile:   The name of the profile used by this Traffic Router
	:status:    The health status of this Traffic Router

		.. seealso:: :ref:`health-proto`

:contentServers: An object containing keys which are the (short) hostnames of the Edge-Tier :term:`cache server` s in the CDN; the values corresponding to those keys are routing information for said servers

	:cacheGroup:       The name of the Cache Group to which the server belongs
	:deliveryServices: An object containing keys which are the names of :term:`Delivery Service`\ s to which this :term:`cache server` is assigned; the values corresponding to those keys are arrays of FQDNs that resolve to this :term:`cache server`

		.. note:: Only Edge-tier :term:`cache server` s can be assigned to a Delivery SErvice, and therefore this field will only be present when ``type`` is ``"EDGE"``.

	:fqdn:            The server's Fully Qualified Domain Name (FQDN)
	:hashCount:       The number of servers to be placed into a single "hash ring" in Traffic Router
	:hashId:          A unique string to be used as the key for hashing servers - as of version 3.0.0 of Traffic Control, this is always the same as the server's (short) hostname and only still exists for legacy compatibility reasons
	:httpsPort:       The port on which the :term:`cache server` listens for incoming HTTPS requests
	:interfaceName:   The name of the main network interface device used by this :term:`cache server`
	:ip6:             The server's IPv6 address
	:ip:              The server's IPv4 address
	:locationId:      This field is exactly the same as ``cacheGroup`` and only exists for legacy compatibility reasons
	:port:            The port on which this :term:`cache server` listens for incoming HTTP requests
	:profile:         The name of the profile used by the :term:`cache server`
	:routingDisabled: An integer representing the boolean concept of whether or not Traffic Routers should route client traffic this :term:`cache server`; one of:

		0
			Do not route traffic to this server
		1
			Route traffic to this server normally

	:status: This :term:`cache server`'s status

		.. seealso:: :ref:`health-proto`

	:type: The type of this :term:`cache server`; one of:

		EDGE
			This is an Edge-tier :term:`cache server`
		MID
			This is a Mid-tier :term:`cache server`

:deliveryServices: An object containing keys which are the 'xml_id's of all of the :term:`Delivery Service`\ s within the CDN

	:anonymousBlockingEnabled: A string containing a boolean that tells whether or not Anonymized IP Addresses are blocked by this :term:`Delivery Service`; one of:

		"true"
			Anonymized IP addresses are blocked by this :term:`Delivery Service`
		"false"
			Anonymized IP addresses are not blocked by this :term:`Delivery Service`

		.. seealso:: :ref:`anonymous_blocking-qht`

	:coverageZoneOnly: A string containing a boolean that tells whether or not this :term:`Delivery Service` routes traffic based only on its Coverage Zone file
	:deepCachingType:  A string that tells when Deep Caching is used by this :term:`Delivery Service`; one of:

		"ALWAYS"
			Deep Caching is always used by this :term:`Delivery Service`
		"NEVER"
			Deep Caching is never used by this :term:`Delivery Service`

	:dispersion: An object describing the "dispersion" - or number of caches within a single Cache Group across which the same content is spread - within the :term:`Delivery Service`

		:limit: The maximum number of caches in which the response to a single request URL will be stored

			.. note:: If this is greater than the number of caches in the Cache Group chosen to service the request, then content will be spread across all of them. That is, it causes no problems.

		:shuffled: A string containing a boolean that tells whether the caches chosen for content dispersion are chosen randomly or based on a consistent hash of the request URL; one of:

			"false"
				Caches will be chosen consistently
			"true"
				Caches will be chosen at random

	:domains:             An array of domains served by this :term:`Delivery Service`
	:geolocationProvider: The name of a provider for IP-to-geographic-location mapping services - currently the only valid value is ``"maxmindGeolocationService"``
	:ip6RoutingEnabled:   A string containing a boolean that tells whether IPv6 traffic can be routed on this :term:`Delivery Service`; one of:

		"false"
			IPv6 traffic will not be routed by this :term:`Delivery Service`
		"true"
			IPv6 traffic will be routed by this :term:`Delivery Service`

	:matchList: An array of methods used by Traffic Router to determine whether or not a request can be serviced by this :term:`Delivery Service`

		:pattern:   A regular expression - the use of this pattern is dependent on the ``type`` field (backslashes are escaped)
		:setNumber: An integral, unique identifier for the set of types to which the ``type`` field belongs
		:type:      The type of match performed using ``pattern`` to determine whether or not to use this :term:`Delivery Service`

			HOST_REGEXP
				Use the :term:`Delivery Service` if ``pattern`` matches the ``Host:`` HTTP header of an HTTP request\ [1]_
			HEADER_REGEXP
				Use the :term:`Delivery Service` if ``pattern`` matches an HTTP header (both the name and value) in an HTTP request\ [1]_
			PATH_REGEXP
				Use the :term:`Delivery Service` if ``pattern`` matches the request path of this :term:`Delivery Service`'s URL
			STEERING_REGEXP
				Use the :term:`Delivery Service` if ``pattern`` matches the ``xml_id`` of one of this :term:`Delivery Service`'s "Steering" target :term:`Delivery Service`\ s

	:missLocation: An object representing the default geographic coordinates to use for a client when lookup of their IP has failed in both the Coverage Zone file(s) and the IP-to-geographic-location database

		:lat:  Geographic latitude
		:long: Geographic longitude

	:protocol: An object that describes how the :term:`Delivery Service` ought to handle HTTP requests both with and without TLS encryption

		:acceptHttps: A string containing a boolean that tells whether HTTPS requests should be normally serviced by this :term:`Delivery Service`; one of:

			"false"
				Refuse to service HTTPS requests
			"true"
				Service HTTPS requests normally

		:redirectToHttps: A string containing a boolean that tells whether HTTP requests ought to be re-directed to use HTTPS; one of:

			"false"
				Do not redirect unencrypted traffic; service it normally
			"true"
				Respond to HTTP requests with instructions to use HTTPS instead

	:regionalGeoBlocking: A string containing a boolean that tells whether Regional Geographic Blocking is enabled on this :term:`Delivery Service`; one of:

		"false"
			Regional Geographic Blocking is not used by this :term:`Delivery Service`
		"true"
			Regional Geographic Blocking is used by this :term:`Delivery Service`

		.. seealso:: :ref:`regionalgeo-qht`

	:routingName: The highest-level part of the FQDNs serviced by this :term:`Delivery Service`
	:soa:         An object defining the Start of Authority (SOA) record for the :term:`Delivery Service`'s TLDs (defined in ``domains``)

		:admin: The name of the administrator for this zone - i.e. the RNAME

			.. note:: This rarely represents a proper email address, unfortunately.

		:expire:  A string containing an integer that sets the number of seconds after which secondary name servers should stop answering requests for this zone if the master does not respond
		:minimum: A string containing an integer that sets the Time To Live (TTL) - in seconds - of the record for the purpose of negative caching
		:refresh: A string containing an integer that sets the number of seconds after which secondary name servers should query the master for the SOA record, to detect zone changes
		:retry:   A string containing an integer that sets the number of seconds after which secondary name servers should retry to request the serial number from the master if the master does not respond

			.. note:: :rfc:`1035` dictates that this should always be less than ``refresh``.

		.. seealso:: `The Wikipedia page on Start of Authority records <https://en.wikipedia.org/wiki/SOA_record>`_.

	:sslEnabled: A string containing a boolean that tells whether this :term:`Delivery Service` uses SSL; one of:

		"false"
			SSL is not used by this :term:`Delivery Service`
		"true"
			SSL is used by this :term:`Delivery Service`

	:ttls: An object that contains keys which are types of DNS records that have values which are strings containing integers that specify the time for which a response to the specific type of record request should remain valid

		.. note:: This overrides ``config.ttls``.

:edgeLocations: An object containing keys which are the names of Edge-Tier Cache Groups within the CDN

	:backupLocations: An object that describes fallbacks for when this Cache Group is unavailable

		:fallbackToClosest: A string containing a boolean which tells whether requests should fall back on the closest available Cache Group when this Cache Group is not available; one of:

			"false"
				Do not fall back on the closest available Cache Group
			"true"
				Fall back on the closest available Cache Group

		:list: If any fallback Cache Groups have been configured for this Cache Group, this key will appear and will be an array of the names of all of those fallback Cache Groups, in the prescribed order

	:latitude:            The geographic latitude of this Cache Group
	:localizationMethods: An array of short names for localization methods available for this Cache Group
	:longitude:           The geographic longitude of this Cache Group

:monitors: An object containing keys which are the (short) hostnames of Traffic Monitors within this CDN

	:fqdn:      The FQDN of this Traffic Monitor
	:httpsPort: The port number on which this Traffic Monitor listens for incoming HTTPS requests
	:ip6:       This Traffic Monitor's IPv6 address
	:ip:        This Traffic Monitor's IPv4 address
	:location:  The name of the Cache Group to which this Traffic Monitor belongs
	:port:      The port number on which this Traffic Monitor listens for incoming HTTP requests
	:profile:   The name of the profile used by this Traffic Monitor

		.. note:: For legacy reasons, this must always start with "RASCAL-".

	:status: The health status of this Traffic Monitor

		.. seealso:: :ref:`health-proto`

:stats: An object containing metadata information regarding the CDN

	:CDN_name: The name of this CDN
	:date:     The UNIX epoch timestamp date in the Traffic Ops server's own timezone
	:tm_host:  The FQDN of the Traffic Ops server
	:tm_path:  A path relative to the root of the Traffic Ops server where a request may be replaced to have this snapshot overwritten by the current *configured state* of the CDN

		.. deprecated:: 1.1
			This field is still present for legacy compatibility reasons, but its contents should be ignored. Instead, make a ``PUT`` request to :ref:`to-api-snapshot-name`.

	:tm_user:    The username of the currently logged-in user
	:tm_version: The full version number of the Traffic Ops server, including release number, git commit hash, and supported Enterprise Linux version

:trafficRouterLocations: An object containing keys which are the names of Cache Groups within the CDN which contain Traffic Routers

	:backupLocations: An object that describes fallbacks for when this Cache Group is unavailable

		:fallbackToClosest: A string containing a boolean which tells whether requests should fall back on the closest available Cache Group when this Cache Group is not available; one of:

			"false"
				Do not fall back on the closest available Cache Group
			"true"
				Fall back on the closest available Cache Group

	:latitude:            The geographic latitude of this Cache Group
	:localizationMethods: An array of short names for localization methods available for this Cache Group
	:longitude:           The geographic longitude of this Cache Group

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: 220bc4XXwaj+s7ODd3QAF5leGj06lnApiN5E8H/B2RgxSphnQIfnwy6WWbBDjonWXPV1IWDCjBMO+rR+lAabMg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 12 Dec 2018 17:36:25 GMT
	Transfer-Encoding: chunked

	{ "response": {
		"config": {
			"api.cache-control.max-age": "10",
			"certificates.polling.interval": "300000",
			"consistent.dns.routing": "true",
			"coveragezone.polling.interval": "3600000",
			"coveragezone.polling.url": "https://trafficops.infra.ciab.test:443/coverage-zone.json",
			"dnssec.dynamic.response.expiration": "300s",
			"dnssec.enabled": "false",
			"domain_name": "mycdn.ciab.test",
			"edge.dns.limit": "6",
			"edge.dns.routing": "true",
			"edge.http.limit": "6",
			"edge.http.routing": "true",
			"federationmapping.polling.interval": "60000",
			"federationmapping.polling.url": "https://${toHostname}/internal/api/1.3/federations.json",
			"geolocation.polling.interval": "86400000",
			"geolocation.polling.url": "https://trafficops.infra.ciab.test:443/GeoLite2-City.mmdb.gz",
			"keystore.maintenance.interval": "300",
			"neustar.polling.interval": "86400000",
			"neustar.polling.url": "https://trafficops.infra.ciab.test:443/neustar.tar.gz",
			"soa": {
				"admin": "twelve_monkeys",
				"expire": "604800",
				"minimum": "30",
				"refresh": "28800",
				"retry": "7200"
			},
			"steeringmapping.polling.interval": "60000",
			"ttls": {
				"A": "3600",
				"AAAA": "3600",
				"DNSKEY": "30",
				"DS": "30",
				"NS": "3600",
				"SOA": "86400"
			},
			"zonemanager.cache.maintenance.interval": "300",
			"zonemanager.threadpool.scale": "0.50"
		},
		"contentServers": {
			"edge": {
				"cacheGroup": "CDN_in_a_Box_Edge",
				"fqdn": "edge.infra.ciab.test",
				"hashCount": 999,
				"hashId": "edge",
				"httpsPort": 443,
				"interfaceName": "eth0",
				"ip": "172.16.239.100",
				"ip6": "fc01:9400:1000:8::100",
				"locationId": "CDN_in_a_Box_Edge",
				"port": 80,
				"profile": "ATS_EDGE_TIER_CACHE",
				"status": "REPORTED",
				"type": "EDGE",
				"deliveryServices": {
					"demo1": [
						"edge.demo1.mycdn.ciab.test"
					]
				},
				"routingDisabled": 0
			},
			"mid": {
				"cacheGroup": "CDN_in_a_Box_Mid",
				"fqdn": "mid.infra.ciab.test",
				"hashCount": 999,
				"hashId": "mid",
				"httpsPort": 443,
				"interfaceName": "eth0",
				"ip": "172.16.239.120",
				"ip6": "fc01:9400:1000:8::120",
				"locationId": "CDN_in_a_Box_Mid",
				"port": 80,
				"profile": "ATS_MID_TIER_CACHE",
				"status": "REPORTED",
				"type": "MID",
				"routingDisabled": 0
			}
		},
		"contentRouters": {
			"trafficrouter": {
				"api.port": "3333",
				"fqdn": "trafficrouter.infra.ciab.test",
				"httpsPort": 443,
				"ip": "172.16.239.60",
				"ip6": "fc01:9400:1000:8::60",
				"location": "CDN_in_a_Box_Edge",
				"port": 80,
				"profile": "CCR_CIAB",
				"status": "ONLINE"
			}
		},
		"deliveryServices": {
			"demo1": {
				"anonymousBlockingEnabled": "false",
				"coverageZoneOnly": "false",
				"dispersion": {
					"limit": 1,
					"shuffled": "true"
				},
				"domains": [
					"demo1.mycdn.ciab.test"
				],
				"geolocationProvider": "maxmindGeolocationService",
				"matchsets": [
					{
						"protocol": "HTTP",
						"matchlist": [
							{
								"regex": ".*\\.demo1\\..*",
								"match-type": "HOST"
							}
						]
					}
				],
				"missLocation": {
					"lat": 42,
					"long": -88
				},
				"protocol": {
					"acceptHttps": "false",
					"redirectToHttps": "false"
				},
				"regionalGeoBlocking": "false",
				"soa": {
					"admin": "traffic_ops",
					"expire": "604800",
					"minimum": "30",
					"refresh": "28800",
					"retry": "7200"
				},
				"sslEnabled": "false",
				"ttls": {
					"A": "",
					"AAAA": "",
					"NS": "3600",
					"SOA": "86400"
				},
				"ip6RoutingEnabled": "true",
				"routingName": "video",
				"deepCachingType": "NEVER"
			}
		},
		"edgeLocations": {
			"CDN_in_a_Box_Edge": {
				"latitude": 38.897663,
				"longitude": -77.036574,
				"backupLocations": {
					"fallbackToClosest": "true"
				},
				"localizationMethods": [
					"GEO",
					"CZ",
					"DEEP_CZ"
				]
			}
		},
		"trafficRouterLocations": {
			"CDN_in_a_Box_Edge": {
				"latitude": 38.897663,
				"longitude": -77.036574,
				"backupLocations": {
					"fallbackToClosest": "false"
				},
				"localizationMethods": [
					"GEO",
					"CZ",
					"DEEP_CZ"
				]
			}
		},
		"monitors": {
			"trafficmonitor": {
				"fqdn": "trafficmonitor.infra.ciab.test",
				"httpsPort": 443,
				"ip": "172.16.239.40",
				"ip6": "fc01:9400:1000:8::40",
				"location": "CDN_in_a_Box_Edge",
				"port": 80,
				"profile": "RASCAL-Traffic_Monitor",
				"status": "ONLINE"
			}
		},
		"stats": {
			"CDN_name": "CDN-in-a-Box",
			"date": 1544635937,
			"tm_host": "trafficops.infra.ciab.test",
			"tm_path": "/tools/write_crconfig/CDN-in-a-Box",
			"tm_user": "admin",
			"tm_version": "traffic_ops-3.0.0-9813.8ad7bd8e.el7"
		}
	}}

.. [1] These only apply to HTTP-routed :term:`Delivery Service`\ s
