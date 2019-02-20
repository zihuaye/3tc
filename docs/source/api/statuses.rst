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

.. _to-api-statuses:

************
``statuses``
************

``GET``
=======
Retrieves a list of all server statuses.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------+----------+--------------------------------------------------------------+
	|    Name     | Required | Description                                                  |
	+=============+==========+==============================================================+
	| description |    no    | Return only statuses with this *exact* description           |
	+-------------+----------+--------------------------------------------------------------+
	|     id      |    no    | Return only the status with this integral, unique identifier |
	+-------------+----------+--------------------------------------------------------------+
	|    name     |    no    | Return only statuses with this name                          |
	+-------------+----------+--------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/statuses?name=REPORTED HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:description: A short description of the status
:id:          The integral, unique identifier of this status
:lastUpdated: The date and time at which this status was last modified, in ISO format
:name:        The name of the status

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: dHNip9kpTGGS1w39/fWcFehNktgmXZus8XaufnmDpv0PyG/3fK/KfoCO3ZOj9V74/CCffps7doEygWeL/xRtKA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 20:56:59 GMT
	Content-Length: 150

	{ "response": [
		{
			"description": "Server is online and reported in the health protocol.",
			"id": 3,
			"lastUpdated": "2018-12-10 19:11:17+00",
			"name": "REPORTED"
		}
	]}
