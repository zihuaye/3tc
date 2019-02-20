<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->

# Traffic Ops API Tests

The Traffic Ops Client API tests are used to validate the clients responses against those from the Traffic Ops API.  

In order to run the tests you will need the following:

1. Port access to both the Postgres port (usually 5432) that your Traffic Ops instance is using as well as the Traffic Ops configured port (usually 443 or 60443).

2. An instance of Postgres running with a `to_test` database that has empty tables.

    To get your to_test database setup do the following:
    
    `$ cd trafficcontrol/traffic_ops/app`
    
    `$ db/admin --env=test reset` 

    NOTE on passwords:
    Check that the passwords defined defined for your `to_test` database match 
    here: `trafficcontrol/traffic_ops/app/conf/test/database.conf`
    and here: `traffic-ops-test.conf` 

    The Traffic Ops users will be created by the tool for accessing the API once the database is accessible.

    For more info see: http://trafficcontrol.apache.org/docs/latest/development/traffic_ops.html?highlight=reset

3. A running Traffic Ops instance running with the `secure` (https) and is pointing to the `to_test` 
   database by running in `MOJO_MODE=test` which will point to your `to_test` database.
    To get your to_test database setup do the following:
    
   	`$ export MOJO_MODE=test`  
   	
   	`$ cd trafficcontrol/traffic_ops/app`
   	
    `$ bin/start.pl --secure`

4. A running Traffic Ops Golang proxy pointing to the to_test database.
	`$ cd trafficcontrol/traffic_ops/traffic_ops_golang`
	`$ cp ../app/conf/cdn.conf $HOME/cdn.conf`
	change `traffic_ops_golang->port` to 8443

    `$ go build && ./traffic_ops_golang -cfg $HOME/cdn.conf -dbcfg ../app/conf/test/database.conf`

## Running the API Tests
The integration tests are run using `go test`, however, there are some flags that need to be provided in order for the tests to work.  

The flags are:

* usage - API Test tool usage
* cfg - the config file needed to run the tests
* env - Environment variables that can be used to override specific config options that are specified in the config file
* env_vars - Show environment variables that can be overridden
* test_data - traffic control
* run - Go runtime flag for executing a specific test case

Example command to run the tests: 
`TO_URL=https://localhost:8443 go test -v -cfg=traffic-ops-test.conf -run TestCDNs`



* It can take several minutes for the API tests to complete, so using the `-v` flag is recommended to see progress.*
