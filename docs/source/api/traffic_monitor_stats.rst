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

.. _to-api-traffic_monitor-stats:

*************************
``traffic_monitor/stats``
*************************
.. deprecated:: TrafficControl 3.0.0
	This endpoint was used by the now-deprecated Traffic Ops UI, and will likely be removed in the future!

.. caution:: This page is a stub! Much of it may be missing or just downright wrong - it needs a lot of love from people with the domain knowledge required to update it.

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  **NOT PRESENT** - this endpoint returns a special, custom JSON response

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:aaData: An array of data points of some kind

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Mon, 03 Dec 2018 14:44:14 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Mon, 03 Dec 2018 18:44:14 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: yRHVMHW+Y78HgaU/UVcrcADq9Jw3ScP+IQEEVqy3R/0A757WM2ZpmGDECDkDp7crWckabMntHRIfaf/6hWJPoQ==
	Content-Length: 57

	{ "aaData": [
		[
			"0",
			"ALL",
			"ALL",
			"ALL",
			"true",
			"ALL",
			"0",
			"0"
		]
	]}
