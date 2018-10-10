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
	"database/sql"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/config"
	_ "github.com/lib/pq"
)

var (
	Config   config.Config
	testData TrafficControl
)

func TestMain(m *testing.M) {
	var err error
	configFileName := flag.String("cfg", "traffic-ops-test.conf", "The config file path")
	tcFixturesFileName := flag.String("fixtures", "tc-fixtures.json", "The test fixtures for the API test tool")
	flag.Parse()

	if Config, err = config.LoadConfig(*configFileName); err != nil {
		fmt.Printf("Error Loading Config %v %v\n", Config, err)
		return
	}

	if err = log.InitCfg(Config); err != nil {
		fmt.Printf("Error initializing loggers: %v\n", err)
		return
	}

	log.Infof(`Using Config values:
			   TO Config File:       %s
			   TO Fixtures:          %s
			   TO URL:               %s
			   TO Session Timeout In Secs:  %d
			   DB Server:            %s
			   DB User:              %s
			   DB Name:              %s
			   DB Ssl:               %t`, *configFileName, *tcFixturesFileName, Config.TrafficOps.URL, Config.Default.Session.TimeoutInSecs, Config.TrafficOpsDB.Hostname, Config.TrafficOpsDB.User, Config.TrafficOpsDB.Name, Config.TrafficOpsDB.SSL)

	//Load the test data
	LoadFixtures(*tcFixturesFileName)

	var db *sql.DB
	db, err = OpenConnection()
	if err != nil {
		fmt.Printf("\nError opening connection to %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOpsDB.User, err)
		os.Exit(1)
	}
	defer db.Close()

	err = Teardown(db)
	if err != nil {
		fmt.Printf("\nError tearingdown data %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOpsDB.User, err)
		os.Exit(1)
	}

	err = SetupTestData(db)
	if err != nil {
		fmt.Printf("\nError setting up data %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOpsDB.User, err)
		os.Exit(1)
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	err = SetupSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
	if err != nil {
		fmt.Printf("\nError creating session to %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOpsDB.User, err)
		os.Exit(1)
	}

	// Now run the test case
	rc := m.Run()
	os.Exit(rc)

}
