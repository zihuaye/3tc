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

.. _to-api-cachegroups-id-queue_update:

***********************************
``cachegroups/{{ID}}/queue_update``
***********************************

``POST``
========
Queue or dequeue updates for all servers assigned to a :term:`Cache Group` limited to a specific CDN.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------------------------------+
	| Name | Description                                                                                             |
	+======+=========================================================================================================+
	| ID   | The integral, unique identifier for the :term:`Cache Group` for which updates are being queued/dequeued |
	+------+---------------------------------------------------------------------------------------------------------+

:action: The action to perform; one of "queue" or "dequeue"
:cdn:    The full name of the CDN in need of update queue/dequeue\ [1]_
:cdnId:  The integral, unique identifier for the CDN in need of update queue/dequeue\ [1]_

.. code-block:: http
	:caption: Request Example

	POST /api/1.3/cachegroups/8/queue_update HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 42
	Content-Type: application/json

	{"action": "queue", "cdn": "CDN-in-a-Box"}

.. [1] Either 'cdn' or 'cdnID' *must* be in the request data (but not both).

Response Structure
------------------
:action:         The action processed, one of "queue" or "dequeue"
:cachegroupId:   The integral, unique identifier of the :term:`Cache Group` for which updates were queued/dequeued
:cachegroupName: The name of the :term:`Cache Group` for which updates were queued/dequeued
:cdn:            The name of the CDN to which the queue/dequeue operation was restricted
:serverNames:    An array of the (short) hostnames of the servers within the :term:`Cache Group` which are also assigned to the CDN specified in the ``"cdn"`` field

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: UAcP7LrflU1RnfR4UqbQrJczlk5rkrcLOtTXJTFvIUXxK1EklZkHkE4vewjDaVIhJJ6YQg8jmPGQpr+x1RHabw==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 20:19:46 GMT
	Content-Length: 115

	{ "response": {
		"cachegroupName": "test",
		"action": "queue",
		"serverNames": [
			"foo"
		],
		"cdn": "CDN-in-a-Box",
		"cachegroupID": 8
	}}
