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

.. _to-api-phys_locations-trimmed:

**************************
``phys_locations/trimmed``
**************************

``GET``
=======
Retrieves only the names of physical locations

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available

Response Structure
------------------
:name: The name of the physical location

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: s4/q6oyQHa+mQ3d3gRGHvVsRyvsrkKxYnP574rVVUji0hHxYDbOnyPPswi4MuuQRm7dZq8cp4/iw9rlLRkBU0g==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 22:35:02 GMT
	Content-Length: 78

	{ "response": [
		{
			"name": "CDN_in_a_Box"
		},
		{
			"name": "Apachecon North America 2018"
		}
	]}
