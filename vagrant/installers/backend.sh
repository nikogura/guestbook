#!/usr/bin/env bash

echo "upgrading system"
sudo apt-get update
sudo apt-get -y upgrade

echo "Installing Git"
sudo apt-get -y install git

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

sudo mkdir -p /etc/guestbook

echo "$CONFIG" > /etc/guestbook/guestbook.json


INITSCRIPT=$(cat <<'EOF'
# guestbook
#
# Simple Guestbook app

description     "Guestbook"

start on runlevel [2345]
stop on runlevel [!2345]

respawn
respawn limit 10 5
umask 022

console none

exec /go/bin/guestbook run &

EOF
)

echo "$INITSCRIPT" | sudo tee /etc/init/guestbook.conf

sudo service guestbook start

