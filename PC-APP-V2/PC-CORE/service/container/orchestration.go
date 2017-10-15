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
    "time"
)

func InitSwarmService(appLife service.ServiceSupervisor) error {
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
                        if sobj != nil && ok {
                            log.Debugf("[ORCHST.CTRL] orchestration instance detected...")
                            swarmsrv = sobj
                        } else {
                            log.Errorf("[ORCHST.CTRL] unable to recieve orchestration instance")
                        }
                    }
                    case <- appLife.StopChannel(): {
                        if swarmsrv != nil {
                            err := swarmsrv.Close()
                            return errors.WithStack(err)
                        }
                        return errors.Errorf("[ERR] null orchestration instance")
                    }
                    case re := <- nodeStatC: {
                        ts, ok := re.Payload.(int64)
                        if !ok {
                            appLife.BroadcastEvent(service.Event{
                                Name:    ivent.IventMonitorNodeRespOrchst,
                                Payload: ivent.EngineStatusMeta{
                                    TimeStamp: ts,
                                    Error:     errors.Errorf("inaccurate timestamp"),
                                }})
                            continue
                        }
                        if swarmsrv == nil {
                            appLife.BroadcastEvent(service.Event{
                                Name:    ivent.IventMonitorNodeRespOrchst,
                                Payload: ivent.EngineStatusMeta{
                                    TimeStamp: ts,
                                    Error:     errors.Errorf("unable to query orchestration engine"),
                                }})
                            continue
                        }

                        nets := swarmsrv.Cluster.Networks()
                        engines := make([]ivent.EngineStatusInfo, 0, len(nets))
                        for i, _ := range nets {
                            n := nets[i]
                            engines = append(engines, ivent.EngineStatusInfo {
                                Name: n.Engine.Name,
                                ID:   n.Engine.ID,
                                IP:   n.Engine.IP,
                                Addr: n.Engine.Addr,
                            })
                        }
                        appLife.BroadcastEvent(service.Event{
                            Name:    ivent.IventMonitorNodeRespOrchst,
                            Payload: ivent.EngineStatusMeta{
                                TimeStamp: ts,
                                Engines:   engines,
                            }})
                    }
                }
            }
            return nil
        },
        service.BindEventWithService(ivent.IventOrchstInstanceSpawn,  swarmSrvC),
        service.BindEventWithService(ivent.IventMonitorNodeReqStatus, nodeStatC))

    beaconManC := make(chan service.Event)
    appLife.RegisterServiceWithFuncs(
        operation.ServiceOrchestrationServer,
        func() error {
            var (
                beaconMan beacon.BeaconManger = nil
                failtimout = time.NewTicker(time.Minute)
            )

            // wait becon agent
            select {
                case <- failtimout.C: {
                    failtimout.Stop()
                    return errors.Errorf("[ORCHST] unable to recieve beacon manager")
                }
                case be := <- beaconManC: {
                    bm, ok := be.Payload.(beacon.BeaconManger)
                    if bm != nil && ok {
                        beaconMan = bm
                    } else {
                        return errors.Errorf("[ERR] invalid beacon manager type")
                    }
                }
            }

            failtimout.Stop()
            log.Info("[ORCHST] beacon manager received...")

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
            appLife.BroadcastEvent(service.Event{
                Name:ivent.IventOrchstInstanceSpawn,
                Payload:swarmsrv})

            log.Debugf("[ORCHST] orchestration service started...")
            err = swarmsrv.ListenAndServeSingleHost()
            return errors.WithStack(err)
        },
        service.BindEventWithService(ivent.IventBeaconManagerSpawn, beaconManC))

    return nil
}
