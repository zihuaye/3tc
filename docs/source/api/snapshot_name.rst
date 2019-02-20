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

.. _to-api-snapshot-name:

*********************
``snapshot/{{name}}``
*********************

``PUT``
=======
Performs a CDN snapshot. Effectively, this propagates the new *configuration* of the CDN to its *operating state*, which replaces the output of the :ref:`to-api-cdns-name-snapshot` endpoint with the output of the :ref:`to-api-cdns-name-snapshot-new` endpoint.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------+
	| Name | Description                                             |
	+======+=========================================================+
	| name | The name of the CDN for which a snapshot shall be taken |
	+------+---------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/snapshot/CDN-in-a-Box HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: z4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg/SpIdNs6c5H0NE8XYXysP+DGNKHfuwvY7kxvUdBeoGlODJ6+SfaPg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 12 Dec 2018 22:00:18 GMT
	Content-Length: 0
	Content-Type: text/plain; charset=utf-8
