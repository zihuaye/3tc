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
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

func TestOrigins(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Users, DeliveryServices, Coordinates, Origins}, func() {
		UpdateTestOrigins(t)
		GetTestOrigins(t)
		OriginTenancyTest(t)
	})
}

func CreateTestOrigins(t *testing.T) {
	// loop through origins, assign FKs and create
	for _, origin := range testData.Origins {
		_, _, err := TOSession.CreateOrigin(origin)
		if err != nil {
			t.Errorf("could not CREATE origins: %v\n", err)
		}
	}
}

func GetTestOrigins(t *testing.T) {
	_, _, err := TOSession.GetOrigins()
	if err != nil {
		t.Errorf("cannot GET origins: %v\n", err)
	}

	for _, origin := range testData.Origins {
		resp, _, err := TOSession.GetOriginByName(*origin.Name)
		if err != nil {
			t.Errorf("cannot GET Origin by name: %v - %v\n", err, resp)
		}
	}
}

func UpdateTestOrigins(t *testing.T) {
	firstOrigin := testData.Origins[0]
	// Retrieve the origin by name so we can get the id for the Update
	resp, _, err := TOSession.GetOriginByName(*firstOrigin.Name)
	if err != nil {
		t.Errorf("cannot GET origin by name: %v - %v\n", *firstOrigin.Name, err)
	}
	remoteOrigin := resp[0]
	updatedPort := 4321
	updatedFQDN := "updated.example.com"

	// update port and FQDN values on origin
	remoteOrigin.Port = &updatedPort
	remoteOrigin.FQDN = &updatedFQDN
	updResp, _, err := TOSession.UpdateOriginByID(*remoteOrigin.ID, remoteOrigin)
	if err != nil {
		t.Errorf("cannot UPDATE Origin by name: %v - %v\n", err, updResp.Alerts)
	}

	// Retrieve the origin to check port and FQDN values were updated
	resp, _, err = TOSession.GetOriginByID(*remoteOrigin.ID)
	if err != nil {
		t.Errorf("cannot GET Origin by ID: %v - %v\n", *remoteOrigin.Name, err)
	}

	respOrigin := resp[0]
	if *respOrigin.Port != updatedPort {
		t.Errorf("results do not match actual: %d, expected: %d\n", *respOrigin.Port, updatedPort)
	}
	if *respOrigin.FQDN != updatedFQDN {
		t.Errorf("results do not match actual: %s, expected: %s\n", *respOrigin.FQDN, updatedFQDN)
	}
}

func OriginTenancyTest(t *testing.T) {
	origins, _, err := TOSession.GetOrigins()
	if err != nil {
		t.Errorf("cannot GET origins: %v\n", err)
	}
	tenant3Origin := tc.Origin{}
	foundTenant3Origin := false
	for _, o := range origins {
		if *o.FQDN == "origin.ds3.example.net" {
			tenant3Origin = o
			foundTenant3Origin = true
		}
	}
	if !foundTenant3Origin {
		t.Error("expected to find origin with tenant 'tenant3' and fqdn 'origin.ds3.example.net'")
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	tenant4TOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "tenant4user", "pa$$word", true, "to-api-v14-client-tests/tenant4user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to log in with tenant4user: %v", err.Error())
	}

	originsReadableByTenant4, _, err := tenant4TOClient.GetOrigins()
	if err != nil {
		t.Error("tenant4user cannot GET origins")
	}

	// assert that tenant4user cannot read origins outside of its tenant
	for _, origin := range originsReadableByTenant4 {
		if *origin.FQDN == "origin.ds3.example.net" {
			t.Error("expected tenant4 to be unable to read origins from tenant 3")
		}
	}

	// assert that tenant4user cannot update tenant3user's origin
	if _, _, err = tenant4TOClient.UpdateOriginByID(*tenant3Origin.ID, tenant3Origin); err == nil {
		t.Error("expected tenant4user to be unable to update tenant3's origin")
	}

	// assert that tenant4user cannot delete an origin outside of its tenant
	if _, _, err = tenant4TOClient.DeleteOriginByID(*origins[0].ID); err == nil {
		t.Errorf("expected tenant4user to be unable to delete an origin outside of its tenant (origin %s)", *origins[0].Name)
	}

	// assert that tenant4user cannot create origins outside of its tenant
	tenant3Origin.FQDN = util.StrPtr("origin.tenancy.test.example.com")
	if _, _, err = tenant4TOClient.CreateOrigin(tenant3Origin); err == nil {
		t.Errorf("expected tenant4user to be unable to create an origin outside of its tenant")
	}
}

func DeleteTestOrigins(t *testing.T) {
	for _, origin := range testData.Origins {
		resp, _, err := TOSession.GetOriginByName(*origin.Name)
		if err != nil {
			t.Errorf("cannot GET Origin by name: %v - %v\n", *origin.Name, err)
		}
		if len(resp) > 0 {
			respOrigin := resp[0]

			delResp, _, err := TOSession.DeleteOriginByID(*respOrigin.ID)
			if err != nil {
				t.Errorf("cannot DELETE Origin by ID: %v - %v\n", err, delResp)
			}

			// Retrieve the Origin to see if it got deleted
			org, _, err := TOSession.GetOriginByName(*origin.Name)
			if err != nil {
				t.Errorf("error deleting Origin name: %s\n", err.Error())
			}
			if len(org) > 0 {
				t.Errorf("expected Origin name: %s to be deleted\n", *origin.Name)
			}
		}
	}
}
