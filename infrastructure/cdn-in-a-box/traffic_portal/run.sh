#!/usr/bin/env bash
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


NAME="Traffic Portal Application"
NODE_BIN_DIR="/usr/bin"
NODE_PATH="/opt/traffic_portal/node_modules"
FOREVER_BIN_DIR="/opt/traffic_portal/node_modules/forever/bin"
APPLICATION_PATH="/opt/traffic_portal/server.js"
PIDFILE="/var/run/traffic_portal.pid"
LOGFILE="/var/log/traffic_portal/traffic_portal.log"
MIN_UPTIME="5000"
SPIN_SLEEP_TIME="2000"

source /to-access.sh

# Wait on SSL certificate generation
until [ -f "$X509_CA_DONE_FILE" ] 
do
  echo "Waiting on Shared SSL certificate generation"
  sleep 3
done

# Source the CIAB-CA shared SSL environment
source $X509_CA_ENV_FILE

# Trust the CIAB-CA at the System level
cp $X509_CA_CERT_FILE /etc/pki/ca-trust/source/anchors
update-ca-trust extract

# Configuration of Traffic Portal
key=$X509_INFRA_KEY_FILE
cert=$X509_INFRA_CERT_FILE
ca=/etc/pki/tls/certs/ca-bundle.crt

# set configs to point to TO_FQDN
sed -i -e "/^\s*base_url:/ s@'.*'@'https://$TO_FQDN:$TO_PORT/api/'@" /etc/traffic_portal/conf/config.js
sed -i -e "/^\s*cert:/ s@'.*'@'$cert'@" /etc/traffic_portal/conf/config.js
sed -i -e "/^\s*key:/ s@'.*'@'$key'@" /etc/traffic_portal/conf/config.js

props=/opt/traffic_portal/public/traffic_portal_properties.json
tmp=$(mktemp)

echo "TO_HOST: $TO_HOST"
echo "TO_HOST: $TO_PORT"
echo "TO_FQDN: $TO_FQDN"

jq --arg TO_FQDN "$TO_FQDN:$TO_PORT" '.properties.api.baseUrl = "https://"+$TO_FQDN' <$props >$tmp
mv $tmp $props

# Enroll the Traffic Portal
to-enroll "tp" ALL || (while true; do echo "enroll failed."; sleep 3 ; done)

# Add node to the path for situations in which the environment is passed.
PATH=$FOREVER_BIN_DIR:$NODE_BIN_DIR:$PATH
forever \
    --pidFile $PIDFILE \
    -a \
    -l $LOGFILE \
    --minUptime $MIN_UPTIME \
    --spinSleepTime $SPIN_SLEEP_TIME \
    start $APPLICATION_PATH

tail -f /dev/null
