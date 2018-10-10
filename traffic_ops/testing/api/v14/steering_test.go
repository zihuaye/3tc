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
)

func TestSteering(t *testing.T) {
	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestProfiles(t)
	CreateTestStatuses(t)
	CreateTestDivisions(t)
	CreateTestRegions(t)
	CreateTestPhysLocations(t)
	CreateTestCacheGroups(t)
	CreateTestServers(t)
	CreateTestDeliveryServices(t)
	CreateTestSteeringTargets(t)

	GetTestSteering(t)

	DeleteTestSteeringTargets(t)
	DeleteTestDeliveryServices(t)
	DeleteTestServers(t)
	DeleteTestCacheGroups(t)
	DeleteTestPhysLocations(t)
	DeleteTestRegions(t)
	DeleteTestDivisions(t)
	DeleteTestStatuses(t)
	DeleteTestProfiles(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)
}

func GetTestSteering(t *testing.T) {
	log.Debugln("GetTestSteering")

	if len(testData.SteeringTargets) < 1 {
		t.Fatalf("get steering: no steering target test data\n")
	}
	st := testData.SteeringTargets[0]
	if st.DeliveryService == nil {
		t.Fatalf("get steering: test data missing ds\n")
	}

	steerings, _, err := TOSession.Steering()
	if err != nil {
		t.Fatalf("steering get: getting steering: %v\n", err)
	}

	if len(steerings) != len(testData.SteeringTargets) {
		t.Fatalf("steering get: expected %v actual %v\n", len(testData.SteeringTargets), len(steerings))
	}

	if steerings[0].ClientSteering {
		t.Fatalf("steering get: ClientSteering expected %v actual %v\n", false, true)
	}
	if len(steerings[0].Targets) != 1 {
		t.Fatalf("steering get: Targets expected %v actual %v\n", 1, len(steerings[0].Targets))
	}
	if steerings[0].Targets[0].Order != 0 {
		t.Fatalf("steering get: Targets Order expected %v actual %v\n", 0, steerings[0].Targets[0].Order)
	}
	if testData.SteeringTargets[0].Value != nil && steerings[0].Targets[0].Weight != int32(*testData.SteeringTargets[0].Value) {
		t.Fatalf("steering get: Targets Order expected %v actual %v\n", testData.SteeringTargets[0].Value, steerings[0].Targets[0].Weight)
	}
	if steerings[0].Targets[0].GeoOrder != nil {
		t.Fatalf("steering get: Targets Order expected %v actual %+v\n", nil, *steerings[0].Targets[0].GeoOrder)
	}
	if steerings[0].Targets[0].Longitude != nil {
		t.Fatalf("steering get: Targets Order expected %v actual %+v\n", nil, *steerings[0].Targets[0].Longitude)
	}
	if steerings[0].Targets[0].Latitude != nil {
		t.Fatalf("steering get: Targets Order expected %v actual %+v\n", nil, *steerings[0].Targets[0].Latitude)
	}
	log.Debugln("GetTestSteering() PASSED")
}
