package dbhelpers

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
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type WhereColumnInfo struct {
	Column  string
	Checker func(string) error
}

const BaseWhere = "\nWHERE"
const BaseOrderBy = "\nORDER BY"

func BuildWhereAndOrderBy(parameters map[string]string, queryParamsToSQLCols map[string]WhereColumnInfo) (string, string, map[string]interface{}, []error) {
	whereClause := BaseWhere
	orderBy := BaseOrderBy
	var criteria string
	var queryValues map[string]interface{}
	var errs []error
	criteria, queryValues, errs = parseCriteriaAndQueryValues(queryParamsToSQLCols, parameters)

	if len(queryValues) > 0 {
		whereClause += " " + criteria
	}
	if len(errs) > 0 {
		return "", "", queryValues, errs
	}

	if orderby, ok := parameters["orderby"]; ok {
		log.Debugln("orderby: ", orderby)
		if colInfo, ok := queryParamsToSQLCols[orderby]; ok {
			log.Debugln("orderby column ", colInfo)
			orderBy += " " + colInfo.Column
		} else {
			log.Debugln("Incorrect name for orderby: ", orderby)
		}
	}
	if whereClause == BaseWhere {
		whereClause = ""
	}
	if orderBy == BaseOrderBy {
		orderBy = ""
	}
	log.Debugf("\n--\n Where: %s \n Order By: %s", whereClause, orderBy)
	return whereClause, orderBy, queryValues, errs
}

func parseCriteriaAndQueryValues(queryParamsToSQLCols map[string]WhereColumnInfo, parameters map[string]string) (string, map[string]interface{}, []error) {
	m := make(map[string]interface{})
	var criteria string

	var criteriaArgs []string
	errs := []error{}
	queryValues := make(map[string]interface{})
	for key, colInfo := range queryParamsToSQLCols {
		if urlValue, ok := parameters[key]; ok {
			var err error
			if colInfo.Checker != nil {
				err = colInfo.Checker(urlValue)
			}
			if err != nil {
				errs = append(errs, errors.New(key+" "+err.Error()))
			} else {
				m[key] = urlValue
				criteria = colInfo.Column + "=:" + key
				criteriaArgs = append(criteriaArgs, criteria)
				queryValues[key] = urlValue
			}
		}
	}
	criteria = strings.Join(criteriaArgs, " AND ")

	return criteria, queryValues, errs
}

// AddTenancyCheck takes a WHERE clause (can be ""), the associated queryValues (can be empty),
// a tenantColumnName that should provide a bigint corresponding to the tenantID of the object being checked (this may require a CAST),
// and an array of the tenantIDs the user has access to; it returns a where clause and associated queryValues including filtering based on tenancy.
func AddTenancyCheck(where string, queryValues map[string]interface{}, tenantColumnName string, tenantIDs []int) (string, map[string]interface{}) {
	if where == "" {
		where = BaseWhere + " " + tenantColumnName + " = ANY(CAST(:accessibleTenants AS bigint[]))"
	} else {
		where += " AND " + tenantColumnName + " = ANY(CAST(:accessibleTenants AS bigint[]))"
	}

	queryValues["accessibleTenants"] = pq.Array(tenantIDs)

	return where, queryValues
}

// GetDSNameFromID loads the DeliveryService's xml_id from the database, from the ID. Returns whether the delivery service was found, and any error.
func GetDSNameFromID(tx *sql.Tx, id int) (tc.DeliveryServiceName, bool, error) {
	name := tc.DeliveryServiceName("")
	if err := tx.QueryRow(`SELECT xml_id FROM deliveryservice WHERE id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return tc.DeliveryServiceName(""), false, nil
		}
		return tc.DeliveryServiceName(""), false, fmt.Errorf("querying xml_id for delivery service ID '%v': %v", id, err)
	}
	return name, true, nil
}

// GetProfileNameFromID returns the profile's name, whether a profile with ID exists, or any error.
func GetProfileNameFromID(id int, tx *sql.Tx) (string, bool, error) {
	name := ""
	if err := tx.QueryRow(`SELECT name from profile where id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying profile name from id: " + err.Error())
	}
	return name, true, nil
}

// GetProfileIDFromName returns the profile's ID, whether a profile with name exists, or any error.
func GetProfileIDFromName(name string, tx *sql.Tx) (int, bool, error) {
	id := 0
	if err := tx.QueryRow(`SELECT id from profile where name = $1`, name).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, errors.New("querying profile id from name: " + err.Error())
	}
	return id, true, nil
}

// Returns true if the cdn exists
func CDNExists(cdnName string, tx *sqlx.Tx) (bool, error) {
	var id int
	if err := tx.QueryRow(`SELECT id FROM cdn WHERE name = $1`, cdnName).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.New("Error querying CDN name: " + err.Error())
	}
	return true, nil
}
