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
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_v2_ASNs = "/api/1.3/asns"
)

// Create a ASN
func (to *Session) CreateASN(entity tc.ASN) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(entity)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_v2_ASNs, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a ASN by ID
func (to *Session) UpdateASNByID(id int, entity tc.ASN) (tc.Alerts, ReqInf, error) {

	var remoteAddr net.Addr
	reqBody, err := json.Marshal(entity)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	route := fmt.Sprintf("%s/%d", API_v2_ASNs, id)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Returns a list of ASNs
func (to *Session) GetASNs() ([]tc.ASN, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v2_ASNs, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ASNsResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GET a ASN by the id
func (to *Session) GetASNByID(id int) ([]tc.ASN, ReqInf, error) {
	route := fmt.Sprintf("%s/%d", API_v2_ASNs, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ASNsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET an ASN by the asn number
func (to *Session) GetASNByASN(asn int) ([]tc.ASN, ReqInf, error) {
	url := fmt.Sprintf("%s?asn=%d", API_v2_ASNs, asn)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ASNsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DELETE an ASN by asn number
func (to *Session) DeleteASNByASN(asn int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s/asn/%d", API_v2_ASNs, asn)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

func (to *Session) DeleteDeliveryServiceServer(dsID int, serverID int) (tc.Alerts, ReqInf, error) {
	route := apiBase + `/deliveryservice_server/` + strconv.Itoa(dsID) + "/" + strconv.Itoa(serverID)
	reqResp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, errors.New("requesting from Traffic Ops: " + err.Error())
	}
	defer reqResp.Body.Close()
	resp := tc.Alerts{}
	if err = json.NewDecoder(reqResp.Body).Decode(&resp); err != nil {
		return tc.Alerts{}, reqInf, errors.New("decoding response: " + err.Error())
	}
	return resp, reqInf, nil
}

func (to *Session) GetDeliveryServiceServers() (tc.DeliveryServiceServerResponse, ReqInf, error) {
	route := apiBase + `/deliveryserviceserver`
	reqResp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.DeliveryServiceServerResponse{}, reqInf, errors.New("requesting from Traffic Ops: " + err.Error())
	}
	defer reqResp.Body.Close()
	resp := tc.DeliveryServiceServerResponse{}
	if err = json.NewDecoder(reqResp.Body).Decode(&resp); err != nil {
		return tc.DeliveryServiceServerResponse{}, reqInf, errors.New("decoding response: " + err.Error())
	}
	return resp, reqInf, nil
}
