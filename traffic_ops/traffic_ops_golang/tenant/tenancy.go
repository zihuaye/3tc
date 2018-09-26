package tenant

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

// tenancy.go defines methods and functions to determine tenancy of resources.

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
)

// DeliveryServiceTenantInfo provides only deliveryservice info needed here
type DeliveryServiceTenantInfo tc.DeliveryServiceNullable

// IsTenantAuthorized returns true if the user has tenant access on this tenant
func (dsInfo DeliveryServiceTenantInfo) IsTenantAuthorized(user *auth.CurrentUser, tx *sql.Tx) (bool, error) {
	if dsInfo.TenantID == nil {
		return false, errors.New("TenantID is nil")
	}
	return IsResourceAuthorizedToUserTx(*dsInfo.TenantID, user, tx)
}

// GetDeliveryServiceTenantInfo returns tenant information for a deliveryservice
func GetDeliveryServiceTenantInfo(xmlID string, tx *sql.Tx) (*DeliveryServiceTenantInfo, error) {
	ds := DeliveryServiceTenantInfo{}
	ds.XMLID = util.StrPtr(xmlID)
	if err := tx.QueryRow(`SELECT tenant_id FROM deliveryservice where xml_id = $1`, &ds.XMLID).Scan(&ds.TenantID); err != nil {
		if err == sql.ErrNoRows {
			return &ds, errors.New("a deliveryservice with xml_id '" + xmlID + "' was not found")
		}
		return nil, errors.New("querying tenant id from delivery service: " + err.Error())
	}
	return &ds, nil
}

// Check checks that the given user has access to the given XMLID. Returns a user error, system error,
// and the HTTP status code to be returned to the user if an error occurred. On success, the user error
// and system error will both be nil, and the error code should be ignored.
func Check(user *auth.CurrentUser, XMLID string, tx *sql.Tx) (error, error, int) {
	dsInfo, err := GetDeliveryServiceTenantInfo(XMLID, tx)
	if err != nil {
		if dsInfo == nil {
			return nil, errors.New("deliveryservice lookup failure: " + err.Error()), http.StatusInternalServerError
		}
		return errors.New("no such deliveryservice: '" + XMLID + "'"), nil, http.StatusBadRequest
	}
	hasAccess, err := dsInfo.IsTenantAuthorized(user, tx)
	if err != nil {
		return nil, errors.New("user tenancy check failure: " + err.Error()), http.StatusInternalServerError
	}
	if !hasAccess {
		return nil, errors.New("Access to this resource is not authorized"), http.StatusForbidden
	}
	return nil, nil, http.StatusOK
}

// CheckID checks that the given user has access to the given delivery service. Returns a user error, a system error, and an HTTP error code. If both the user and system error are nil, the error code should be ignored.
func CheckID(tx *sql.Tx, user *auth.CurrentUser, dsID int) (error, error, int) {
	ok, err := IsTenancyEnabledTx(tx)
	if err != nil {
		return nil, errors.New("checking tenancy enabled: " + err.Error()), http.StatusInternalServerError
	}
	if !ok {
		return nil, nil, http.StatusOK
	}

	dsTenantID, ok, err := getDSTenantIDByIDTx(tx, dsID)
	if err != nil {
		return nil, errors.New("checking tenant: " + err.Error()), http.StatusInternalServerError
	}
	if !ok {
		return errors.New("delivery service " + strconv.Itoa(dsID) + " not found"), nil, http.StatusNotFound
	}
	if dsTenantID == nil {
		return nil, nil, http.StatusOK
	}

	authorized, err := IsResourceAuthorizedToUserTx(*dsTenantID, user, tx)
	if err != nil {
		return nil, errors.New("checking tenant: " + err.Error()), http.StatusInternalServerError
	}
	if !authorized {
		return errors.New("not authorized on this tenant"), nil, http.StatusForbidden
	}
	return nil, nil, http.StatusOK
}

