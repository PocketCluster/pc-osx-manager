# LOG

**11/03/2016**

* DONE
  - `pc-node-agent/config` : Whole `yaml` config, `network/interfaces` fixation (`dhcp` -> `static`) and tests 
* TODO 
  - Complete test cases on `pc-node-agent/mcast`, `pc-node-agent/slagent`, `pc-node-agent/slcontext` 
  - `dns-server` is fixed to `pc-master:53535`. The needs `etc/hosts` to be edited.

**11/02/2016**

* DONE
  - `pc-node-agent/SlaveLocatingService` added  
  - Dependency table : `pc-node-agent/SlaveLocatingService` -> `locator` -> `slagent`  
* TODO 
  - Complete test cases on `pc-node-agent/config`, `pc-node-agent/mcast`, `pc-node-agent/slagent`, `pc-node-agent/slcontext` 
