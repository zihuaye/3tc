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

.. _to-api-cachegroups-id-parameters:

*********************************
``cachegroups/{{ID}}/parameters``
*********************************
Gets all the parameters associated with a :term:`Cache Group`

.. seealso:: :ref:`param-prof`

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------------------------------------------------------+
	| Parameter | Description                                              |
	+===========+==========================================================+
	| ID        | The integral, unique identifier of a :term:`Cache Group` |
	+-----------+----------------------------------------------------------+


Response Structure
------------------
:configFile:  Configuration file associated with the parameter
:id:          A numeric, unique identifier for this parameter
:lastUpdated: The date and time at which this parameter was last updated, in an ISO-like format
:name:        Name of the parameter
:secure:      If ``true``, the parameter value is only visible to "admin"-role users
:value:       Value of the parameter

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 19:56:23 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Wed, 14 Nov 2018 23:56:23 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: DfqPtySzVMpnBYqVt/45sSRG/1pRTlQdIcYuQZ0CQt79QSHLzU5e4TbDqht6ntvNP041LimKsj5RzPlPX1n6tg==
	Content-Length: 135

	{ "response": [
		{
			"lastUpdated": "2018-11-14 18:22:43.754786+00",
			"value": "foobar",
			"secure": false,
			"name": "foo",
			"id": 124,
			"configFile": "bar"
		}
	]}
