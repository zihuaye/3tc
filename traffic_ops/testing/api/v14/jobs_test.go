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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestJobs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		CreateTestJobs(t)
		GetTestJobs(t)
	})
}

func CreateTestJobs(t *testing.T) {
	toDSes, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v - %v\n", err, toDSes)
	}
	dsNameIDs := map[string]int64{}
	for _, ds := range toDSes {
		dsNameIDs[ds.XMLID] = int64(ds.ID)
	}

	for i, job := range testData.Jobs {
		job.Request.StartTime = time.Now().UTC()
		job.Request.DeliveryServiceID = dsNameIDs[job.DSName]
		testData.Jobs[i] = job
	}

	for _, job := range testData.Jobs {
		id, ok := dsNameIDs[job.DSName]
		if !ok {
			t.Fatalf("can't create test data job: delivery service '%v' not found in Traffic Ops", job.DSName)
		}
		job.Request.DeliveryServiceID = id
		_, _, err := TOSession.CreateJob(job.Request)
		if err != nil {
			t.Errorf("could not CREATE job: %v\n", err)
		}
	}
}

func GetTestJobs(t *testing.T) {
	toJobs, _, err := TOSession.GetJobs(nil, nil)
	if err != nil {
		t.Fatalf("error getting jobs: " + err.Error())
	}

	toDSes, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v - %v\n", err, toDSes)
	}

	dsIDNames := map[int64]string{}
	for _, ds := range toDSes {
		dsIDNames[int64(ds.ID)] = ds.XMLID
	}

	for _, testJob := range testData.Jobs {
		found := false
		for _, toJob := range toJobs {
			if toJob.DeliveryService != dsIDNames[testJob.Request.DeliveryServiceID] {
				continue
			}
			if !strings.HasSuffix(toJob.AssetURL, testJob.Request.Regex) {
				continue
			}
			toJobTime, err := time.Parse(tc.JobTimeFormat, toJob.StartTime)
			if err != nil {
				t.Errorf("job ds %v regex %v start time expected format '%+v' actual '%+v' error '%+v'", testJob.Request.DeliveryServiceID, testJob.Request.Regex, tc.JobTimeFormat, toJob.StartTime, err)
				continue
			}
			toJobTime = toJobTime.Round(time.Minute)
			testJobTime := testJob.Request.StartTime.Round(time.Minute)
			if !toJobTime.Equal(testJobTime) {
				t.Errorf("test job ds %v regex %v start time expected '%+v' actual '%+v'", testJob.Request.DeliveryServiceID, testJob.Request.Regex, testJobTime, toJobTime)
				continue
			}
			found = true
			break
		}
		if !found {
			t.Errorf("test job ds %v regex %v expected: exists, actual: not found", testJob.Request.DeliveryServiceID, testJob.Request.Regex)
		}
	}
}
