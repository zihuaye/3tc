#!/usr/bin/env bash
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#
# Defines bash functions to consistently interact with the Traffic Ops API
#
# Build FQDNs
export CDN_FQDN="$CDN_SUBDOMAIN.$TLD_DOMAIN"
export INFRA_FQDN="$INFRA_SUBDOMAIN.$TLD_DOMAIN"
export DB_FQDN="$DB_SERVER.$INFRA_FQDN"
export DNS_FQDN="$DNS_SERVER.$INFRA_FQDN"
export EDGE_FQDN="$EDGE_HOST.$INFRA_FQDN"
export MID_FQDN="$MID_HOST.$INFRA_FQDN"
export ORIGIN_FQDN="$ORIGIN_HOST.$INFRA_FQDN"
export TO_FQDN="$TO_HOST.$INFRA_FQDN"
export TO_PERL_FQDN="$TO_PERL_HOST.$INFRA_FQDN"
export TM_FQDN="$TM_HOST.$INFRA_FQDN"
export TP_FQDN="$TP_HOST.$INFRA_FQDN"
export TR_FQDN="$TR_HOST.$INFRA_FQDN"
export TS_FQDN="$TS_HOST.$INFRA_FQDN"
export TV_FQDN="$TV_HOST.$INFRA_FQDN"

