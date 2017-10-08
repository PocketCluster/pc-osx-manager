#!/usr/bin/env bash

if [[ -f "/opt/pocket/bin/pocketd.update" ]]; then
  mv /opt/pocket/bin/pocketd.update /opt/pocket/bin/pocketd
fi

if [[ -f "/opt/pocket/bin/update.sh" ]]; then
  rm /opt/pocket/bin/update.sh
fi

/usr/sbin/service pocket restart
