# Node `pocketd` install & config guide

### `dhcpagent` is to receieve dhcp event.

1. [DEBUG] place a script in `/etc/dhcp/dhclient-exit-hooks.d/dhcpagent` with following content for debugging

  ```sh
  # Notifies DHCP event
  # Copyright 2017 PocketCluster.io

  echo "$(date): entering ${1%/*}, dumping variables." >> "/tmp/dh-client-env.log"
  /opt/pocket/bin/dhcpevent -mode=dhcpagent -dev=jsonprint | python -mjson.tool >>  "/tmp/dh-client-env.log"
  ```
  - to view the log, do the following

  ```sh
  #!/usr/bin/env bash

  echo "----------------------------------- show journal log ------------------------------------"
  /bin/journalctl -b0 --system _COMM=dhclient

  echo "-------------------------------- show noti exit hook log --------------------------------"
  cat /tmp/dh-client-env.log
  ```
2. [RELEASE] for production, remove all extra debugging info from `dhcpagent` with `0644` permission

  ```sh
  # Notifies DHCP event
  # Copyright 2017 PocketCluster.io

  /opt/pocket/bin/pocketd dhcpagent
  ```
  - you might add ` > /dev/null 2>&1` at the end.
  - <https://unix.stackexchange.com/questions/119648/redirecting-to-dev-null>

### `systemd service`

1. `pocketd.service`

  ```sh
  [Unit]
  Description=PocketCluster Node Agent
  After=network.target
  
  [Service]
  Type=simple
  PIDFile=/var/run/pocketd.pid
  Restart=always
  ExecStart=/opt/pocket/bin/pocketd
  
  [Install]
  WantedBy=multi-user.target
  ```
2. Activate the service with foloowing command

  ```sh
  mv pocketd.service /etc/systemd/system/ && chown root:root /etc/systemd/system/pocketd.service
  
  systemctl daemon-reload
  systemctl start pocketd
  systemctl enable pocketd
  systemctl status pocketd.service
  ```

### Update process

1. Copy `update.sh` -> `/opt/pocket/bin/update.sh` and `pocketd.update` -> `/opt/pocket/bin/pocketd.update` and `/bin/bash /opt/pocket/bin/update.sh`

  - `update.sh` (update script)

  ```sh
  #!/usr/bin/env bash

  if [[ -f "/opt/pocket/bin/pocketd.update" ]]; then
    mv /opt/pocket/bin/pocketd.update /opt/pocket/bin/pocketd
  fi

  if [[ -f "/opt/pocket/bin/update.sh" ]]; then
    rm /opt/pocket/bin/update.sh
  fi

  /usr/sbin/service pocket restart
  ```
