#!/usr/bin/env bash

echo "upgrading system"
sudo bash -c "DEBIAN_FRONTEND=noninteractive apt-get update && DEBIAN_FRONTEND=noninteractive sudo apt-get upgrade -y"

echo "Installing Git"
sudo bash -c "DEBIAN_FRONTEND=noninteractive apt-get install -y nginx"


NGINX=$(cat <<'EOF'
server {
    listen 80 default_server;
    server_name localhost;

    root /usr/share/nginx/html;
    index index.html index.htm;

    location /guestbook {
        proxy_pass http://backend:8080/guestbook;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}

EOF
)

echo -e "$NGINX" | sudo tee /etc/nginx/sites-available/guestbook

sudo ln -s /etc/nginx/sites-available/guestbook /etc/nginx/sites-enabled/guestbook
sudo rm /etc/nginx/sites-enabled/default

sudo mv /tmp/mountain-scene-welcome-sign-3.gif /usr/share/nginx/html/mountain-scene-welcome-sign-3.gif
sudo chmod 644 /usr/share/nginx/html/mountain-scene-welcome-sign-3.gif

sudo service nginx restart


