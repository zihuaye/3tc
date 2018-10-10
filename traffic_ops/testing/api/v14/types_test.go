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

func TestTypes(t *testing.T) {

	CreateTestTypes(t)
	UpdateTestTypes(t)
	GetTestTypes(t)
	DeleteTestTypes(t)

}

func CreateTestTypes(t *testing.T) {
	log.Debugln("---- CreateTestTypes ----")

	for _, typ := range testData.Types {
		resp, _, err := TOSession.CreateType(typ)
		if err != nil {
			t.Errorf("could not CREATE types: %v\n", err)
		}
		log.Debugln("Response: ", resp)
	}

}

func UpdateTestTypes(t *testing.T) {
	log.Debugln("---- UpdateTestTypes ----")

	firstType := testData.Types[0]
	// Retrieve the Type by name so we can get the id for the Update
	resp, _, err := TOSession.GetTypeByName(firstType.Name)
	if err != nil {
		t.Errorf("cannot GET Type by name: %v - %v\n", firstType.Name, err)
	}
	remoteType := resp[0]
	expectedTypeName := "testType1"
	remoteType.Name = expectedTypeName
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateTypeByID(remoteType.ID, remoteType)
	if err != nil {
		t.Errorf("cannot UPDATE Type by id: %v - %v\n", err, alert)
	}

	// Retrieve the Type to check Type name got updated
	resp, _, err = TOSession.GetTypeByID(remoteType.ID)
	if err != nil {
		t.Errorf("cannot GET Type by name: %v - %v\n", firstType.Name, err)
	}
	respType := resp[0]
	if respType.Name != expectedTypeName {
		t.Errorf("results do not match actual: %s, expected: %s\n", respType.Name, expectedTypeName)
	}

	log.Debugln("Response Type: ", respType)

	respType.Name = firstType.Name
	alert, _, err = TOSession.UpdateTypeByID(respType.ID, respType)
	if err != nil {
		t.Errorf("cannot restore UPDATE Type by id: %v - %v\n", err, alert)
	}
}

func GetTestTypes(t *testing.T) {
	log.Debugln("---- GetTestTypes ----")

	for _, typ := range testData.Types {
		resp, _, err := TOSession.GetTypeByName(typ.Name)
		if err != nil {
			t.Errorf("cannot GET Type by name: %v - %v\n", err, resp)

		}

		log.Debugln("Response: ", resp)
	}
}

func DeleteTestTypes(t *testing.T) {
	log.Debugln("---- DeleteTestTypes ----")

	for _, typ := range testData.Types {
		// Retrieve the Type by name so we can get the id for the Update
		resp, _, err := TOSession.GetTypeByName(typ.Name)
		if err != nil || len(resp) == 0 {
			t.Errorf("cannot GET Type by name: %v - %v\n", typ.Name, err)
		}
		respType := resp[0]

		delResp, _, err := TOSession.DeleteTypeByID(respType.ID)
		if err != nil {
			t.Errorf("cannot DELETE Type by name: %v - %v\n", err, delResp)
		}

		// Retrieve the Type to see if it got deleted
		types, _, err := TOSession.GetTypeByName(typ.Name)
		if err != nil {
			t.Errorf("error deleting Type name: %s\n", err.Error())
		}
		if len(types) > 0 {
			t.Errorf("expected Type name: %s to be deleted\n", typ.Name)
		}
	}
}
