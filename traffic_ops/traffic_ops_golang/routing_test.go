package main

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"fmt"

	"bytes"
	"context"
	"net/http/httptest"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
)

type key int

const AuthWasCalled key = iota

type routeTest struct {
	Method      string
	Path        string
	ExpectMatch bool
	Params      map[string]string
}

// TODO: This should be expanded to include POST/PUT/DELETE and other params
var testRoutes = []routeTest{
	routeTest{
		Method:      `GET`,
		Path:        `api/1.1/cdns`,
		ExpectMatch: true,
		Params:      map[string]string{},
	},
	routeTest{
		Method:      `POST`,
		Path:        `api/1.4/users/login`,
		ExpectMatch: false,
		Params:      map[string]string{},
	},
	routeTest{
		Method:      `POST`,
		Path:        `api/1.1/cdns`,
		ExpectMatch: true,
		Params:      map[string]string{},
	},
	routeTest{
		Method:      `POST`,
		Path:        `api/1.1/users`,
		ExpectMatch: true,
		Params:      map[string]string{},
	},
	routeTest{
		Method:      `PUT`,
		Path:        `api/1.1/deliveryservices/3`,
		ExpectMatch: true,
		Params:      map[string]string{"id": "3"},
	},
	routeTest{
		Method:      `DELETE`,
		Path:        `api/1.1/servers/777`,
		ExpectMatch: true,
		Params:      map[string]string{"id": "777"},
	},
	routeTest{
		Method:      `GET`,
		Path:        `api/1.4/cdns/1`,
		ExpectMatch: true,
		Params:      map[string]string{"id": "1"},
	},
	routeTest{
		Method:      `GET`,
		Path:        `api/1.4/notatypeweknowabout`,
		ExpectMatch: false,
		Params:      map[string]string{},
	},
	routeTest{
		Method:      `GET`,
		Path:        `api/1.3/asns.json`,
		ExpectMatch: true,
		Params:      map[string]string{},
	},
	routeTest{
		Method:      `GET`,
		Path:        `api/99999.99999/cdns`,
		ExpectMatch: false,
		Params:      map[string]string{},
	},
	routeTest{
		Method:      `GET`,
		Path:        `blahblah/api/1.2/cdns`,
		ExpectMatch: false,
		Params:      map[string]string{},
	},
	routeTest{
		Method:      `GET`,
		Path:        `internal/api/1.2/federations.json`,
		ExpectMatch: false,
		Params:      map[string]string{},
	},
}

func TestCompileRoutes(t *testing.T) {
	url, err := url.Parse("https://to.test")
	if err != nil {
		t.Error("error parsing test url")
	}
	d := ServerData{Config: config.Config{URL: url, Secrets: []string{"n0SeCr3t$"}}}
	// TODO: not currently checking rawRoutes or catchall
	routeSlice, _ /*rawRoutes*/, _ /*catchall*/, err := Routes(d)
	if err != nil {
		t.Error("error fetching routes: ", err.Error())
	}

	authBase := AuthBase{secret: d.Secrets[0], override: nil}
	routes, versions := CreateRouteMap(routeSlice, nil, authBase, 1)
	if len(routes) == 0 {
		t.Error("no routes handler defined")
	}
	if len(versions) == 0 {
		t.Error("no versions defined")
	}
	compiledRoutes := CompileRoutes(routes)

	for _, rt := range testRoutes {
		t.Logf("testing path %s %s", rt.Method, rt.Path)
		var found bool
		params := map[string]string{}
		for _, compiledRoute := range compiledRoutes[rt.Method] {
			match := compiledRoute.Regex.FindStringSubmatch(rt.Path)
			if len(match) == 0 {
				continue
			}
			found = true
			for i, v := range compiledRoute.Params {
				params[v] = match[i+1]
			}
		}
		if found != rt.ExpectMatch {
			if rt.ExpectMatch {
				t.Errorf("expected %s %s to have a route match", rt.Method, rt.Path)
			} else {
				t.Errorf("expected %s %s to have no route match", rt.Method, rt.Path)
			}
			continue
		}
		if !reflect.DeepEqual(params, rt.Params) {
			t.Errorf("%s %s: expected params %v, got %v", rt.Method, rt.Path, rt.Params, params)
		}
	}
}

func TestCreateRouteMap(t *testing.T) {
	authBase := AuthBase{"secret", func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), AuthWasCalled, "true")
			handlerFunc(w, r.WithContext(ctx))
		}
	}}

	PathOneHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authWasCalled := getAuthWasCalled(ctx)
		fmt.Fprintf(w, "%s %s", "path1", authWasCalled)
	}

	PathTwoHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authWasCalled := getAuthWasCalled(ctx)
		fmt.Fprintf(w, "%s %s", "path2", authWasCalled)
	}

	PathThreeHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authWasCalled := getAuthWasCalled(ctx)
		fmt.Fprintf(w, "%s %s", "path3", authWasCalled)
	}

	routes := []Route{
		{1.2, http.MethodGet, `path1`, PathOneHandler, auth.PrivLevelReadOnly, true, nil},
		{1.2, http.MethodGet, `path2`, PathTwoHandler, 0, false, nil},
		{1.2, http.MethodGet, `path3`, PathThreeHandler, 0, false, []Middleware{}},
	}

	rawRoutes := []RawRoute{}
	routeMap, _ := CreateRouteMap(routes, rawRoutes, authBase, 60)

	route1Handler := routeMap["GET"][0].Handler

	w := httptest.NewRecorder()
	r, err := http.NewRequest("", "/", nil)
	if err != nil {
		t.Error("Error creating new request")
	}

	route1Handler(w, r)

	if bytes.Compare(w.Body.Bytes(), []byte("path1 true")) != 0 {
		t.Errorf("Got: %s \nExpected to receive path1 true\n", w.Body.Bytes())
	}

	route2Handler := routeMap["GET"][1].Handler

	w = httptest.NewRecorder()

	route2Handler(w, r)

	if bytes.Compare(w.Body.Bytes(), []byte("path2 false")) != 0 {
		t.Errorf("Got: %s \nExpected to receive path2 false\n", w.Body.Bytes())
	}

	if v, ok := w.HeaderMap["Access-Control-Allow-Credentials"]; !ok || len(v) != 1 || v[0] != "true" {
		t.Errorf(`Expected Access-Control-Allow-Credentials: [ "true" ]`)
	}

	route3Handler := routeMap["GET"][2].Handler
	w = httptest.NewRecorder()
	route3Handler(w, r)
	if bytes.Compare(w.Body.Bytes(), []byte("path3 false")) != 0 {
		t.Errorf("Got: %s \nExpected to receive path3 false\n", w.Body.Bytes())
	}
	if v, ok := w.HeaderMap["Access-Control-Allow-Credentials"]; ok {
		t.Errorf("Unexpected Access-Control-Allow-Credentials: %s", v)
	}
}

func getAuthWasCalled(ctx context.Context) string {
	val := ctx.Value(AuthWasCalled)
	if val != nil {
		return val.(string)
	}
	return "false"
}