export TO_URL=${TO_URL:-https://$TO_FQDN:$TO_PORT}
export TO_USER=${TO_USER:-$TO_ADMIN_USER}
export TO_PASSWORD=${TO_PASSWORD:-$TO_ADMIN_PASSWORD}

export CURLOPTS=${CURLOPTS:--LfsS}
export CURLAUTH=${CURLAUTH:--k}
export COOKIEJAR=$(mktemp)

login=$(mktemp)

cleanup() {
	rm -f "$COOKIEJAR" "$login"
}

trap cleanup EXIT

cookie_current() {
	local cookiefile=$1
	[[ -s $cookiefile ]] || return 1

	# get expiration from cookiejar -- compare to current time
	exp=$(awk '/mojolicious/ {print $5}' $cookiefile | tail -n 1)
	cur=$(date +%s)

	# return value is the comparison itself
	(( $exp > $cur ))
}

to-auth() {
	# These are required
	if [[ -z $TO_URL || -z $TO_USER || -z $TO_PASSWORD ]]; then
		echo TO_URL TO_USER TO_PASSWORD must all be set
		return 1
	fi

	# if cookiejar is current, nothing to do..
	cookie_current $COOKIEJAR && return

	local url=$TO_URL/api/1.3/user/login
	local datatype='Accept: application/json'
	cat >"$login" <<-CREDS
{ "u" : "$TO_USER", "p" : "$TO_PASSWORD" }
CREDS
	res=$(curl $CURLAUTH $CURLOPTS -H "$datatype" --cookie "$COOKIEJAR" --cookie-jar "$COOKIEJAR" -X POST --data @"$login" "$url")
	if [[ $res != *"Successfully logged in."* ]]; then
		echo "Login failed: $res"
		return 1
	fi
}

to-ping() {
	# ping endpoint does not require authentication
	curl $CURLAUTH $CURLOPTS -X GET "$TO_URL/api/1.3/ping"
}

to-get() {
	to-auth && \
		curl $CURLAUTH $CURLOPTS --cookie "$COOKIEJAR" -X GET "$TO_URL/$1"
}

to-post() {
	local t
	local data
	if [[ -z "$2" ]]; then
		data=""
	elif [[ -f "$2" ]]; then
		data="--data @$2"
	else
		t=$(mktemp)
		echo $2 >$t
		data="--data @$t"
	fi
	to-auth && \
	    curl $CURLAUTH $CURLOPTS --cookie "$COOKIEJAR" -X POST $data "$TO_URL/$1"
	[[ -n $t ]] && rm "$t"    
}

to-put() {
	if [[ -z "$2" ]]; then
		data=""
	elif [[ -f "$2" ]]; then
		data="--data @$2"
	else
		data="--data $2"
	fi
	to-auth && \
	    curl $CURLAUTH $CURLOPTS --cookie "$COOKIEJAR" -X PUT $data "$TO_URL/$1"
}

to-delete() {
	to-auth && \
		curl $CURLAUTH $CURLOPTS --cookie "$COOKIEJAR" -X DELETE "$TO_URL/$1"
}

# Constructs a server's JSON definiton and places it into the enroller's structure for loading
# args:
#         serverType - the type of the server to be created; one of "edge", "mid", "tm"
to-enroll() {

	while true; do 
		[ -d "$ENROLLER_DIR" ] && break
		echo "Waiting for $ENROLLER_DIR ..."
		sleep 2
	done

	while true; do 
		[ "$serverType" = "to" ] && break
		[ -f "$ENROLLER_DIR/initial-load-done" ] && break
		echo "Waiting for traffic-ops to do initial load ..."
		sleep 2
	done
	if [[ ! -d ${ENROLLER_DIR}/servers ]]; then
		echo "${ENROLLER_DIR}/servers not found -- contents:"
		find ${ENROLLER_DIR} -ls
	fi
	local serverType="$1"

	if [[ ! -z "$2" ]]; then
		export MY_CDN="$2"
	else
		export MY_CDN="CDN-in-a-Box"
	fi

	export MY_NET_INTERFACE='eth0'
	export MY_HOSTNAME="$(hostname -s)"
	export MY_DOMAINNAME="$(dnsdomainname)"
	export MY_IP="$(ifconfig $MY_NET_INTERFACE | grep 'inet ' | tr -s ' ' | cut -d ' ' -f 3)"
	export MY_GATEWAY="$(route -n | grep $MY_NET_INTERFACE | grep -E '^0\.0\.0\.0' | tr -s ' ' | cut -d ' ' -f2)"
	export MY_NETMASK="$(ifconfig $MY_NET_INTERFACE | grep 'inet ' | tr -s ' ' | cut -d ' ' -f 5)"
	export MY_IP6_ADDRESS="$(ifconfig $MY_NET_INTERFACE | grep inet6 | grep global | awk '{ print $2 }')"
	export MY_IP6_GATEWAY="$(route -n6 | grep UG | awk '{print $2}')"

	case "$serverType" in
		"edge" )
			export MY_TYPE="EDGE"
			export MY_PROFILE="ATS_EDGE_TIER_CACHE"
			export MY_STATUS="REPORTED"
			if [[ ! -z "$3" ]]; then
				export MY_CACHE_GROUP="$3"
			else
				export MY_CACHE_GROUP="CDN_in_a_Box_Edge"
			fi
			;;
		"mid" )
			export MY_TYPE="MID"
			export MY_PROFILE="ATS_MID_TIER_CACHE"
			export MY_STATUS="REPORTED"
			if [[ ! -z "$3" ]]; then
				export MY_CACHE_GROUP="$3"
			else
				export MY_CACHE_GROUP="CDN_in_a_Box_Mid"
			fi
			;;
		"tm" )
			export MY_TYPE="RASCAL"
			export MY_PROFILE="RASCAL-Traffic_Monitor"
			export MY_STATUS="ONLINE"
			if [[ ! -z "$3" ]]; then
				export MY_CACHE_GROUP="$3"
			else
				export MY_CACHE_GROUP="CDN_in_a_Box_Edge"
			fi
			;;
		"to" ) 
			export MY_TYPE="TRAFFIC_OPS"
			export MY_PROFILE="TRAFFIC_OPS"
			export MY_STATUS="ONLINE"
			if [[ ! -z "$3" ]]; then
				export MY_CACHE_GROUP="$3"
			else
				export MY_CACHE_GROUP="CDN_in_a_Box_Edge"
			fi
			;;
		"tr" )
			export MY_TYPE="CCR"
			export MY_PROFILE="CCR_CIAB"
			export MY_STATUS="ONLINE"
			if [[ ! -z "$3" ]]; then
				export MY_CACHE_GROUP="$3"
			else
				export MY_CACHE_GROUP="CDN_in_a_Box_Edge"
			fi
			;;
		"tp" )
			export MY_TYPE="TRAFFIC_PORTAL"
			export MY_PROFILE="TRAFFIC_PORTAL"
			export MY_STATUS="ONLINE"
			if [[ ! -z "$3" ]]; then
				export MY_CACHE_GROUP="$3"
			else
				export MY_CACHE_GROUP="CDN_in_a_Box_Edge"
			fi
			;;
		"tv" )
			export MY_TYPE="RIAK"
			export MY_PROFILE="RIAK_ALL"
			export MY_STATUS="ONLINE"
			if [[ ! -z "$3" ]]; then
				export MY_CACHE_GROUP="$3"
			else
				export MY_CACHE_GROUP="CDN_in_a_Box_Edge"
			fi
			;;
		* )
			echo "Usage: to-enroll SERVER_TYPE" >&2
			echo "(SERVER_TYPE must be a recognized server type)" >&2
			return 1
			;;
	esac

	# replace env references in the file
	envsubst < "/server_template.json" > "${ENROLLER_DIR}/servers/$HOSTNAME.json"

	sleep 3
}
