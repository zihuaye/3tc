package status

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
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	validation "github.com/go-ozzo/ozzo-validation"
)

//we need a type alias to define functions on
type TOStatus struct {
	api.APIInfoImpl `json:"-"`
	tc.StatusNullable
}

func (v *TOStatus) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOStatus) InsertQuery() string           { return insertQuery() }
func (v *TOStatus) NewReadObj() interface{}       { return &tc.Status{} }
func (v *TOStatus) SelectQuery() string           { return selectQuery() }
func (v *TOStatus) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"id":          dbhelpers.WhereColumnInfo{"id", api.IsInt},
		"description": dbhelpers.WhereColumnInfo{"description", nil},
		"name":        dbhelpers.WhereColumnInfo{"name", nil},
	}
}
func (v *TOStatus) UpdateQuery() string { return updateQuery() }
func (v *TOStatus) DeleteQuery() string { return deleteQuery() }

func (status TOStatus) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (status TOStatus) GetKeys() (map[string]interface{}, bool) {
	if status.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *status.ID}, true
}

func (status *TOStatus) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	status.ID = &i
}

func (status TOStatus) GetAuditName() string {
	if status.Name != nil {
		return *status.Name
	}
	if status.ID != nil {
		return strconv.Itoa(*status.ID)
	}
	return "unknown"
}

func (status TOStatus) GetType() string { return "status" }

func (status TOStatus) Validate() error {
	errs := validation.Errors{
		"name": validation.Validate(status.Name, validation.NotNil, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
}

func (st *TOStatus) Read() ([]interface{}, error, error, int) { return api.GenericRead(st) }
func (st *TOStatus) Update() (error, error, int)              { return api.GenericUpdate(st) }
func (st *TOStatus) Create() (error, error, int)              { return api.GenericCreate(st) }
func (st *TOStatus) Delete() (error, error, int)              { return api.GenericDelete(st) }

func selectQuery() string {
	return `
SELECT
  description,
  id,
  last_updated,
  name
FROM
  status s
`
}

func updateQuery() string {
	return `UPDATE
status SET
name=:name,
description=:description
WHERE id=:id RETURNING last_updated`
}

func insertQuery() string {
	return `INSERT INTO status (
name,
description) VALUES (
:name,
:description) RETURNING id,last_updated`
}

func deleteQuery() string {
	return `DELETE FROM status WHERE id=:id`
}
