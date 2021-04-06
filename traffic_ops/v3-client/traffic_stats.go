package client

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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

// GetCurrentStats gets current stats for each CDNs and a total across them
func (to *Session) GetCurrentStats() (tc.TrafficStatsCDNStatsResponse, toclientlib.ReqInf, error) {
	resp := tc.TrafficStatsCDNStatsResponse{}
	reqInf, err := to.get("/current_stats", nil, &resp)
	return resp, reqInf, err
}
