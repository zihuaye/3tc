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

.. _to-api-cachegroups-id:

**********************
``cachegroups/{{ID}}``
**********************
Extracts information about a single :term:`Cache Group`

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+--------------+---------------------------------------------------------------+
	| Parameter    | Description                                                   |
	+==============+===============================================================+
	| ID           | The integral, unique identifier of a :term:`Cache Group`      |
	+--------------+---------------------------------------------------------------+

Response Structure
------------------
:fallbackToClosest:   If ``true``, Traffic Router will direct clients to peers of this :term:`Cache Group` in the event that it becomes unavailable
:id:                  Integral, unique identifier for the :term:`Cache Group`
:lastUpdated:         The date and time at which this :term:`Cache Group` was last updated, in an ISO-like format
:latitude:            Latitude of the :term:`Cache Group`
:localizationMethods: An array of strings that name the localization methods enabled for this :term:`Cache Group`. Each of the three available localization methods may be present, with the following meanings:

	CZ
		Lookup in the Traffic Router's "Coverage Zone" file is enabled
	DEEP_CZ
		Lookup in the Traffic Router's "Deep Coverage Zone" file is enabled
	GEO
		Use of a geographical location-to-IP mapping database is enabled

:longitude:                     Longitude of the :term:`Cache Group`
:name:                          The name of the :term:`Cache Group`
:parentCachegroupId:            Integral, unique identifier of the :term:`Cache Group` that is this :term:`Cache Group`\ 's parent
:parentCachegroupName:          The name of the :term:`Cache Group` that is this :term:`Cache Group`\ 's parent
:secondaryParentCachegroupId:   Integral, unique identifier of the :term:`Cache Group` that is this :term:`Cache Group`\ 's secondary parent
:secondaryParentCachegroupName: The name of the :term:`Cache Group` that is this :term:`Cache Group`\ 's secondary parent
:shortName:                     Abbreviation of the :term:`Cache Group` Name
:typeId:                        The integral, unique identifier for the 'Type' of :term:`Cache Group`
:typeName:                      The name of the type of this :term:`Cache Group`

.. note:: The default value of ``fallbackToClosest`` is 'true', and if it is 'null' Traffic Control components will still interpret it as 'true'.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: EXO+TK1CIwQ5lzTXQGqlLDzU641pLLCQbyqz5Z8QUYSPAjjn5cqC9W3c0ioDiCdK9bUWvHP3E4/ERBzkBTi06g==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 18:35:53 GMT
	Content-Length: 357

	{ "response": [
		{
			"id": 8,
			"name": "test",
			"shortName": "test",
			"latitude": 0,
			"longitude": 0,
			"parentCachegroupName": "CDN_in_a_Box_Mid",
			"parentCachegroupId": 6,
			"secondaryParentCachegroupName": null,
			"secondaryParentCachegroupId": null,
			"fallbackToClosest": null,
			"localizationMethods": [
				"DEEP_CZ",
				"CZ"
			],
			"typeName": "EDGE_LOC",
			"typeId": 23,
			"lastUpdated": "2018-11-14 18:23:33+00"
		}
	]}


``PUT``
=======
Update :term:`Cache Group`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+--------------+---------------------------------------------------------------+
	| Parameter    | Description                                                   |
	+==============+===============================================================+
	| ID           | The integral, unique identifier of a :term:`Cache Group`      |
	+--------------+---------------------------------------------------------------+

:fallbackToClosest: An optional field which, if present and ``true``, will cause Traffic Router to direct clients to peers of this :term:`Cache Group` in the event that it becomes unavailable

	.. note:: The default value of ``fallbackToClosest`` is ``true``, and if it is ``null`` or ``undefined`` Traffic Control components will still interpret it as ``true``.

:latitude:            An optional field which, if specified, will set the latitude of the new :term:`Cache Group`\ [1]_
:localizationMethods: An optional array of strings that name the localization methods enabled for this :term:`Cache Group`. Each of the three available localization methods may be present, with the following meanings:

	CZ
		Lookup in the Traffic Router's "Coverage Zone" file will be enabled
	DEEP_CZ
		Lookup in the Traffic Router's "Deep Coverage Zone" file will be enabled
	GEO
		Use of a geographical location-to-IP mapping database will be enabled

