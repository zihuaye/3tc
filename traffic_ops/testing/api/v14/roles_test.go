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
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	roleGood         = 0
	roleInvalidCap   = 1
	roleNeedCap      = 2
	roleBadPrivLevel = 3
)

func TestRoles(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Roles}, func() {
		UpdateTestRoles(t)
		GetTestRoles(t)
	})
}

func CreateTestRoles(t *testing.T) {
	expectedAlerts := []tc.Alerts{tc.Alerts{[]tc.Alert{tc.Alert{"role was created.", "success"}}}, tc.Alerts{[]tc.Alert{tc.Alert{"can not add non-existent capabilities: [invalid-capability]", "error"}}}}
	for i, role := range testData.Roles {
		var alerts tc.Alerts
		alerts, _, status, err := TOSession.CreateRole(role)
		log.Debugln("Status Code: ", status)
		log.Debugln("Response: ", alerts)
		if err != nil {
			log.Debugf("error: %v", err)
			//t.Errorf("could not CREATE role: %v\n", err)
		}
		if !reflect.DeepEqual(alerts, expectedAlerts[i]) {
			t.Errorf("got alerts: %v but expected alerts: %v", alerts, expectedAlerts[i])
		}
	}
}

func UpdateTestRoles(t *testing.T) {
	log.Debugf("testData.Roles contains: %++v\n", testData.Roles)
	firstRole := testData.Roles[0]
	// Retrieve the Role by role so we can get the id for the Update
	resp, _, status, err := TOSession.GetRoleByName(*firstRole.Name)
	log.Debugln("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot GET Role by role: %v - %v\n", firstRole.Name, err)
	}
	log.Debugf("got response: %++v\n", resp)
	remoteRole := resp[0]
	expectedRole := "new admin2"
	remoteRole.Name = &expectedRole
	var alert tc.Alerts
	alert, _, status, err = TOSession.UpdateRoleByID(*remoteRole.ID, remoteRole)
	log.Debugln("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot UPDATE Role by id: %v - %v\n", err, alert)
	}

	// Retrieve the Role to check role got updated
	resp, _, status, err = TOSession.GetRoleByID(*remoteRole.ID)
	log.Debugln("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot GET Role by role: %v - %v\n", firstRole.Name, err)
	}
	respRole := resp[0]
	if *respRole.Name != expectedRole {
		t.Errorf("results do not match actual: %s, expected: %s\n", *respRole.Name, expectedRole)
	}

	// Set the name back to the fixture value so we can delete it after
	remoteRole.Name = firstRole.Name
	alert, _, status, err = TOSession.UpdateRoleByID(*remoteRole.ID, remoteRole)
	log.Debugln("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot UPDATE Role by id: %v - %v\n", err, alert)
	}

}

func GetTestRoles(t *testing.T) {
	role := testData.Roles[roleGood]
	resp, _, status, err := TOSession.GetRoleByName(*role.Name)
	log.Debugln("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot GET Role by role: %v - %v\n", err, resp)
	}

}

func DeleteTestRoles(t *testing.T) {

	role := testData.Roles[roleGood]
	// Retrieve the Role by name so we can get the id
	resp, _, status, err := TOSession.GetRoleByName(*role.Name)
	log.Debugln("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot GET Role by name: %v - %v\n", role.Name, err)
	}
	respRole := resp[0]

	delResp, _, status, err := TOSession.DeleteRoleByID(*respRole.ID)
	log.Debugln("Status Code: ", status)

	if err != nil {
		t.Errorf("cannot DELETE Role by role: %v - %v\n", err, delResp)
	}

	// Retrieve the Role to see if it got deleted
	roleResp, _, status, err := TOSession.GetRoleByName(*role.Name)
	log.Debugln("Status Code: ", status)

	if err != nil {
		t.Errorf("error deleting Role role: %s\n", err.Error())
	}
	if len(roleResp) > 0 {
		t.Errorf("expected Role : %s to be deleted\n", *role.Name)
	}

}
