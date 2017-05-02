package slvcomm

import (
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/udpnet/ucast"
)

type BeaconSendPack struct {
    Message    []byte
    DstAddr    string
}

func ServeUcastBeaconOnWaitGroup(wg *sync.WaitGroup, chRead chan ucast.BeaconPack, chSend chan BeaconSendPack, chClose chan bool) error {
    var chErrors = make(chan error)
    wg.Add(1)
    go func () {
        defer wg.Done()

        channel, err := ucast.NewPocketLocatorChannel(wg)
        if err != nil {
            chErrors <- errors.WithStack(err)
            return
        }
        defer channel.Close()
        chErrors <- nil

        for {
            select {
                case <- chClose:
                    return

                case s := <-chSend:
                    err := channel.Send(s.DstAddr, s.Message)
                    if err != nil {
                        log.Debug(errors.WithStack(err))
                    }

                case r := <- channel.ChRead:
                    chRead <- r
            }
        }
    }()

    return <- chErrors
}

