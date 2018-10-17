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

.. _ds_requests:

*************************
Delivery Service Requests
*************************
When enabled in ``traffic_portal_properties.json``, Delivery Service requests are created when *all* users attempt to create, update or delete a Delivery Service. This allows users with higher level permissions (ops or admin) to review the changes for completeness and accuracy before deploying the changes. In addition, most Delivery Service changes require cache configuration updates (aka queue updates) and/or a CDN snapshot. Both of these actions are reserved for users with elevated permissions.

A list of the Delivery Service requests associated with your tenant can be found under 'Services' -> 'Delivery Service Requests'

.. figure:: ../traffic_portal/images/tp_table_ds_requests.png
	:width: 60%
	:align: center
	:alt: A screenshot of the Traffic Portal UI depicting an example list of Delivery Service Requests

	Example Delivery Service Request Listing

Who Can Create a Delivery Service Request and How?
==================================================
Users with the Portal role (or above) can create Delivery Service Requests by doing one of three things:

- Creating a new Delivery Service
- Updating an existing Delivery Service
- Deleting an exiting Delivery Service

By performing one of these actions, a Delivery Service Request will be created for you with a status of 'draft' or 'submitted'. You determine the status of your request upon submission. Only change the status of your request to 'submitted' once the request is ready for review and deployment.

Who Can Fulfill a Delivery Service Request and How?
===================================================
Users with elevated permissions (Operations or above) can fulfill (apply the changes) or reject the Delivery Service Request. In fact, they can do all of the following:

Update the contents of the Delivery Service request
	This will update the "Last Edited By" field to indicate who last updated the request.

Assign or Unassign the Delivery Service Request
	Assignment is currently limited to current user. This is optional as fulfillment will auto-assign the request to the user doing the fulfillment.

Reject the Delivery Service Request
	Rejecting a Delivery Service Request will set status to 'rejected' and the request can no longer be modified. This will auto-assign the request to the user doing the rejection.

Fulfill the Delivery Service Request
	Fulfilling a Delivery Service Request will show the requested changes and, once committed, will apply the desired changes and set status to 'pending'. The request is pending because many types of changes will require cache configuration updates (aka queue updates) and/or a CDN snapshot. Once queue updates and/or CDN snapshot is complete, the request should be marked 'complete'.

Complete the Delivery Service Request
	Only after the Delivery Service Request has been fulfilled and the changes have been applied can a Delivery Service Request be marked as 'complete'. Marking a Delivery Service as 'complete' is currently a manual step because some changes require cache configuration updates (aka queue updates) and/or a CDN snapshot. Once that is done and the changes have been deployed, the request status should be changed from 'pending' to 'complete'.

Delete the Delivery Service request
	Delivery Service Requests with a status of 'draft' or 'submitted' can always be deleted entirely if appropriate.
