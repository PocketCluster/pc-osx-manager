package main

import (
    "github.com/stkim1/pc-core/lib/swarm"
)

/*
swarm manage
--debug
--host =            :3376
--advertise=        pc-master:3376

--tlsverify =       true
--tlscacert =       /Users/almightykim/Workspace/DKIMG/CERT/ca-cert.pub
--tlscert =         /Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.cert
--tlskey =          /Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.key

nodes://192.168.1.150:2375,
192.168.1.151:2375,
192.168.1.152:2375,
192.168.1.153:2375,
192.168.1.161:2375,
192.168.1.162:2375,
192.168.1.163:2375,
192.168.1.164:2375,
192.168.1.165:2375,
192.168.1.166:2375
*/

func main() {
    context := swarm.NewContext(
        "0.0.0.0:3376",
        "192.168.1.150:2375,192.168.1.151:2375,192.168.1.152:2375,192.168.1.153:2375,192.168.1.161:2375,192.168.1.162:2375,192.168.1.163:2375,192.168.1.164:2375,192.168.1.165:2375,192.168.1.166:2375",
        "/Users/almightykim/Workspace/DKIMG/CERT/ca-cert.pub",
        "/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.cert",
        "/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.key",
    )
    context.Manage()
}
