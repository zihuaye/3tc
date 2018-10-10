package v14

/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestCDNs(t *testing.T) {

	CreateTestCDNs(t)
	UpdateTestCDNs(t)
	GetTestCDNs(t)
	DeleteTestCDNs(t)

}

func CreateTestCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.CreateCDN(cdn)
		log.Debugln("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE cdns: %v\n", err)
		}
	}

}

func UpdateTestCDNs(t *testing.T) {

	firstCDN := testData.CDNs[0]
	// Retrieve the CDN by name so we can get the id for the Update
	resp, _, err := TOSession.GetCDNByName(firstCDN.Name)
	if err != nil {
		t.Errorf("cannot GET CDN by name: '%s', %v\n", firstCDN.Name, err)
	}
	remoteCDN := resp[0]
	expectedCDNDomain := "domain2"
	remoteCDN.DomainName = expectedCDNDomain
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateCDNByID(remoteCDN.ID, remoteCDN)
	if err != nil {
		t.Errorf("cannot UPDATE CDN by id: %v - %v\n", err, alert)
	}

	// Retrieve the CDN to check CDN name got updated
	resp, _, err = TOSession.GetCDNByID(remoteCDN.ID)
	if err != nil {
		t.Errorf("cannot GET CDN by name: '$%s', %v\n", firstCDN.Name, err)
	}
	respCDN := resp[0]
	if respCDN.DomainName != expectedCDNDomain {
		t.Errorf("results do not match actual: %s, expected: %s\n", respCDN.DomainName, expectedCDNDomain)
	}

}

func GetTestCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.GetCDNByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot GET CDN by name: %v - %v\n", err, resp)
		}
	}
}

func DeleteTestCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		// Retrieve the CDN by name so we can get the id for the Update
		resp, _, err := TOSession.GetCDNByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot GET CDN by name: %v - %v\n", cdn.Name, err)
		}
		if len(resp) > 0 {
			respCDN := resp[0]

			_, _, err := TOSession.DeleteCDNByID(respCDN.ID)
			if err != nil {
				t.Errorf("cannot DELETE CDN by name: '%s' %v\n", respCDN.Name, err)
			}

			// Retrieve the CDN to see if it got deleted
			cdns, _, err := TOSession.GetCDNByName(cdn.Name)
			if err != nil {
				t.Errorf("error deleting CDN name: %s\n", err.Error())
			}
			if len(cdns) > 0 {
				t.Errorf("expected CDN name: %s to be deleted\n", cdn.Name)
			}
		}
	}
}
