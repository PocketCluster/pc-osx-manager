## TASKS

```
- TODO/ + HALF DONE/ * COMPLETED
```

## L.T.G

- [ ] MUSL ARM64 w/ GOLANG
- [ ] TinyCore ARM64 (from Busybox)
- [ ] Have a UI queue to stack ui requests so that duplicated requests will not cause double execution

### V0.1.4
- [ ] `health` check service should count *ETCD*, *Registry* service startup as well.
- [ ] make sure `interface` building tips from [`S39hostname`](/PC-APP-V1/pc-rpi-agent/S39hostname) being applied to [PocketAgent](/PC-APP-V1/pc-rpi-agent).
- [ ] Close health monitor + external net listener first before close any other services so that supervisor won't crash
- [ ] Fix External net listener open timing (after all services are ready)
- [ ] Fix Teleport generate log in Log directory
- [ ] `CoreNode`, `SlaveNode` : `LastAlive` data for health check on OSX. Save the last alive time to report in OSX.
- [ ] `CoreNode` set `SlaveID` with `MachineID`.
- [ ] Cluster User Defined Name
- [ ] Repartition logic from Boot2Docker
- [ ] Disable Docker when there is no TLS certs available
- [ ] Append 'CA' <http://rt.openssl.org/Ticket/History.html?user=guest&pass=guest&id=3979>
- [ ] Run `dockerd` with script and check conditions there like boot2docker
- [ ] Combine Docker related stuff in a unified package
- [ ] Combine `pc-core` related consts, variable in a unified package
- [ ] Docker User Namespace Remap
  * <https://docs.oracle.com/cd/E37670_01/E75728/html/ol-docker-userns-remap.html>
  * <https://docs.docker.com/edge/engine/reference/commandline/dockerd/#daemon-user-namespace-options>
- [ ] Singleton lock for it’s property access
- [ ] UUID for ID (Too long. we will do it when udp packet fragmentation is supported)
- [ ] Complete Slave node `bounded` state with valid checks
- [ ] Instant status check on `bounded` (i.e. as soon as master pings slave, master gets response)
- [ ] `Makefile` to remove all *_test.go and TEST.go
- [ ] Clean up core node including `vi`
- [ ] More Test for `pc-vbox-comm`
- [ ] Pass `uid` + `gid` to all slave nodes to have them setup the __user__ with the same `uid` 
  * Look `Resources.bundle/setup/raspberry_user_setup.sh` and extract essential tasks
  * UID/GID pair in vbox disk for core. Same pair in slavenode Teleport docker setup
- [ ] Look carefully `etcd` TLS configuration behave with `docker-compose` connected.
  * At initial, etcd `network/bridged` API point isn't available and dockerd complains. This seems to be normal though.
  * <https://dims-devguide.readthedocs.io/en/latest/dockerdatacenter.html>
- [ ] **Slave** : Remove unnecessary `const` in `slagent` package such as `SLAVE_CLUSTER_MEMBERS`
- [ ] **Slave** : Remove unnecessary `field` in `PocketSlaveDiscovery` & `PocketSlaveDiscovery` struct
- [ ] **Slave** : search logic in master `Beacon` package
- [ ] **Slave** : testing on timeout mechanism for receiving master meta
- [ ] **Slave** <-> Master `timezone` incompatibility
- [ ] **Slave** : config should be able to tell if Slave node is bounded or not by reading config
- [ ] **Slave** : After changing hostname, please updata `/etc/hosts`
- [ ] _Master_ interface refresh logic
- [ ] _Master_ SQLite encryption
- [ ] _Master_ Remove unnecessary `const` & `field` in `msagent` package
- [ ] [_RPI_] cannot acquire proper interface name from netifaces for default gateway
- [x] Provide `pc-master` DNS for `docker-compose`
- [x] [_ODROID_] netmask format fix for network/interfaces (`ffffff00` -> `255.255.255.0`)
- [x] Vagrant Fix Network Interface Order (eth0 : internal/ eth1 : external)
- [x] Move Docker port to `0.0.0.0:2376` for ssl connection
- [x] AESKEY regeneration when `MasterBeacon` goes to `BindBroken` state
- [x] Shorten `msgpack` name field to reduce message package size
- [x] Network Broadcast Form -> 192.168.2.211/24
- [x] _Master_ Private/Public/CA
- [x] _Master_ teleport bolts -> SQLite
- [-] Slave sends SSH key in cryptocheck to keyexchange of master
- [-] [_PINE64_] `fdisk` new partition incorrectly begins new partition sector from 2048
- [-] Use FQDN name for certification and name server ( e.g. `pc-core.Q1oqc1921lq.cluster.pocketcluster.io` )
  * `pc-master` queries to main DNS.


#### UI Layer
- [ ] check `pc-core` & `pc-nodes` connections for
  1. teleport clients
  2. swarm nodes
  3. beacon/locator
