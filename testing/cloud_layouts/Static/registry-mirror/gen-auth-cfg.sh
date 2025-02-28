#!/bin/bash

generate_hash() {
    local password=$1
    docker run --rm -e PASSWORD="$password" python:3-alpine \
    sh -c "pip install bcrypt > /dev/null && python -c 'import bcrypt, os; print(bcrypt.hashpw(os.environ[\"PASSWORD\"].encode(), bcrypt.gensalt(rounds=10)).decode())'"
}

MIRROR_HASH=$(generate_hash "$1")
CLUSTER_HASH=$(generate_hash "$2")

cat <<EOF
server:
  addr: ":5001"
token:
  certificate: "/ssl/server.pem"
  key: "/ssl/server.key"
  issuer: "auth-service"
  expiration: 900
users:
  mirror:
    password: |
        ${MIRROR_HASH}
  cluster:
    password: |
        ${CLUSTER_HASH}
acl:
  - match: {account: "mirror"}
    actions: ["*"]
  - match: {account: "cluster"}
    actions: ["pull"]
EOF
