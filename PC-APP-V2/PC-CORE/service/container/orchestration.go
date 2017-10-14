package container

import (
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/beacon"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/operation"
    swarmemb "github.com/stkim1/pc-core/extlib/swarm"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
)

// nost status info for health monitor
type EngineStatusInfo struct {
    Name    string
    ID      string
    Addr    string
}

func InitSwarmService(appLife service.ServiceSupervisor) error {
    const (
        iventSwarmInstanceSpawn string  = "ivent.swarm.instance.spawn"
    )
    var (
        swarmSrvC = make(chan service.Event)
        nodeStatC = make(chan service.Event)
    )
    appLife.RegisterServiceWithFuncs(
        operation.ServiceOrchestrationOperation,
        func() error {
            var (
                swarmsrv *swarmemb.SwarmService = nil
            )
            for {
                select {
                    case se := <- swarmSrvC: {
                        sobj, ok := se.Payload.(*swarmemb.SwarmService)
                        if ok {
                            log.Debugf("[SWARM-CTRL] orchestration instance detected...")
                            swarmsrv = sobj
                        } else {
                            log.Errorf("[SWARM-CTRL] unable to recieve orchestration instance")
                        }
                    }
                    case <- appLife.StopChannel(): {
                        if swarmsrv != nil {
                            err := swarmsrv.Close()
                            return errors.WithStack(err)
                        }
                        return errors.Errorf("[ERR] null SWARM instance")
                    }
                    case <- nodeStatC: {
                        if swarmsrv == nil {
                            appLife.BroadcastEvent(service.Event{
                                Name:    ivent.IventMonitorNodeRsltOrchst,
                                Payload: errors.Errorf("unable to query orchestration engine")})
                            continue
                        }

                        nets := swarmsrv.Cluster.Networks()
                        engines := make([]EngineStatusInfo, 0, len(nets))
                        for i, _ := range nets {
                            n := nets[i]
                            log.Debugf("node ip %v | addr %v", n.Engine.Addr)
                            engines = append(engines, EngineStatusInfo {
                                Name: n.Engine.Name,
                                ID:   n.Engine.ID,
                                Addr: n.Engine.IP,
                            })
                        }
                        appLife.BroadcastEvent(service.Event{
                            Name:    ivent.IventMonitorNodeRsltOrchst,
                            Payload: engines})
                    }
                }
            }
            return nil
        },
        service.BindEventWithService(iventSwarmInstanceSpawn,         swarmSrvC),
        service.BindEventWithService(ivent.IventMonitorNodeReqOrchst, nodeStatC))

    beaconManC := make(chan service.Event)
    appLife.RegisterServiceWithFuncs(
        operation.ServiceOrchestrationServer,
        func() error {
            be := <- beaconManC
            beaconMan, ok := be.Payload.(beacon.BeaconManger)
            if !ok {
                return errors.Errorf("[ERR] invalid beacon manager type")
            }
            ctx := context.SharedHostContext()
            caCert, err := ctx.CertAuthCertificate()
            if err != nil {
                return errors.WithStack(err)
            }
            hostCrt, err := ctx.MasterHostCertificate()
            if err != nil {
                return errors.WithStack(err)
            }
            hostPrv, err := ctx.MasterHostPrivateKey()
            if err != nil {
                return errors.WithStack(err)
            }
            swarmctx, err := swarmemb.NewContextWithCertAndKey(caCert, hostCrt, hostPrv, beaconMan)
            if err != nil {
                return errors.WithStack(err)
            }
            swarmsrv, err := swarmemb.NewSwarmService(swarmctx)
            if err != nil {
                return errors.WithStack(err)
            }
            appLife.BroadcastEvent(service.Event{Name:iventSwarmInstanceSpawn, Payload:swarmsrv})
            log.Debugf("[SWARM] swarm service started...")
            err = swarmsrv.ListenAndServeSingleHost()
            return errors.WithStack(err)
        },
        service.BindEventWithService(ivent.IventBeaconManagerSpawn, beaconManC))

    return nil
}
