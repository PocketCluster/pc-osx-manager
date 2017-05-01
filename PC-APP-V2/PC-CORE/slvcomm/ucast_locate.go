package slvcomm

import (
    "strconv"
    "sync"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/ucast"
)

func ServeUcastLocationOnWaitGroup(wg *sync.WaitGroup) error {
    var chErrors = make(chan error)
    go func (w *sync.WaitGroup) {
        defer w.Done()

        channel, err := ucast.NewPocketLocatorChannel()
        if err != nil {
            chErrors <- errors.WithStack(err)
            return
        }
        defer channel.Close()
        chErrors <- nil

        for i := 0; i < 10; i++ {
            err := channel.Send("192.168.1.220", []byte("HELLO! - " + strconv.Itoa(i)))
            if err != nil {
                log.Fatal(err.Error())
            }
        }
        time.After(time.Millisecond)
    }(wg)

    return <- chErrors
}

