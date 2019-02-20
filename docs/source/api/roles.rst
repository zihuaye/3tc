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

.. _to-api-roles:

*********
``roles``
*********

``GET``
=======
Retrieves all user roles.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:description: A description of the role
:id:          The integral, unique identifier for this role
:name:        The name of the role
:privLevel:   An integer that allows for comparison between roles

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: caagePqpL6u9Mn1UBIDCJSgiLAKOHm72/DcrkxCuS7oLekMe87BkGhyJzkhQqUJh/CTmokr9x053GQ5FjhSKhg==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 15:34:05 GMT
	Transfer-Encoding: chunked

	{ "response": [
		{
			"id": 4,
			"name": "disallowed",
			"description": "Block all access",
			"privLevel": 0,
			"capabilities": []
		}
	]}

.. note:: The response example for this method of this endpoint has been truncated to only the last element of the resultant array, as the full response was hundreds of lines long.
