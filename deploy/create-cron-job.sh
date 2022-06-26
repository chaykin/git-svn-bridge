#!/bin/bash

echo "*/10 * * * * /opt/git-bridge/sync-all.sh CORE MODULES TEST_PRJ" >> cron.tmp

crontab cron.tmp
rm cron.tmp

