#!/bin/bash

SCRIPT_DIR="$( dirname "$0" )"
source $SCRIPT_DIR/create-repo.sh

sudo apt update
sudo apt install cron git git-svn subversion apache2 apache2-utils -y
sudo a2enmod env cgi alias rewrite

git config --global user.name "Git Svn Bridge"
git config --global user.email "example@company.ru"

sudo mkdir /opt/git-bridge
sudo cp $SCRIPT_DIR/git-svn-bridge /opt/git-bridge
sudo cp $SCRIPT_DIR/config.yml /opt/git-bridge
sudo cp $SCRIPT_DIR/sync-all.sh /opt/git-bridge
sudo cp -r $SCRIPT_DIR/gitHookTemplates /opt/git-bridge
sudo chown -R $USER /opt/git-bridge

pushd /opt/git-bridge
	createRepo "CORE" "PRJS/CORE"
	STATUS=$?
	if [ $STATUS -ne 0 ]; then
		echo "Could not create repo: CORE"
		exit $STATUS
	fi

	createRepo "MODULES" "PRJS/MODULES"
	STATUS=$?
	if [ $STATUS -ne 0 ]; then
		echo "Could not create repo: MODULES"
		exit $STATUS
	fi

	createRepo "TEST_PRJ" "PRJS/TEST_PRJ"
	STATUS=$?
	if [ $STATUS -ne 0 ]; then
		echo "Could not create repo: TEST_PRJ"
		exit $STATUS
	fi
popd

htpasswd -b -c /opt/git-bridge/.htpasswd test_user test_password

sudo chown -R $USER /opt/git-bridge
sudo chgrp -R www-data /opt/git-bridge
sudo chmod -R 770 /opt/git-bridge

sudo git config --global --add safe.directory "/opt/git-bridge/repos/git/CORE"
sudo git config --global --add safe.directory "/opt/git-bridge/repos/bridge/CORE"
sudo git config --global --add safe.directory "/opt/git-bridge/repos/git/MODULES"
sudo git config --global --add safe.directory "/opt/git-bridge/repos/bridge/MODULES"
sudo git config --global --add safe.directory "/opt/git-bridge/repos/git/TEST_PRJ"
sudo git config --global --add safe.directory "/opt/git-bridge/repos/bridge/TEST_PRJ"

$SCRIPT_DIR/apache-configure.sh
$SCRIPT_DIR/create-cron-job.sh
echo "======================== DONE ========================"