// GetUserTenantListTx returns a Tenant list that the specified user has access to.
// NOTE: This method does not use the use_tenancy parameter and if this method is being used
// to control tenancy the parameter must be checked. The method IsResourceAuthorizedToUser checks the use_tenancy parameter
// and should be used for this purpose in most cases.
func GetUserTenantListTx(user auth.CurrentUser, tx *sql.Tx) ([]tc.TenantNullable, error) {
	query := `WITH RECURSIVE q AS (SELECT id, name, active, parent_id, last_updated FROM tenant WHERE id = $1
	UNION SELECT t.id, t.name, t.active, t.parent_id, t.last_updated  FROM tenant t JOIN q ON q.id = t.parent_id)
	SELECT id, name, active, parent_id, last_updated FROM q;`

	rows, err := tx.Query(query, user.TenantID)
	if err != nil {
		return nil, errors.New("querying user tenant list: " + err.Error())
	}
	defer rows.Close()

	tenants := []tc.TenantNullable{}
	for rows.Next() {
		t := tc.TenantNullable{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Active, &t.ParentID, &t.LastUpdated); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	return tenants, nil
}

func GetUserTenantIDListTx(tx *sql.Tx, userTenantID int) ([]int, error) {
	query := `
WITH RECURSIVE q AS (SELECT id, name, active, parent_id FROM tenant WHERE id = $1
UNION SELECT t.id, t.name, t.active, t.parent_id  FROM tenant t JOIN q ON q.id = t.parent_id)
SELECT id FROM q;
`
	rows, err := tx.Query(query, userTenantID)
	if err != nil {
		return nil, errors.New("querying user tenant ID list: " + err.Error())
	}
	defer rows.Close()

	tenants := []int{}
	for rows.Next() {
		tenantID := 0
		if err := rows.Scan(&tenantID); err != nil {
			return nil, err
		}
		tenants = append(tenants, tenantID)
	}
	return tenants, nil
}

// IsTenancyEnabledTx returns true if tenancy is enabled or false otherwise
func IsTenancyEnabledTx(tx *sql.Tx) (bool, error) {
	query := `SELECT COALESCE(value::boolean,FALSE) AS value FROM parameter WHERE name = 'use_tenancy' AND config_file = 'global' UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY`
	useTenancy := false
	err := tx.QueryRow(query).Scan(&useTenancy)
	if err != nil {
		return false, errors.New("checking if tenancy is enabled: " + err.Error())
	}
	return useTenancy, nil
}

// IsResourceAuthorizedToUserTx returns a boolean value describing if the user has access to the provided resource tenant id and an error
// if use_tenancy is set to false (0 in the db) this method will return true allowing access.
func IsResourceAuthorizedToUserTx(resourceTenantID int, user *auth.CurrentUser, tx *sql.Tx) (bool, error) {
	// $1 is the user tenant ID and $2 is the resource tenant ID
	query := `WITH RECURSIVE q AS (SELECT id, active FROM tenant WHERE id = $1
	UNION SELECT t.id, t.active FROM TENANT t JOIN q ON q.id = t.parent_id),
	tenancy AS (SELECT COALESCE(value::boolean,FALSE) AS value FROM parameter WHERE name = 'use_tenancy' AND config_file = 'global' UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY)
	SELECT id, active, tenancy.value AS use_tenancy FROM tenancy, q WHERE id = $2 UNION ALL SELECT -1, false, tenancy.value AS use_tenancy FROM tenancy FETCH FIRST 1 ROW ONLY;`

	var tenantID int
	var active bool
	var useTenancy bool

	log.Debugln("\nQuery: ", query)
	err := tx.QueryRow(query, user.TenantID, resourceTenantID).Scan(&tenantID, &active, &useTenancy)

	switch {
	case err != nil:
		log.Errorf("Error checking user tenant %v access on resourceTenant  %v: %v", user.TenantID, resourceTenantID, err.Error())
		return false, err
	default:
		if !useTenancy {
			return true, nil
		}
		if active && tenantID == resourceTenantID {
			return true, nil
		} else {
			fmt.Printf("default")
			return false, nil
		}
	}
}

// getDSTenantIDByIDTx returns the tenant ID, whether the delivery service exists, and any error.
// Note the id may be nil, even if true is returned, if the delivery service exists but its tenant_id field is null.
// TODO move somewhere generic
func getDSTenantIDByIDTx(tx *sql.Tx, id int) (*int, bool, error) {
	tenantID := (*int)(nil)
	if err := tx.QueryRow(`SELECT tenant_id FROM deliveryservice where id = $1`, id).Scan(&tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("querying tenant ID for delivery service ID '%v': %v", id, err)
	}
	return tenantID, true, nil
}

type Tenanter interface {
	TenantID() *int
	Name() string
	GetType() string
}

// FilterAuthorized takes a slice of objects, and returns only the objects the given user is authorized for.
func FilterAuthorized(objs []Tenanter, user *auth.CurrentUser, tx *sql.Tx) ([]interface{}, error) {
	tenancyEnabled, err := IsTenancyEnabledTx(tx)
	if err != nil {
		return nil, errors.New("Error checking if tenancy enabled.")
	}
	if !tenancyEnabled {
		newObjs := []interface{}{}
		for _, obj := range objs {
			newObjs = append(newObjs, obj)
		}
		return newObjs, nil
	}

	newObjs := []interface{}{}
	for _, obj := range objs {
		if obj.TenantID() == nil {
			return nil, fmt.Errorf("FilterAuthorized for %T %s %s: no tenant ID", obj, obj.Name(), obj.GetType())
		}
		// TODO add/use a helper func to make a single SQL call, for performance
		ok, err := IsResourceAuthorizedToUserTx(*obj.TenantID(), user, tx)
		if err != nil {
			return nil, fmt.Errorf("FilterAuthorized for %T %s %s: no tenant ID", obj, obj.Name(), obj.GetType())
		}
		if !ok {
			continue
		}
		newObjs = append(newObjs, obj)
	}
	return newObjs, nil
}
