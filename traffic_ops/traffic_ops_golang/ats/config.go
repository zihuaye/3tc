package ats

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

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

const configSuffix = ".config"

const HeaderRewritePrefix = "hdr_rw_"
const RegexRemapPrefix = "regex_remap_"
const CacheUrlPrefix = "cacheurl_"

const RemapFile = "remap.config"

const HeaderCommentDateFormat = "Mon Jan 2 15:04:05 MST 2006"

func GetConfigFile(prefix string, xmlId string) string {
	return prefix + xmlId + configSuffix
}

func GetNameVersionString(tx *sql.Tx) (string, error) {
	qry := `
SELECT
  p.name,
  p.value
FROM
  parameter p
WHERE
  (p.name = 'tm.toolname' OR p.name = 'tm.url') AND p.config_file = 'global'
`
	rows, err := tx.Query(qry)
	if err != nil {
		return "", errors.New("querying: " + err.Error())
	}
	defer rows.Close()
	toolName := ""
	url := ""
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return "", errors.New("scanning: " + err.Error())
		}
		if name == "tm.toolname" {
			toolName = val
		} else if name == "tm.url" {
			url = val
		}
	}
	return toolName + " (" + url + ")", nil
}

// getCDNNameFromNameOrID returns the CDN name from a parameter which may be the name or ID.
// This also checks and verifies the existence of the given CDN, and returns an appropriate user error if it doesn't exist.
// Returns the name, any user error, any system error, and any error code.
func getCDNNameFromNameOrID(tx *sql.Tx, cdnNameOrID string) (string, error, error, int) {
	if cdnID, err := strconv.Atoi(cdnNameOrID); err == nil {
		cdnName, ok, err := dbhelpers.GetCDNNameFromID(tx, int64(cdnID))
		if err != nil {
			return "", nil, fmt.Errorf("getting CDN name from id %v: %v", cdnID, err), http.StatusInternalServerError
		} else if !ok {
			return "", errors.New("cdn not found"), nil, http.StatusNotFound
		}
		return string(cdnName), nil, nil, http.StatusOK
	}

	cdnName := cdnNameOrID
	if ok, err := dbhelpers.CDNExists(cdnName, tx); err != nil {
		return "", nil, fmt.Errorf("checking CDN name '%v' existence: %v", cdnName, err), http.StatusInternalServerError
	} else if !ok {
		return "", errors.New("cdn not found"), nil, http.StatusNotFound
	}
	return cdnName, nil, nil, http.StatusOK
}

func headerComment(tx *sql.Tx, name string) (string, error) {
	nameVersionStr, err := GetNameVersionString(tx)
	if err != nil {
		return "", errors.New("getting name version string: " + err.Error())
	}
	return "# DO NOT EDIT - Generated for " + name + " by " + nameVersionStr + " on " + time.Now().Format(HeaderCommentDateFormat) + "\n", nil
}

