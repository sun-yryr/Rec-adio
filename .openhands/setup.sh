#!/bin/bash

CURRENT_DIR=$(cd $(dirname $0); pwd)

sudo apt-get update
mkdir -p /tmp/setup
cd /tmp/setup
rm -rf /usr/local/go
wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
source ~/.bashrc

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
  | sh -s -- -b $(go env GOPATH)/bin v2.1.6

go install github.com/nats-io/natscli/nats@latest

BIN="/usr/local/bin" && \
  VERSION="1.53.0" && \
  curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m)" \
    -o "${BIN}/buf" && \
  chmod +x "${BIN}/buf"

cd $CURRENT_DIR

go build ./...
