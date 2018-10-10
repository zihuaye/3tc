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

package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func (to *Session) Federations() ([]tc.AllFederation, ReqInf, error) {
	type FederationResponse struct {
		Response []tc.AllFederation `json:"response"`
	}
	data := FederationResponse{}
	inf, err := get(to, apiBase+"/federations", &data)
	return data.Response, inf, err
}

func (to *Session) AllFederations() ([]tc.AllFederation, ReqInf, error) {
	type FederationResponse struct {
		Response []tc.AllFederation `json:"response"`
	}
	data := FederationResponse{}
	inf, err := get(to, apiBase+"/federations/all", &data)
	return data.Response, inf, err
}

func (to *Session) AllFederationsForCDN(cdnName string) ([]tc.AllFederation, ReqInf, error) {
	// because the Federations JSON array is heterogeneous (array members may be a AllFederation or AllFederationCDN), we have to try decoding each separately.
	type FederationResponse struct {
		Response []json.RawMessage `json:"response"`
	}
	data := FederationResponse{}
	inf, err := get(to, apiBase+"/federations/all?cdnName="+cdnName, &data)
	if err != nil {
		return nil, inf, err
	}

	feds := []tc.AllFederation{}
	for _, raw := range data.Response {
		fed := tc.AllFederation{}
		if err := json.Unmarshal([]byte(raw), &fed); err != nil {
			// we don't actually need the CDN, but we want to return an error if we got something unexpected
			cdnFed := tc.AllFederationCDN{}
			if err := json.Unmarshal([]byte(raw), &cdnFed); err != nil {
				return nil, inf, errors.New("Traffic Ops returned an unexpected object: '" + string(raw) + "'")
			}
		}
		feds = append(feds, fed)
	}
	return feds, inf, nil
}

func (to *Session) CreateFederationDeliveryServices(federationID int, deliveryServiceIDs []int, replace bool) (ReqInf, error) {
	req := tc.FederationDSPost{DSIDs: deliveryServiceIDs, Replace: &replace}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return ReqInf{CacheHitStatus: CacheHitStatusMiss}, err
	}
	resp := map[string]interface{}{}
	inf, err := makeReq(to, http.MethodPost, apiBase+`/federations/`+strconv.Itoa(federationID)+`/deliveryservices`, jsonReq, &resp)
	return inf, err
}
