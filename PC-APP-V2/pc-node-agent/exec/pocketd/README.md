# Node `pocketd` install & config guide

### 1. `dhcpagent` is to receieve dhcp event.

1. place `dhcpagent` script in `/etc/dhcp/dhclient-exit-hooks.d` with following content for debugging

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
2. for production, remove all extra debugging info from `dhcpagent`

  ```sh
  # Notifies DHCP event
  # Copyright 2017 PocketCluster.io

  /opt/pocket/bin/pocketd dhcpagent
  ```
  - you might add ` > /dev/null 2>&1` at the end.
  - <https://unix.stackexchange.com/questions/119648/redirecting-to-dev-null>