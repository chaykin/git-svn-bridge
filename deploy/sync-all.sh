#!/bin/bash

cd /opt/git-bridge
./git-svn-bridge sync -r $1 -r $2 -r $3

