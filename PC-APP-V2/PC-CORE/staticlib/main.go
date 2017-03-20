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

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/teleport/lib/utils"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/config"
    "github.com/stkim1/pcrypto"
)
import (
    "github.com/tylerb/graceful"
    "github.com/davecgh/go-spew/spew"
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

func setLogger(debug bool) {
    // debug setup
    if debug {
        utils.InitLoggerDebug()
        log.Info("DEBUG mode logger output configured")
    } else {
        utils.InitLoggerCLI()
        log.Info("NORMAL mode logger configured")
    }
}

func main() {
    setLogger(true)

    var wg sync.WaitGroup
    srv := RunWebServer(&wg)

    go func() {
        wg.Wait()
    }()

    // setup context
    ctx := context.SharedHostContext()
    context.FindSystemInterfaceStatus()
    context.FindSystemGatewayStatus()
    config.SetupBaseConfigPath(ctx)

    // open database
    db, err := config.OpenStorageInstance(ctx)
    if err != nil {
        log.Info(err)
    }
    // cert engine
    certStorage, err := pcrypto.NewPocketCertStorage(db)
    if err != nil {
        log.Info(err)
    }

    cfg := service.MakeCoreConfig(ctx, true)
    cfg.AssignCertStorage(certStorage)

    log.Info(spew.Sdump(ctx))

    // Perhaps the first thing main() function needs to do is initiate OSX main
    C.osxmain(0, nil)

    srv.Stop(time.Second)
    config.CloseStorageInstance(db)
    fmt.Println("pc-core terminated!")
}