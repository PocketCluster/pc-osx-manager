package main

import (

    log "github.com/Sirupsen/logrus"
//    "github.com/pkg/errors"

//    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/udpnet/ucast"
)

func main() {
    agent := NewPocketSupervisor()
/*
    agent.RegisterFunc(func() error {
        caster, err := mcast.NewSearchCaster()
        if err != nil {
            return err
        }
        return caster.Close()
    })
*/
    agent.RegisterFunc(func() error {
        agent, err := ucast.NewBeaconAgent()
        if err != nil {
            return err
        } else {
            go func() {
                for v := range agent.ChRead {
                    log.Debugf("Received message %v", v.Message)
                }
            }()
        }
        return agent.Close()
    })

    err := agent.Start()
    if err != nil {
        log.Debug(err)
    }
    agent.Wait()
}