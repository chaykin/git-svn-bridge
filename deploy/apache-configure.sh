#!/bin/bash

SCRIPT_DIR="$( dirname "$0" )"

sudo cp $SCRIPT_DIR/git-bridge.conf /etc/apache2/sites-available/git-bridge.conf

sudo a2dissite 000-default.conf
sudo a2ensite git-bridge.conf

sudo sed -i 's/APACHE_RUN_USER=www-data/APACHE_RUN_USER=bridge-admin/' /etc/apache2/envvars
sudo sed -i 's/APACHE_RUN_GROUP=www-data/APACHE_RUN_GROUP=bridge-admin/' /etc/apache2/envvars
sudo chown -R bridge-admin:bridge-admin /var/www/html/

sudo systemctl restart apache2
