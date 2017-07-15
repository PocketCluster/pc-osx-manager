package main

import (
    "sync"
    "time"
    "net"

    log "github.com/Sirupsen/logrus"
)

func echoServer(w *sync.WaitGroup) {
    defer w.Done()
    l, err := net.Listen("tcp4", "127.0.0.1:10068")
    if err != nil {
        log.Panicln(err)
    }
    log.Debugf("Listening to connections on port :10068")

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Debugf(err.Error())
            return
        }

        buf := make([]byte, 1024)

        for {
            size, err := conn.Read(buf)
            if err != nil {
                log.Debugf(err.Error())
                return
            }
            data := buf[:size]
            conn.Write(data)
        }
    }

}

func echoClient(w *sync.WaitGroup, done chan bool) {
    defer w.Done()

    var (
        count, errorCount, success int = 0, 0, 0
        buf []byte  = make([]byte, 10240)
        conn net.Conn = nil
        err error = nil
    )
    log.Debugf("[REPORTER] starting reporter service ...")

    for {
        select {
            case <- done: {
                return
            }
            default: {
                conn, err = net.DialTimeout("tcp4", net.JoinHostPort("127.0.0.1", "10068"), time.Second * 3)
                if err != nil {
                    log.Debugf("[REPORTER] connection error (%v)", err.Error())
                } else {
                    errorCount = 0
                    success = 0
                    //err = conn.SetDeadline(time.Now().Add(time.Second * time.Duration(3)))
                    if err != nil {
                        log.Debugf("[REPORTER] deadline error (%v)", err.Error())
                    } else {
                        for {
                            select {
                                case <- done: {
                                    conn.Close()
                                    return
                                }
                                default: {
                                    if 5 <= errorCount {
                                        conn.Close()
                                        conn = nil
                                        break
                                    }

                                    count, err = conn.Write([]byte("hello"))
                                    if err != nil {
                                        log.Debugf("[REPORTER] write error (%v)", err.Error())
                                        errorCount++
                                        continue
                                    }

                                    count, err = conn.Read(buf)
                                    if err != nil {
                                        log.Debugf("[REPORTER] read error (%v)", err.Error())
                                        errorCount++
                                        continue
                                    }
                                    success++
                                    log.Debugf("[REPORTER] All OK! %d | %s", success, string(buf[:count]))
                                    time.Sleep(time.Second * time.Duration(1))
                                }
                            }
                        }
                    }
                }
            }
        }
    }

}

func main() {
    done := make(chan bool)
    log.SetLevel(log.DebugLevel)
    var wg sync.WaitGroup
    wg.Add(2)
    go echoServer(&wg)
    go echoClient(&wg, done)
    wg.Wait()
}