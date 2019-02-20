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

.. _to-api-user-current:

****************
``user/current``
****************

``GET``
=======
.. deprecated:: 1.4
	As a username is needed to log in, any administrator or application must necessarily know the current username at any given time. Thus, use the ``username`` query parameter of a ``GET`` request to :ref:`to-api-users` instead.

Retrieves the profile for the authenticated user.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:addressLine1:     The user's address - including street name and number
:addressLine2:     An additional address field for e.g. apartment number
:city:             The name of the city wherein the user resides
:company:          The name of the company for which the user works
:country:          The name of the country wherein the user resides
:email:            The user's email address
:fullName:         The user's full name, e.g. "John Quincy Adams"
:gid:              A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX group ID of the user - now it is always ``null``
:id:               An integral, unique identifier for this user
:lastUpdated:      The date and time at which the user was last modified, in ISO format
:newUser:          A meta field with no apparent purpose that is usually ``null`` unless explicitly set during creation or modification of a user via some API endpoint
:phoneNumber:      The user's phone number
:postalCode:       The postal code of the area in which the user resides
:publicSshKey:     The user's public key used for the SSH protocol
:registrationSent: If the user was created using the :ref:`to-api-users-register` endpoint, this will be the date and time at which the registration email was sent - otherwise it will be ``null``
:role:             The integral, unique identifier of the highest-privilege role assigned to this user
:rolename:         The name of the highest-privilege role assigned to this user
:stateOrProvince:  The name of the state or province where this user resides
:tenant:           The name of the tenant to which this user belongs
:tenantId:         The integral, unique identifier of the tenant to which this user belongs
:uid:              A deprecated field only kept for legacy compatibility reasons that used to contain the UNIX user ID of the user - now it is always ``null``
:username:         The user's username

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: HQwu9FxFyinXSVFK5+wpEhSxU60KbqXuokFbMZ3OoerOoM5ZpWpglsHz7mRch8VAw0dzwsJzpPJivj07RiKaJg==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 13 Dec 2018 15:14:45 GMT
	Content-Length: 382

	{ "response": {
		"username": "admin",
		"localUser": true,
		"addressLine1": null,
		"addressLine2": null,
		"city": null,
		"company": null,
		"country": null,
		"email": null,
		"fullName": null,
		"gid": null,
		"id": 2,
		"newUser": false,
		"phoneNumber": null,
		"postalCode": null,
		"publicSshKey": null,
		"role": 1,
		"rolename": "admin",
		"stateOrProvince": null,
		"tenant": "root",
		"tenantId": 1,
		"uid": null,
		"lastUpdated": "2018-12-12 16:26:32+00"
	}}

``PUT``
=======
.. deprecated:: 1.4
	Use the ``PUT`` method of the :ref:`to-api-users` instead.

.. warning:: Users that login via LDAP pass-back cannot be modified

Updates the date for the authenticated user.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  ``undefined``

Request Structure
-----------------
:addressLine1:       An optional field which should contain the user's address - including street name and number
:addressLine2:       An optional field which should contain an additional address field for e.g. apartment number
:city:               An optional field which should contain the name of the city wherein the user resides
:company:            An optional field which should contain the name of the company for which the user works
:confirmLocalPasswd: The 'confirm' field in a new user's password specification - must match ``localPasswd``
:country:            An optional field which should contain the name of the country wherein the user resides
:email:              The user's email address

	.. versionchanged:: 1.4
		Prior to version 1.4, the email was validated using the `Email::Valid Perl package <https://metacpan.org/pod/Email::Valid>`_ but is now validated (circuitously) by `GitHub user asaskevich's regular expression <https://github.com/asaskevich/govalidator/blob/9a090521c4893a35ca9a228628abf8ba93f63108/patterns.go#L7>`_ . Note that neither method can actually distinguish a valid, deliverable, email address but merely ensure the email is in a commonly-found format.

:fullName:        The user's full name, e.g. "John Quincy Adams"
:localPasswd:     The user's password
:newUser:         An optional meta field with no apparent purpose - don't use this
:phoneNumber:     An optional field which should contain the user's phone number
:postalCode:      An optional field which should contain the user's postal code
:publicSshKey:    An optional field which should contain the user's public encryption key used for the SSH protocol
:role:            The number that corresponds to the highest permission role which will be permitted to the user
:stateOrProvince: An optional field which should contain the name of the state or province in which the user resides
:tenantId:        The integral, unique identifier of the tenant to which the new user shall belong

	.. note:: This field is optional if and only if tenancy is not enabled in Traffic Control

:username: The user's new username

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/user/current HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 483
	Content-Type: application/json

	{ "user": {
		"addressLine1": "not a real address",
		"addressLine2": "not a real address either",
		"city": "not a real city",
		"company": "not a real company",
		"country": "not a real country",
		"email": "not@real.email",
		"fullName": "Not a real fullName",
		"phoneNumber": "not a real phone number",
		"postalCode": "not a real postal code",
		"publicSshKey": "not a real ssh key",
		"stateOrProvince": "not a real state or province",
		"tenantId": 1,
		"role": 1,
		"username": "admin"
	}}

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
	Date: Thu, 13 Dec 2018 21:05:49 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Fri, 14 Dec 2018 01:05:49 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: sHFqZQ4Cv7IIWaIejoAvM2Fr/HSupcX3D16KU/etjw+4jcK9EME3Bq5ohLC+eQ52BDCKW2Ra+AC3TfFtworJww==
	Content-Length: 79

	{ "alerts": [
		{
			"level": "success",
			"text": "User profile was successfully updated"
		}
	]}
