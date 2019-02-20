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
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestSteeringTargets(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices, SteeringTargets}, func() {
		GetTestSteeringTargets(t)
		UpdateTestSteeringTargets(t)
	})
}

func CreateTestSteeringTargets(t *testing.T) {
	log.Debugln("CreateTestSteeringTargets")
	for _, st := range testData.SteeringTargets {
		if st.Type == nil {
			t.Errorf("creating steering target: test data missing type\n")
		}
		if st.DeliveryService == nil {
			t.Errorf("creating steering target: test data missing ds\n")
		}
		if st.Target == nil {
			t.Errorf("creating steering target: test data missing target\n")
		}

		{
			respTypes, _, err := TOSession.GetTypeByName(*st.Type)
			if err != nil {
				t.Errorf("creating steering target: getting type: %v\n", err)
			} else if len(respTypes) < 1 {
				t.Errorf("creating steering target: getting type: not found\n")
			}
			st.TypeID = util.IntPtr(respTypes[0].ID)
		}
		{
			respDS, _, err := TOSession.GetDeliveryServiceByXMLID(string(*st.DeliveryService))
			if err != nil {
				t.Errorf("creating steering target: getting ds: %v\n", err)
			} else if len(respDS) < 1 {
				t.Errorf("creating steering target: getting ds: not found\n")
			}
			dsID := uint64(respDS[0].ID)
			st.DeliveryServiceID = &dsID
		}
		{
			respTarget, _, err := TOSession.GetDeliveryServiceByXMLID(string(*st.Target))
			if err != nil {
				t.Errorf("creating steering target: getting target ds: %v\n", err)
			} else if len(respTarget) < 1 {
				t.Errorf("creating steering target: getting target ds: not found\n")
			}
			targetID := uint64(respTarget[0].ID)
			st.TargetID = &targetID
		}

		resp, _, err := TOSession.CreateSteeringTarget(st)
		log.Debugln("Response: ", resp)
		if err != nil {
			t.Errorf("creating steering target: %v\n", err)
		}
	}
	log.Debugln("CreateTestSteeringTargets() PASSED")
}

func UpdateTestSteeringTargets(t *testing.T) {
	log.Debugln("UpdateTestSteeringTargets")

	if len(testData.SteeringTargets) < 1 {
		t.Errorf("updating steering target: no steering target test data\n")
	}
	st := testData.SteeringTargets[0]
	if st.DeliveryService == nil {
		t.Errorf("updating steering target: test data missing ds\n")
	}
	if st.Target == nil {
		t.Errorf("updating steering target: test data missing target\n")
	}

	respDS, _, err := TOSession.GetDeliveryServiceByXMLID(string(*st.DeliveryService))
	if err != nil {
		t.Errorf("updating steering target: getting ds: %v\n", err)
	}
	if len(respDS) < 1 {
		t.Errorf("updating steering target: getting ds: not found\n")
	}
	dsID := respDS[0].ID

	sts, _, err := TOSession.GetSteeringTargets(dsID)
	if err != nil {
		t.Errorf("updating steering targets: getting steering target: %v\n", err)
	}
	if len(sts) < 1 {
		t.Errorf("updating steering targets: getting steering target: got 0\n")
	}
	st = sts[0]

	expected := util.JSONIntStr(-12345)
	if st.Value != nil && *st.Value == expected {
		expected += 1
	}
	st.Value = &expected

	_, _, err = TOSession.UpdateSteeringTarget(st)
	if err != nil {
		t.Errorf("updating steering targets: updating: %+v\n", err)
	}

	sts, _, err = TOSession.GetSteeringTargets(dsID)
	if err != nil {
		t.Errorf("updating steering targets: getting updated steering target: %v\n", err)
	}
	if len(sts) < 1 {
		t.Errorf("updating steering targets: getting updated steering target: got 0\n")
	}
	actual := sts[0]

	if actual.DeliveryServiceID == nil {
		t.Errorf("steering target update: ds id expected %v actual %v\n", dsID, nil)
	} else if *actual.DeliveryServiceID != uint64(dsID) {
		t.Errorf("steering target update: ds id expected %v actual %v\n", dsID, *actual.DeliveryServiceID)
	}
	if actual.TargetID == nil {
		t.Errorf("steering target update: ds id expected %v actual %v\n", dsID, nil)
	} else if *actual.TargetID != *st.TargetID {
		t.Errorf("steering target update: ds id expected %v actual %v\n", *st.TargetID, *actual.TargetID)
	}
	if actual.TypeID == nil {
		t.Errorf("steering target update: ds id expected %v actual %v\n", *st.TypeID, nil)
	} else if *actual.TypeID != *st.TypeID {
		t.Errorf("steering target update: ds id expected %v actual %v\n", *st.TypeID, *actual.TypeID)
	}
	if actual.DeliveryService == nil {
		t.Errorf("steering target update: ds expected %v actual %v\n", *st.DeliveryService, nil)
	} else if *st.DeliveryService != *actual.DeliveryService {
		t.Errorf("steering target update: ds name expected %v actual %v\n", *st.DeliveryService, *actual.DeliveryService)
	}
	if actual.Target == nil {
		t.Errorf("steering target update: target expected %v actual %v\n", *st.Target, nil)
	} else if *st.Target != *actual.Target {
		t.Errorf("steering target update: target expected %v actual %v\n", *st.Target, *actual.Target)
	}
	if actual.Type == nil {
		t.Errorf("steering target update: type expected %v actual %v\n", *st.Type, nil)
	} else if *st.Type != *actual.Type {
		t.Errorf("steering target update: type expected %v actual %v\n", *st.Type, *actual.Type)
	}
	if actual.Value == nil {
		t.Errorf("steering target update: ds expected %v actual %v\n", *st.Value, nil)
	} else if *st.Value != *actual.Value {
		t.Errorf("steering target update: value expected %v actual %v\n", *st.Value, actual.Value)
	}
	log.Debugln("UpdateTestSteeringTargets() PASSED")
}

