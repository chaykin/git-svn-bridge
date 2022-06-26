#!/bin/bash

createRepo() {
	REPO_NAME=$1
	REPO_PATH=$2
	
	./git-svn-bridge create "https://localhost/$REPO_PATH" --name $REPO_NAME
	
	{ echo "1" & sleep 0.1; echo "1" & sleep 0.1; echo "example@company.ru" & sleep 0.1; echo "Git Svn Bridge"; } | ./git-svn-bridge add-user $REPO_NAME "Git Svn Bridge"

	COUNT=0
	while true
	do
		if (($COUNT > 0)); then
			echo "Continue ($COUNT) init $REPO_NAME"
			./git-svn-bridge init $REPO_NAME --continue
			STATUS=$?
			echo "\r\nContinue init $REPO_NAME result: $STATUS"
		else
			echo "Init $REPO_NAME first time"
			./git-svn-bridge init $REPO_NAME
			STATUS=$?
			echo "\r\nInit $REPO_NAME first time result: $STATUS"
		fi

		if [ $STATUS -ne 0 ]; then
			((COUNT++))
			if (($COUNT >= 10)); then
				echo "Too many tries!"
				return $STATUS
			fi
			sleep 1
		else
			break
		fi
	done
	echo "\r\nSuccessfully inited $REPO_NAME" 
	
	echo "Configure repo $REPO_NAME" 
	pushd "/opt/git-bridge/repos/git/$REPO_NAME"
	git config http.receivepack true
	git update-server-info
	chown -Rf www-data:www-data "/opt/git-bridge/repos/git/$REPO_NAME"
	popd
	
	echo "Successfully created $REPO_NAME" 
	return 0
}
