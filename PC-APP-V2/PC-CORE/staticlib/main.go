package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_osxmain -framework Cocoa

extern int osxmain(int argc, const char * argv[]);
*/
import "C"
import (
    "net/http"
    "fmt"
    "time"
    "sync"

    "github.com/tylerb/graceful"
)

func RunWebServer(wg *sync.WaitGroup) *graceful.Server {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        fmt.Fprintf(w, "Welcome to the home page!")
    })
    srv := &graceful.Server{
        Timeout: 10 * time.Second,
        NoSignalHandling: true,
        Server: &http.Server{
            Addr: ":3001",
            Handler: mux,
        },
    }

    go func() {
        defer wg.Done()
        srv.ListenAndServe()
    }()
    return srv
}

func main() {

    var wg sync.WaitGroup
    srv := RunWebServer(&wg)

    go func() {
        wg.Wait()
    }()

    // Perhaps the first thing main() function needs to do is initiate OSX main
    C.osxmain(0, nil)

    srv.Stop(time.Second)
    fmt.Println("pc-core terminated!")
}