func GetTestSteeringTargets(t *testing.T) {
	log.Debugln("GetTestSteeringTargets")

	if len(testData.SteeringTargets) < 1 {
		t.Errorf("updating steering target: no steering target test data\n")
	}
	st := testData.SteeringTargets[0]
	if st.DeliveryService == nil {
		t.Errorf("updating steering target: test data missing ds\n")
	}

	respDS, _, err := TOSession.GetDeliveryServiceByXMLID(string(*st.DeliveryService))
	if err != nil {
		t.Errorf("creating steering target: getting ds: %v\n", err)
	} else if len(respDS) < 1 {
		t.Errorf("steering target get: getting ds: not found\n")
	}
	dsID := respDS[0].ID

	sts, _, err := TOSession.GetSteeringTargets(dsID)
	if err != nil {
		t.Errorf("steering target get: getting steering target: %v\n", err)
	}

	if len(sts) != len(testData.SteeringTargets) {
		t.Errorf("steering target get: expected %v actual %v\n", len(testData.SteeringTargets), len(sts))
	}

	expected := testData.SteeringTargets[0]
	actual := sts[0]

	if actual.DeliveryServiceID == nil {
		t.Errorf("steering target get: ds id expected %v actual %v\n", dsID, nil)
	} else if *actual.DeliveryServiceID != uint64(dsID) {
		t.Errorf("steering target get: ds id expected %v actual %v\n", dsID, *actual.DeliveryServiceID)
	}
	if actual.DeliveryService == nil {
		t.Errorf("steering target get: ds expected %v actual %v\n", expected.DeliveryService, nil)
	} else if *expected.DeliveryService != *actual.DeliveryService {
		t.Errorf("steering target get: ds name expected %v actual %v\n", expected.DeliveryService, actual.DeliveryService)
	}
	if actual.Target == nil {
		t.Errorf("steering target get: target expected %v actual %v\n", expected.Target, nil)
	} else if *expected.Target != *actual.Target {
		t.Errorf("steering target get: target expected %v actual %v\n", expected.Target, actual.Target)
	}
	if actual.Type == nil {
		t.Errorf("steering target get: type expected %v actual %v\n", expected.Type, nil)
	} else if *expected.Type != *actual.Type {
		t.Errorf("steering target get: type expected %v actual %v\n", expected.Type, actual.Type)
	}
	if actual.Value == nil {
		t.Errorf("steering target get: ds expected %v actual %v\n", expected.Value, nil)
	} else if *expected.Value != *actual.Value {
		t.Errorf("steering target get: value expected %v actual %v\n", *expected.Value, *actual.Value)
	}
	log.Debugln("GetTestSteeringTargets() PASSED")
}

func DeleteTestSteeringTargets(t *testing.T) {
	log.Debugln("DeleteTestSteeringTargets")
	dsIDs := []uint64{}
	for _, st := range testData.SteeringTargets {
		if st.DeliveryService == nil {
			t.Errorf("deleting steering target: test data missing ds\n")
		}
		if st.Target == nil {
			t.Errorf("deleting steering target: test data missing target\n")
		}

		respDS, _, err := TOSession.GetDeliveryServiceByXMLID(string(*st.DeliveryService))
		if err != nil {
			t.Errorf("deleting steering target: getting ds: %v\n", err)
		} else if len(respDS) < 1 {
			t.Errorf("deleting steering target: getting ds: not found\n")
		}
		dsID := uint64(respDS[0].ID)
		st.DeliveryServiceID = &dsID

		dsIDs = append(dsIDs, dsID)

		respTarget, _, err := TOSession.GetDeliveryServiceByXMLID(string(*st.Target))
		if err != nil {
			t.Errorf("deleting steering target: getting target ds: %v\n", err)
		} else if len(respTarget) < 1 {
			t.Errorf("deleting steering target: getting target ds: not found\n")
		}
		targetID := uint64(respTarget[0].ID)
		st.TargetID = &targetID

		_, _, err = TOSession.DeleteSteeringTarget(int(*st.DeliveryServiceID), int(*st.TargetID))
		if err != nil {
			t.Errorf("deleting steering target: deleting: %+v\n", err)
		}
	}

	for _, dsID := range dsIDs {
		sts, _, err := TOSession.GetSteeringTargets(int(dsID))
		if err != nil {
			t.Errorf("deleting steering targets: getting steering target: %v\n", err)
		}
		if len(sts) != 0 {
			t.Errorf("deleting steering targets: after delete, getting steering target: expected 0 actual %+v\n", len(sts))
		}
	}
	log.Debugln("DeleteTestSteeringTargets() PASSED")
}
