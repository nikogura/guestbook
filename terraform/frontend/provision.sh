#!/usr/bin/env bash

set -e

export DEBIAN_FRONTEND=noninteractive

echo "upgrading system"
sudo apt-get update

sudo sudo apt-get -y upgrade

echo "Installing Git"
sudo apt-get -y install nginx

NGINX=$(cat <<'EOF'
server {
    listen 80 default_server;
    server_name localhost;

    root /usr/share/nginx/html;
    index index.html index.htm;

    location /guestbook/ {
        proxy_pass http://{{BACKEND_ELB}}:8080/guestbook/;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}

EOF
)

echo -e "$NGINX" | sudo tee /etc/nginx/sites-available/guestbook

sudo ln -s /etc/nginx/sites-available/guestbook /etc/nginx/sites-enabled/guestbook
sudo rm /etc/nginx/sites-enabled/default

sudo wget -q https://github.com/nikogura/guestbook/raw/master/vagrant/files/mountain-scene-welcome-sign-3.gif -O /usr/share/nginx/html/mountain-scene-welcome-sign-3.gif
sudo chmod 644 /usr/share/nginx/html/mountain-scene-welcome-sign-3.gif

sudo service nginx stop

sudo service nginx start