- [ ] UI mode flow
  0. sign up
  1. setup cluster / restart cluster
  2. install a package
  3. run a package
  4. turn off the package
  5. turn off cluster

#### VirtualBox
- [x] get the actual core count
- [x] vbox files permissions
- [x] multiple shared folder
- [x] vbox machine start stop
- [x] force get session from running machine to reboot
  * -> We cannot hijack a session. We force to stop and restart with command line
- [x] check if setting has been changed
- [x] make sure ACPI shutdown works properly
- [x] discard machine setting if change has made (IMachine::discardSettings())
  * -> when force machine powerdown settings changed! 
- [ ] shared folder permissions
- [ ] modify machine settings
- [ ] disable snapshot, save, teleport,  restore, fault-tolerant-sync, setting up, (these are illiegal statates)
- [ ] vbox disk creation from XPCOM
- [ ] vbox machine “options”
  * "VBoxInternal/CPUM/EnableHVP", “1"


#### Beacon/Locator
- [ ] Cleanup the package / Reduce Garbage collection (no allocation on every socket, clean interface, ...)
- [ ] Per-node PUB/PRI key pair on Master
- [ ] Additional Command (UPDATE/REBOOT/SHUTDOWN/BREAKCONN)
- [ ] Expiring discovery broadcast after some point
- [ ] TLS-TCP connection for confirmed and large packet size


### V0.1.3
-
```
* root permission to create/edit/copy config files  
* launching service (Webserver, Salt, Multicast)
* copy & modify /etc/hosts, /etc/salt/minion, /etc/salt/master
* create /pocket, /bigpkg
* unify linux version to truty64
* Add PC-MASTER key to known_hosts
* install salt during setup process
* CHECK if brew, java are installed
* Move monitoring activation at the end of package installation
* backup original salt configcd /
* remove all ssh files when all done
* before package installation complete following steps
	- basic file, timezone, locale setup
	- basic ssh login credential
* Termination script
* CHECK if local ssh is open

* [RPI] JVM installation
* [RPI] grap ip4 addresses & select ethernet interface
* [RPI] Create special package directory : /bigpkg/archive
* [RPI] store bigpkgs in /bigpkg/archive
* [RPI] unify node name to pc-node[?]
* [RPI] time broadcast from OSX
* [RPI] timezone broadcast from OSX
* [RPI] check if cluster is online then execute the next move.
* [RPI] install process check if it is neccessary to download bigpkg
* [RPI] check /pocket, /bigpkg directory
* [RPI] check repartition!
* [RPI] add pc-node4 ~ 6 to slaves file1
* [RPI] check if swap file created.

* [VAG] CHECK if vagrant, virtualbox is installed
* [VAG] fix virtualbox setup process
	1 install vbox interface (load vbox environment!)
	2 install base config with sudo priviledeges
	3 install salt
	4 install vagrant instances
* [VAG] create vbox interface & check its ip addresss
* [VAG] while starting, fix vagrant interface according to the setup process

- [RPI] detailed installation process on progress label
- [VAG] detailed private_network interface building -> private_network, name:vboxnet2, ip:10.211.55.201
- [VAG] remove redandunt files (apt-get autoremove, clean)
- [VAG] install language package (apt-get locale)
- Finish button to remove window
- change /etc/hosts of every datanode
- [RPI] send cluster member ip only when firstly setting up, or member changes
- [VAG] vagrant box update, 
- Until all stopping process is completed, do not show complete indicator on menu
- Vagrant reload/reboot


- Install multiple packages
- Disable close button on installation
- salt 'pc-master' pkg.install R taps='homebrew/science' does not work
- separate pre-requisite library installation
- ability to skip first installation part 
- spark directory setups
- dependedn project needs to start first
- spark slave execution script fix
- spark/hadoop needs to be re-configured to run secondary master in slaves.sh scripts
	- this script prevent a binary package to be reused
- [VAG] should be able to restart when it fails to boot up everyone
- Update package menu when a new package is installed
- Check hdfs sparkJob directory creation
- when job is still active stop proceed
- check download size and if that does not match, re-try


- SPARK: change work dir path (user.dir)
- SPARK: change metadata_db path 
* SPARK: chceck jps IF HADOOP PROCESS IS ALIVE
* SPARK: install scala 2.11.7 to slave nodes -> what if different version of same package installed and collide. for example scala 2.11.7 && scala 3.11.6?

- [VAG] need UUID
- [RPI-OSX] async network status change alert
	- change host file
- register port to system automatically
- [RPI] Add raspberry heartbeat checker, 
-(?) remove localhost from known_hosts
- relaunching service (From NSLogger)
- Download file and run checksum
- [RPI] command type

===============================================================================

- Hadoop with Secondary Node
- Hadoop Setup completely
- Spark Install

===============================================================================

- Let people clear XCODE lisence	
```
