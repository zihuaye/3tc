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

.. _to-api-servers-id-queue_update:

*******************************
``servers/{{ID}}/queue_update``
*******************************
.. deprecated:: 1.1
	Use the ``PUT`` method of the :ref:`to-api-servers-id` endpoint instead.

``POST``
========
Queue or dequeue updates for a specific server.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------------------+
	| Name | Description                                                                                 |
	+======+=============================================================================================+
	|  ID  | The integral, unique identifier of the server on which updates are being queued or dequeued |
	+------+---------------------------------------------------------------------------------------------+

:action: A string describing what action to take regarding server updates; one of:

	queue
		Enqueue updates for the server, propagating configuration changes to the actual server
	dequeue
		Cancels any pending updates on the server

.. code-block:: http
	:caption: Request Example

	POST /api/1.1/servers/13/queue_update HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 22
	Content-Type: application/json

	{
		"action": "dequeue"
	}

Response Structure
------------------
:action: The action processed, one of:

	queue
		Enqueued updates for the server, propagating configuration changes to the actual server
	dequeue
		Canceled any pending updates on the server

:serverId: The integral, unique identifier of the server on which ``action`` was taken

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Mon, 10 Dec 2018 18:20:04 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Mon, 10 Dec 2018 22:20:04 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: 9Mmo9hIFZyF5gAvfdJD//VH9eNgiHVLinXt88H0GlJSHhwND8gMxaFyC+f9XZfiNAoGd1MKi1934ZJGmaIR6qQ==
	Content-Length: 49

	{
		"response": {
			"serverId": "13",
			"action": "dequeue"
		}
	}
