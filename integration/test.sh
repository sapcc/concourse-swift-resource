#!/bin/bash

set -e
set -o pipefail
#ensure we are in the integration test directory
cd "$( dirname "${BASH_SOURCE[0]}" )"

DIFF=diff
if type colordiff>/dev/null; then
  DIFF=colordiff
fi

export OS_USERNAME=$(jq -r .username resource.json)
export OS_AUTH_URL=$(jq -r .auth_url resource.json)
export OS_PASSWORD=$(jq -r .password resource.json)
export OS_USER_DOMAIN_NAME=$(jq -r .domain resource.json)
export OS_PROJECT_DOMAIN_NAME=$(jq -r .domain resource.json)
export OS_PROJECT_ID=$(jq -r .tenant_id resource.json)
export SWIFTCLIENT_INSECURE=$(jq -r .disable_tls_verify resource.json)
CONTAINER=$(jq -r .container resource.json)

CHECK() {
  jq "$1" resource.json | ../bin/check
}

OUT() {
  jq "$1" resource.json | ../bin/out . 
}

IN() {
  mkdir -p $2
  jq "$1" resource.json | ../bin/in $2
}

MATCH() {
  if ! jq --arg a "$1" --arg b "$2" -n -e '$a == $b' > /dev/null; then
    colordiff -y <(echo "$1" |jq .) <(echo "$2" | jq .)
    exit 1
  fi
}

#Clean objects from previous integration tests
objects=$(swift list $CONTAINER)
[ -n "$objects" ] && swift delete $CONTAINER $objects >/dev/null

echo "Check empty container"
expected='[]'
response=$(CHECK '{source:.}')
MATCH "$expected" "$response"

echo "Putting file_0.1.0"
expected='{"version":{"path":"file_0.1.0"},"metadata":[{"name":"Version","value":"0.1.0"},{"name":"Size","value":"4"}]}'
response=$(OUT '{source:., params: {from: "out/file_0.1.0"}}')
MATCH "$expected" "$response"

echo "Putting file_0.2.0"
expected='{"version":{"path":"file_0.2.0"},"metadata":[{"name":"Version","value":"0.2.0"},{"name":"Size","value":"4"}]}'
response=$(OUT '{source:., params: {from: "out/file_0.2.0"}}')
MATCH "$expected" "$response"

echo Check without version
expected='[{"path":"file_0.2.0"}]'
response=$(CHECK '{source:.}')
MATCH "$expected" "$response"

echo Check for newer versions than file_0.2.0 
expected='[]'
response=$(CHECK '{source:., version:{path:"file_0.2.0"}}')
MATCH "$expected" "$response"

echo Check for versions newer than file_0.0.1
expected='[{"path":"file_0.1.0"},{"path":"file_0.2.0"}]'
response=$(CHECK '{source:., version:{path:"file_0.0.1"}}')
MATCH "$expected" "$response"

rm -rf in/file_* in/version in/filename
echo Get version file_0.2.0
last_modified=$(swift stat $CONTAINER file_0.2.0 | sed -n 's/.*Last Modified: \(.*\)/\1/p')
expected='{"version":{"path":"file_0.2.0"},"metadata":[{"name":"Version","value":"0.2.0"},{"name":"Size","value":"4"},{"name":"Last Modified","value":"'$last_modified'"}]}'
response=$(IN '{source:., version:{path:"file_0.2.0"}}' in/)
MATCH "$expected" "$response"
ls in/file_0.2.0 in/version in/filename > /dev/null



