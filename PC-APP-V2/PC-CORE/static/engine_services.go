package main

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/coreos/etcd/embed"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/beacon"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/extlib/registry"
    swarmemb "github.com/stkim1/pc-core/extlib/swarm"
    "github.com/stkim1/pc-core/service"
)

const (
    iventBeaconManagerSpawn string  = "ivent.beacon.manager.spawn"
)

func initStorageServie(a *appMainLife, config *embed.PocketConfig) error {
    a.RegisterServiceWithFuncs(
        operation.ServiceStorageProcess,
        func() error {
            etcd, err := embed.StartPocketEtcd(config)
            if err != nil {
                return errors.WithStack(err)
            }
            // startup preps
            select {
                case <-etcd.Server.ReadyNotify(): {
                    log.Debugf("[ETCD] server is ready to run")
                }
                case <-time.After(120 * time.Second): {
                    etcd.Server.Stop() // trigger a shutdown
                    return errors.Errorf("[ETCD] Server took too long to start!")
                }
            }
            // until server goes down, errors and stop signal will be constantly checked
            for {
                select {
                    case err = <-etcd.Err(): {
                        log.Debugf("[ETCD] error : %v", err)
                    }
                    case <- a.StopChannel(): {
                        etcd.Close()
                        log.Debugf("[ETCD] server shuts down")
                        return nil
                    }
                }
            }
            return nil
        })

    return nil
}

func initRegistryService(a *appMainLife, config *registry.PocketRegistryConfig) error {
    a.RegisterServiceWithFuncs(
        operation.ServiceContainerRegistry,
        func() error {
            reg, err := registry.NewPocketRegistry(config)
            if err != nil {
                return errors.WithStack(err)
            }
            err = reg.Start()
            if err != nil {
                return errors.WithStack(err)
            }
            log.Debugf("[REGISTRY] server start successful")

            // wait for service to stop
            <- a.StopChannel()
            err = reg.Stop(time.Second)
            log.Debugf("[REGISTRY] server exit. Error : %v", err)
            return errors.WithStack(err)
        })
    return nil
}

func initSwarmService(a *appMainLife) error {
    const (
        iventSwarmInstanceSpawn string  = "ivent.swarm.instance.spawn"
    )
    var (
        swarmSrvC = make(chan service.Event)
    )
    a.RegisterServiceWithFuncs(
        operation.ServiceSwarmEmbeddedOperation,
        func() error {
            var (
                swarmsrv *swarmemb.SwarmService = nil
            )
            for {
                select {
                    case se := <- swarmSrvC: {
                        srv, ok := se.Payload.(*swarmemb.SwarmService)
                        if ok {
                            swarmsrv = srv
                        }
                    }
                    case <- a.StopChannel(): {
                        if swarmsrv != nil {
                            err := swarmsrv.Close()
                            return errors.WithStack(err)
                        }
                        return errors.Errorf("[ERR] null SWARM instance")
                    }
                }
            }
            return nil
        },
        service.BindEventWithService(iventSwarmInstanceSpawn, swarmSrvC))

    beaconManC := make(chan service.Event)
    a.RegisterServiceWithFuncs(
        operation.ServiceSwarmEmbeddedServer,
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
            a.BroadcastEvent(service.Event{Name:iventSwarmInstanceSpawn, Payload:swarmsrv})
            err = swarmsrv.ListenAndServeSingleHost()
            return errors.WithStack(err)
        },
        service.BindEventWithService(iventBeaconManagerSpawn, beaconManC))

    return nil
}
