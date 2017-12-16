#!/usr/bin/env bash

echo "installing db repo"
echo "deb http://apt.postgresql.org/pub/repos/apt/ trusty-pgdg main" | sudo tee -a /etc/apt/sources.list

echo "upgrading system"
DEBIAN_FRONTEND=noninteractive sudo apt-get update && DEBIAN_FRONTEND=noninteractive sudo apt-get upgrade

echo "installing db repo key"
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -

echo "installing database"
DEBIAN_FRONTEND=noninteractive sudo apt-get install --allow-unauthenticated  -y postgresql-all


echo "installing schema and user"
sudo -u postgres /usr/bin/createdb guestbook
sudo -u postgres /usr/bin/createuser guestbook
sudo -u postgres psql -c "ALTER USER guestbook WITH PASSWORD 'guestbook';"


echo "host        guestbook             guestbook             backend            md5" | sudo tee -a /etc/postgresql/9.3/main/pg_hba.conf
sudo sed -i "/# - Connection Settings -/a listen_addresses = '*'" /etc/postgresql/9.3/main/postgresql.conf

sudo service postgresql restart



