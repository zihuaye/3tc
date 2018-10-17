# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

"""
This module contains miscellaneous utilities, typically dealing with string
manipulation or user input/output
"""
import logging
from sys import stderr
import requests

def getYesNoResponse(prmpt:str, default:str = None) -> bool:
	"""
	Utility function to get an interactive yes/no response to the prompt `prmpt`

	:param prmpt: The prompt to display to users
	:param default: The default response; should be one of ``'y'``, ``"yes"``, ``'n'`` or ``"no"``
		(case insensitive)
	:raises AttributeError: if 'prmpt' and/or 'default' is/are not strings
	:returns: the parsed response as a boolean
	"""
	if default:
		prmpt = prmpt.rstrip().rstrip(':') + '['+default+"]:"
	while True:
		choice = input(prmpt).lower()

		if choice in {'y', 'yes'}:
			return True
		if choice in {'n', 'no'}:
			return False
		if not choice and default is not None:
			return default.lower() in {'y', 'yes'}

		print("Please enter a yes/no response.", file=stderr)

def getTextResponse(uri:str, cookies:dict = None, verify:bool = True) -> str:
	"""
	Gets the plaintext response body of an HTTP ``GET`` request

	:param uri: The full path to a resource for the request
	:param cookies: An optional dictionary of cookie names mapped to values
	:param verify: If :const:`True`, the SSL keys used to communicate with the full URI will be
		verified

	:raises ConnectionError: when an error occurs trying to communicate with the server
	:raises ValueError: if the server's response cannot be interpreted as a UTF-8 string - e.g.
		when the response body is raw binary data but the response headers claim it's UTF-16
	"""
	logging.info("Getting plaintext response via 'HTTP GET %s'", uri)

	response = requests.get(uri, cookies=cookies, verify=verify)

	if response.status_code not in range(200, 300):
		logging.warning("Status code (%d) seems to indicate failure!", response.status_code)
		logging.debug("Response: %r\n%r", response.headers, response.content)

	return response.text

def getJSONResponse(uri:str, cookies:dict = None, verify:bool = True) -> dict:
	"""
	Retrieves a JSON object from some HTTP API

	:param uri: The URI to fetch
	:param cookies: A dictionary of cookie names mapped to values
	:param verify: If this is :const:`True`, the SSL keys will be verified during handshakes with
		'https' URIs
	:returns: The decoded JSON object
	:raises ConnectionError: when an error occurs trying to communicate with the server
	:raises ValueError: when the request completes successfully, but the response body
		does not represent a JSON-encoded object.
	"""

	logging.info("Getting JSON response via 'HTTP GET %s", uri)

	try:
		response = requests.get(uri, cookies=cookies, verify=verify)
	except (ValueError, ConnectionError, requests.exceptions.RequestException) as e:
		raise ConnectionError from e

	if response.status_code not in range(200, 300):
		logging.warning("Status code (%d) seems to indicate failure!", response.status_code)
		logging.debug("Response: %r\n%r", response.headers, response.content)

	return response.json()
