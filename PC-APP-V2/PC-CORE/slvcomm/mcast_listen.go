package slvcomm

import (
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/mcast"
)

func ServeMcastListenerOnWaitGroup(iface string, wg *sync.WaitGroup) error {
    var chErrors = make(chan error)

    go func (i string, w *sync.WaitGroup) {
        defer w.Done()
        listener, err := mcast.NewMcastListener(i)
        if err != nil {
            log.Error(err)
            chErrors <- errors.WithStack(err)
            return
        }

        chErrors <- nil
        for v := range listener.ChRead {
            log.Println(string(v.Message) + " : " + v.Address.IP.String())
        }
    }(iface, wg)

    return <- chErrors
}
