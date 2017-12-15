#!/usr/bin/env bash

echo "upgrading system"
DEBIAN_FRONTEND=noninteractive sudo apt-get update && DEBIAN_FRONTEND=noninteractive sudo apt-get upgrade

echo "Installing Git"
DEBIAN_FRONTEND=noninteractive sudo apt-get install -y git

echo "Installing Go"

wget -q https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz

sudo tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz

echo "Installing Go Tools"

mkdir -p /go/bin
mkdir -p /go/src
mkdir -p /go/pkg

sudo chown -R vagrant:vagrant /go

echo "export PATH=$PATH:/usr/local/go/bin:/go/bin" > /etc/profile.d/go.sh
echo "export GOPATH=/go" >> /etc/profile.d/go.sh

. /etc/profile.d/go.sh

echo "Installing Govendor"
go get github.com/kardianos/govendor

echo "Installing App"
PKG="github.com/nikogura/guestbook"

go get $PKG

CONFIG=$(cat <<EOF
{
  "state": {
    "manager": {
      "type": "gorm",
	  "dialect": "postgres",
      "connect_string": "postgresql://guestbook:guestbook@db:5432/guestbook?sslmode=disable"
    }
  },
  "server": {
    "addr": "0.0.0.0:8080"
  }
}
EOF
)

echo "$CONFIG" > /home/vagrant/guestbook.json


