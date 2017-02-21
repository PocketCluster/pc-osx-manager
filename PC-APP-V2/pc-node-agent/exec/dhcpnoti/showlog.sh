#!/usr/bin/env bash

echo "----------------------------------- show journal log ------------------------------------"
/bin/journalctl -b0 --system _COMM=dhclient

echo "-------------------------------- show noti exit hook log --------------------------------"
cat /tmp/dh-client-env.log