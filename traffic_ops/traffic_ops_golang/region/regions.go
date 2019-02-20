package region

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
	"errors"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

//we need a type alias to define functions on
type TORegion struct {
	api.APIInfoImpl `json:"-"`
	tc.Region
}

func (v *TORegion) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = t }
func (v *TORegion) InsertQuery() string           { return insertQuery() }
func (v *TORegion) NewReadObj() interface{}       { return &tc.Region{} }
func (v *TORegion) SelectQuery() string           { return selectQuery() }
func (v *TORegion) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name":     dbhelpers.WhereColumnInfo{"r.name", nil},
		"division": dbhelpers.WhereColumnInfo{"r.division", nil},
		"id":       dbhelpers.WhereColumnInfo{"r.id", api.IsInt},
	}
}
func (v *TORegion) UpdateQuery() string { return updateQuery() }
func (v *TORegion) DeleteQuery() string { return deleteQuery() }

func (region TORegion) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (region TORegion) GetKeys() (map[string]interface{}, bool) {
	return map[string]interface{}{"id": region.ID}, true
}

func (region *TORegion) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	region.ID = i
}

func (region *TORegion) GetAuditName() string {
	return region.Name
}

func (region *TORegion) GetType() string {
	return "region"
}

func (region *TORegion) Validate() error {
	if len(region.Name) < 1 {
		return errors.New(`Region 'name' is required.`)
	}
	return nil
}

func (rg *TORegion) Read() ([]interface{}, error, error, int) { return api.GenericRead(rg) }
func (rg *TORegion) Update() (error, error, int)              { return api.GenericUpdate(rg) }
func (rg *TORegion) Create() (error, error, int)              { return api.GenericCreate(rg) }
func (rg *TORegion) Delete() (error, error, int)              { return api.GenericDelete(rg) }

func selectQuery() string {
	return `SELECT
r.division,
d.name as divisionname,
r.id,
r.last_updated,
r.name
FROM region r
JOIN division d ON r.division = d.id`
}

func updateQuery() string {
	query := `UPDATE
region SET
division=:division,
name=:name
WHERE id=:id RETURNING last_updated`
	return query
}

func insertQuery() string {
	query := `INSERT INTO region (
division,
name) VALUES (
:division,
:name) RETURNING id,last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM region
WHERE id=:id`
	return query
}