:longitude:                 An optional field which, if specified, will set the longitude of the new :term:`Cache Group`\ [1]_
:name:                      The desired name of the :term:`Cache Group` entry
:parentCachegroup:          An optional field which, if specified, should be the integral, unique identifier of :term:`Cache Group` to use as the new :term:`Cache Group`\ 's parent
:secondaryParentCachegroup: An optional field which, if specified, should be the integral, unique identifier of :term:`Cache Group` to use as the new :term:`Cache Group`\ 's parent
:shortName:                 A more human-friendly abbreviation of the :term:`Cache Group`\ 's name
:typeId:                    The integral, unique identifier of the desired type of the new :term:`Cache Group` - by default the valid options are: "EDGE_LOC", "MID_LOC" or "ORG_LOC"

	.. note:: Rather than the actual name of the type, be sure to use the "database ID" of the desired type. Typically this will require looking up the types via the API first, as the IDs of even these default types is not deterministic.

.. code-block:: http
	:caption: Request Example

	PUT /api/1.3/cachegroups/8 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 118
	Content-Type: application/json

	{"latitude": 0.0, "longitude": 0.0, "name": "test", "shortName": "test", "typeId": 23, "localizationMethods": ["GEO"]}

Response Structure
------------------
:fallbackToClosest:   If ``true``, Traffic Router will direct clients to peers of this :term:`Cache Group` in the event that it becomes unavailable
:id:                  Integral, unique identifier for the :term:`Cache Group`
:lastUpdated:         The date and time at which this :term:`Cache Group` was last updated, in an ISO-like format
:latitude:            Latitude of the :term:`Cache Group`
:localizationMethods: An array of strings that name the localization methods enabled for this :term:`Cache Group`. Each of the three available localization methods may be present, with the following meanings:

	CZ
		Lookup in the Traffic Router's "Coverage Zone" file is enabled
	DEEP_CZ
		Lookup in the Traffic Router's "Deep Coverage Zone" file is enabled
	GEO
		Use of a geographical location-to-IP mapping database is enabled

:longitude:                     Longitude of the :term:`Cache Group`
:name:                          The name of the :term:`Cache Group`
:parentCachegroupId:            Integral, unique identifier of the :term:`Cache Group` that is this :term:`Cache Group`\ 's parent
:parentCachegroupName:          The name of the :term:`Cache Group` that is this :term:`Cache Group`\ 's parent
:secondaryParentCachegroupId:   Integral, unique identifier of the :term:`Cache Group` that is this :term:`Cache Group`\ 's secondary parent
:secondaryParentCachegroupName: The name of the :term:`Cache Group` that is this :term:`Cache Group`\ 's secondary parent
:shortName:                     Abbreviation of the :term:`Cache Group` Name
:typeId:                        The integral, unique identifier for the 'Type' of :term:`Cache Group`
:typeName:                      The name of the type of this :term:`Cache Group`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: t1W65/2kj25QyHt0Ib0xpBaAR2sXu2kOsRZ49WjKZp/AK5S1YWhX7VNWCuUGiN1VNM4QRNqODC/7ewhYDFUncA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 19:14:28 GMT
	Content-Length: 385

	{ "alerts": [
		{
			"text": "cg was updated.",
			"level": "success"
		}
	],
	"response": {
		"id": 8,
		"name": "test",
		"shortName": "test",
		"latitude": 0,
		"longitude": 0,
		"parentCachegroupName": null,
		"parentCachegroupId": null,
		"secondaryParentCachegroupName": null,
		"secondaryParentCachegroupId": null,
		"fallbackToClosest": null,
		"localizationMethods": [
			"GEO"
		],
		"typeName": null,
		"typeId": 23,
		"lastUpdated": "2018-11-14 19:14:28+00"
	}}

.. [1] While these fields are technically optional, note that if they are not specified many things may break. For this reason, Traffic Portal requires them when creating or editing :term:`Cache Group`\ s.

``DELETE``
==========
Delete :term:`Cache Group`. :term:`Cache Group`\ s which have assigned servers or child :term:`Cache Group`\ s cannot be deleted.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+--------------+------------------------------------------------------------------------+
	| Parameter    | Description                                                            |
	+==============+========================================================================+
	| ID           | The integral, unique identifier of a :term:`Cache Group` to be deleted |
	+--------------+------------------------------------------------------------------------+

Response Structure
------------------
.. code block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: 5jZBgO7h1eNF70J/cmlbi3Hf9KJPx+WLMblH/pSKF3FWb/10GUHIN35ZOB+lN5LZYCkmk3izGbTFkiruG8I41Q==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 20:31:04 GMT
	Content-Length: 57

	{ "alerts": [
		{
			"text": "cg was deleted.",
			"level": "success"
		}
	]}

