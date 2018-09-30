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

.. _dnssec-qht:

****************
Configure DNSSEC
****************

.. seealso:: :ref:`tr-dnssec`

.. Note:: In order for Traffic Ops to successfully store keys in Traffic Vault, at least one Riak Server needs to be configured in Traffic Ops. See the `Traffic Vault admin page <../traffic_vault.html>`_ for more information.

.. Note:: Currently DNSSEC is only supported for DNS delivery services.

#. Go to 'CDNs' and click on the desired CDN.

	.. figure:: dnssec/00.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the CDNs page

		CDNs Page

#. Click on 'Manage DNSSEC Keys' under the 'More' drop-down menu.

	.. figure:: dnssec/01.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the CDN details page

		CDN Details Page

#. Click on the 'Generate DNSSEC Keys' button.

	.. figure:: dnssec/02.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the CDN DNSSEC Key Management page

		DNSSEC Key Management Page

#. A modal will pop up asking you to confirm that you want to proceed.

	.. figure:: dnssec/03.png
		:width: 30%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the CDN DNSSEC Key Generation confirmation modal

		Confirmation Modal

#. Input the required information (reasonable defaults should be generated for you). When done, click on the green 'Generate' button.:

	.. note:: Depending upon the number of Delivery Services in the CDN, generating DNSSEC keys may take several seconds.

	.. figure:: dnssec/04.png
		:width: 50%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the CDN DNSSEC Key Generation page

		DNSSEC Key Generation Page

#. You will be prompted to confirm the changes by typing the name of the CDN into a text box. After doing so, click on the red 'Confirm' button.

	.. figure:: dnssec/05.png
		:width: 30%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the confirmation modal for committing changes to DNSSEC Keys.

		DNSSEC Key Change Confirmation


#. In order for DNSSEC to work properly, the DS Record information needs to be added to the parent zone of the CDN's domain (e.g. If the CDN's domain is 'ciab.cdn.local' the parent zone is 'cdn.local'). If you control your parent zone you can enter this information yourself, otherwise you will need to work with your DNS team to get the DS Record added to the parent zone.

#. Once DS Record information has been added to the parent zone, DNSSEC needs to be activated for the CDN so that Traffic Router will sign responses. Go back to the CDN details page for this CDN, and set the 'DNSSEC Enabled' field to 'true', then click the green 'Update' button.

	.. figure:: dnssec/06.png
		:width: 60%
		:align: center
		:alt: Screenshot of the Traffic Portal UI depicting the details page for a CDN when changing its 'DNSSEC Enabled' field

		Change 'DNSSEC Enabled' to 'true'

#. DNSSEC should now be active on your CDN and Traffic Router should be signing responses. This should be tested e.g. with this ``dig`` command: ``dig edge.cdn.local. +dnssec``.

#. When KSK expiration is approaching (default 365 days), it is necessary to manually generate a new KSK for the TLD (Top Level Domain) and add the DS Record to the parent zone. In order to avoid signing errors, it is suggested that an effective date is chosen which allows time for the DS Record to be added to the parent zone before the new KSK becomes active.
