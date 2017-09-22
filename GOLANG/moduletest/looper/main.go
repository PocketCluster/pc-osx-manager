package main

import (
    "time"
    "sync"
    "sync/atomic"
    "runtime"

    log "github.com/Sirupsen/logrus"
    //"github.com/gravitational/trace"
    "io/ioutil"
)

func updateSupplement(waiter *sync.WaitGroup, isWorking *atomic.Value) {
    waiter.Add(1)
    isWorking.Store(true)
    log.Infof("[updateSupplement] update cycle begin")
    time.Sleep(time.Second * 10)
    log.Infof("[updateSupplement] update cycle done")
    isWorking.Store(false)
    waiter.Done()
}

func main() {

    if false {
        ioutil.WriteFile("timestamp.txt", []byte(time.Now().Format(time.RFC3339)), 0600)
    } else {
        ts, err := ioutil.ReadFile("timestamp.txt")
        if err != nil {
            log.Error(err.Error())
        } else {
            lastRec, err := time.Parse(time.RFC3339, string(ts))
            if err == nil {
                log.Info(lastRec.String())
            }
        }
    }
    return


    runtime.GOMAXPROCS(2)
    var (
        trigger     = time.NewTicker(time.Second * 2)
        mainTicker  = time.NewTicker(time.Second)
        abort       = make(chan bool)


        waiter sync.WaitGroup
        abortCounter int = 0
    )

    go func (timer <- chan time.Time, quit <- chan bool, wg *sync.WaitGroup) {
        log.Info("[updateSupplement] let's tick a notch!")
        wg.Add(1)
        defer wg.Done()

        var isWorking atomic.Value
        isWorking.Store(false)

        for {
            select {
            case launch := <- timer:
                if isWorking.Load().(bool) {
                    log.Infof("[updateSupplement] %v inprogress", launch)
                } else {
                    go updateSupplement(wg, &isWorking)
                }

            case <- quit:
                log.Info("[updateSupplement] time to quit...")
                return
            }
        }
    }(trigger.C, abort, &waiter)

    mainLooper:
    for {
        select {
        case <- mainTicker.C:
            abortCounter++
            log.Infof("[main] Count : %v", abortCounter)
            if (60 < abortCounter) {
                abort <- true
                mainTicker.Stop()
                break mainLooper
            }
        }
    }

    log.Info("[main] We've broken out of main loop! Let's wait!!!")
    waiter.Wait()
    trigger.Stop()
    close(abort)
    log.Info("[main] Now all channels are closed!!!")
}