/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

var pd = require('./pageData.js');
var cfunc = require('../common/commonFunctions.js');

describe('Traffic Portal CDNs Test Suite', function() {
	var pageData = new pd();
	var commonFunctions = new cfunc();
	var myNewCDN = 'pTestCDN';
	var myDomainName = 'ptest.com';
	var mydnssec = 'true';

	it('should go to the CDNs page', function() {
		console.log("Go to the CDNs page");
		browser.get(browser.baseUrl + "/#!/cdns");
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/cdns");
	});

	it('should open new CDN form page', function() {
		console.log("Open new CDN form page");
		browser.driver.findElement(by.name('createCdnButton')).click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/cdns/new");
	});

	it('should fill out form, create button is enabled and submit', function () {
		console.log("Filling out form, check create button is enabled and submit");
		expect(pageData.createButton.isEnabled()).toBe(false);
		pageData.dnssecEnabled.click();
		pageData.dnssecEnabled.sendKeys(mydnssec);
		pageData.name.sendKeys(myNewCDN);
		pageData.domainName.sendKeys(myDomainName);
		expect(pageData.createButton.isEnabled()).toBe(true);
		pageData.createButton.click();
		expect(browser.getCurrentUrl().then(commonFunctions.urlPath)).toEqual(commonFunctions.urlPath(browser.baseUrl)+"#!/cdns");
	});

	it('should verify the new CDN and then update CDN', function() {
		console.log("verifying the new CDN and then updating CDN");
		browser.sleep(250);
		pageData.searchFilter.sendKeys(myNewCDN);
		browser.sleep(250);
		element.all(by.repeater('cdn in ::cdns')).filter(function(row){
			return row.element(by.name('name')).getText().then(function(val){
				return val === myNewCDN;
			});
		}).get(0).click();
		browser.sleep(1000);
		pageData.domainName.clear();
		pageData.domainName.sendKeys('ptestUpdated.com');
		pageData.dnssecEnabled.click();
		pageData.dnssecEnabled.sendKeys('false');
		pageData.updateButton.click();
		expect(pageData.domainName.getText() === 'ptestUpdated.com');
	});

	it('should delete the new CDN', function() {
		console.log("Deleting " + myNewCDN);
		pageData.deleteButton.click();
		pageData.confirmWithNameInput.sendKeys(myNewCDN);
		pageData.deletePermanentlyButton.click();
	});
});
