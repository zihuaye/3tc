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

.. _to-api-profiles-name-name-copy-copy:

****************************************
``profiles/name/{{name}}/copy/{{copy}}``
****************************************

``POST``
========
Copy profile to a new profile. The new profile name must not exist.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------+
	| Name | Description                                          |
	+======+======================================================+
	| name | The name of the new profile                          |
	+------+------------------------------------------------------+
	| copy | The name of profile from which the copy will be made |
	+------+------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/profiles/name/GLOBAL_copy/copy/GLOBAL HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:description:     The description of the new profile
:id:              An integral, unique identifier for the new profile
:idCopyFrom:      The integral, unique identifier for the profile from which the copy was made
:name:            The name of the new profile
:profileCopyFrom: The name of the profile from which the copy was made

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Fri, 07 Dec 2018 22:03:54 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Sat, 08 Dec 2018 02:03:54 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: r6V9viEZui1WCns0AUGEx1MtxjjXiU8SZVOtSQjeq7ZJDLl5s8fMmjJdR/HRWduHn7Ax6GzYhoKwnIjMyc7ZWg==
	Content-Length: 252

	{ "alerts": [
		{
			"level": "success",
			"text": "Created new profile [ GLOBAL_copy ] from existing profile [ GLOBAL ]"
		}
	],
	"response": {
		"idCopyFrom": 1,
		"name": "GLOBAL_copy",
		"profileCopyFrom": "GLOBAL",
		"id": 17,
		"description": "Global Traffic Ops profile, DO NOT DELETE"
	}}
