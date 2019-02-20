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

.. _to-api-cachegroupparameters-id-parameterID:

***********************************************
``cachegroupparameters/{{ID}}/{{parameterID}}``
***********************************************

``DELETE``
==========
De-associate a parameter with a :term:`Cache Group`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------------+-------------------------------------------------------------------------------------------------+
	| Name        | Description                                                                                     |
	+=============+=================================================================================================+
	| ID          | Unique identifier for the :term:`Cache Group` which will have the parameter association deleted |
	+-------------+-------------------------------------------------------------------------------------------------+
	| parameterID | Unique identifier for the parameter which will be removed from a :term:`Cache Group`            |
	+-------------+-------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.1/cachegroupparameters/8/124 HTTP/1.1
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
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 18:26:40 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Wed, 14 Nov 2018 22:26:40 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: Cuj+ZPAKsDLp4FpbJDcwsWY0yVQAi1Um1CWraeTIQEMlyJSBEm17oKQWDjzTrvqqV8Prhu3gzlcHoVPzEpbQ1Q==
	Content-Length: 84

	{ "alerts": [
		{
			"level": "success",
			"text": "Profile parameter association was deleted."
		}
	]}
