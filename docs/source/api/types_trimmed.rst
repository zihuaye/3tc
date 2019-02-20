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

.. _to-api-types-trimmed:

*****************
``types/trimmed``
*****************
.. deprecated:: 1.1
	A type is not a large object. Just get the names from :ref:`to-api-types` instead.

``GET``
=======
Retrieves only the names of all of the types of things configured in Traffic Ops. Yes, that is as specific as a description of a 'type' can be.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available

Response Structure
------------------
:name: The name of the type

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Connection: keep-alive
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Set-Cookie: mojolicious=...; expires=Thu, 13 Dec 2018 03:37:01 GMT; path=/; HttpOnly
	Whole-Content-SHA512: Wh4z9VkNcOI8UzSTM77N+JFx5bP8yxRR4rg1fZIH40DI+0suOD36YhePUMMqMl6DIlIWjrnkj+iojuQ09oTzeg==
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Date: Wed, 12 Dec 2018 23:37:01 GMT
	Access-Control-Allow-Origin: *
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Content-Length: 1104
	Content-Type: application/json
	Server: Mojolicious (Perl)

	{ "response": [
		{
			"name": "AAAA_RECORD"
		},
		{
			"name": "ANY_MAP"
		}
	]}

.. note:: The response example for this endpoint has been truncated to only the first two elements of the resulting array, as the output was hundreds of lines long.
