package beacon

import (
    "time"

    "github.com/docker/docker/pkg/discovery"
    . "gopkg.in/check.v1"

    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/model"
)

func (s *ManagerSuite) TestWatch(c *C) {
    var (
        _ = insertTestNodes(1, c)
        comm = &DebugCommChannel{}
        noti = &DebugBeaconNotiReceiver{}
        masterTS, slaveTS time.Time = time.Now(), time.Now()
        man, err = NewBeaconManager(masterAgentName, noti, comm)
    )
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].SlaveNode(), NotNil)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterBindBroken)

    piface, err := slcontext.PrimaryNetworkInterface()
    c.Assert(err, IsNil)

    meta := &slagent.PocketSlaveAgentMeta{
        MetaVersion:         slagent.SLAVE_META_VERSION,
        MasterBoundAgent:    masterAgentName,
        SlaveID:             "00:00:00:00:00:00",
        DiscoveryAgent:      &slagent.PocketSlaveDiscovery {
            Version:             slagent.SLAVE_DISCOVER_VERSION,
            SlaveResponse:       slagent.SLAVE_LOOKUP_AGENT,
            SlaveAddress:        piface.PrimaryIP4Addr(),
            SlaveGateway:        piface.GatewayAddr,
        },
    }
    psm, err := slagent.PackedSlaveMeta(meta)
    c.Assert(err, IsNil)

    // check if this successfully generate new beacon and move the transition
    err = man.TransitionWithSearchData(mcast.CastPack{Address:*slaveAddr, Message:psm}, masterTS)
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].SlaveNode(), NotNil)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterBindRecovery)

    // prepare slave node data
    slaveTS = masterTS.Add(time.Second)
    mb := man.(*beaconManger).beaconList[0]
    aescryptor := mb.(*masterBeacon).state.(DebugState).AESCryptor()

    // create slave meta
    sa, err := slagent.SlaveBoundedStatus("pc-node1", mb.SlaveNode().SlaveUUID, slaveTS)
    c.Assert(err, IsNil)
    mp, err := slagent.PackedSlaveStatus(sa)
    c.Assert(err, IsNil)
    encrypted, err := aescryptor.EncryptByAES(mp)
    c.Assert(err, IsNil)
    ma := &slagent.PocketSlaveAgentMeta{
        MetaVersion:         slagent.SLAVE_META_VERSION,
        MasterBoundAgent:    masterAgentName,
        SlaveID:             "00:00:00:00:00:00",
        EncryptedStatus:     encrypted,
    }
    psa, err := slagent.PackedSlaveMeta(ma)
    c.Assert(err, IsNil)

    // test slave meta
    masterTS = slaveTS.Add(time.Second)
    err = man.TransitionWithBeaconData(ucast.BeaconPack{Address:*slaveAddr, Message:psa}, masterTS)
    c.Assert(err, IsNil)
    c.Assert(len(man.(*beaconManger).beaconList), Equals, 1)
    c.Assert(man.(*beaconManger).beaconList[0].CurrentState(), Equals, MasterBounded)


    man.Initialize("foo", 1000, 0, nil)
    stopCh := make(chan struct{})
    ch, errCh := man.Watch(stopCh)

    // We have to drain the error channel otherwise Watch will get stuck.
    go func() {
        for range errCh {
        }
    }()

    addr, err := model.IP4AddrToString(piface.PrimaryIP4Addr())
    c.Assert(err, IsNil)

    expected := discovery.Entries{
        &discovery.Entry{Host: addr, Port: "2375"},
    }
    obtained := <-ch
    c.Assert(obtained.Equals(expected), Equals, true)

    // Stop and make sure it closes all channels.
    close(stopCh)
    c.Assert(<-ch, IsNil)
    c.Assert(<-errCh, IsNil)
}
