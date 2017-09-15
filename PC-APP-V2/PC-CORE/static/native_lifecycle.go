// +build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -Wl,-U,_PCNativeThreadID,-U,_PCNativeMainStart,-U,_PCNativeMainStop

#include <pthread.h>
#include "PCLifeCycle.h"
#include "PCNativeThread.h"

*/
import "C"
import (
    "runtime"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/app"
    "github.com/stkim1/pc-core/config"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/event/lifecycle"
    "github.com/stkim1/pc-core/event/crash"
    "github.com/stkim1/pc-core/route"
)

type appMainLife struct {
    app.App
    route.Router
    service.ServiceSupervisor
}

var (
    initThreadID uint64
    theApp *appMainLife = nil
)

func init() {
    // Lock the goroutine responsible for initialization to an OS thread.
    // This means the goroutine running main (and calling PCNativeMainStart below)
    // is locked to the OS thread that started the program. This is
    // necessary for the correct delivery of Cocoa events to the process.
    //
    // A discussion on this topic:
    // https://groups.google.com/forum/#!msg/golang-nuts/IiWZ2hUuLDA/SNKYYZBelsYJ
    runtime.LockOSThread()
    initThreadID = uint64(C.PCNativeThreadID())
    theApp = &appMainLife{
        App: app.NewApp(),
        Router: route.NewRouter(func(_, _, _ string) error {
            return errors.Errorf("invalid path root (/)")
        }),
        ServiceSupervisor: service.NewServiceSupervisor(),
    }
}

// this was app package of main()
func appLifeCycle(f func(*appMainLife)) int {
    if tid := uint64(C.PCNativeThreadID()); tid != initThreadID {
        log.Fatalf("[CRITICAL] engine main called on thread %d, but inititaed from %d", tid, initThreadID)
    }

    go func() {
        f(theApp)
        C.PCNativeMainStop()
        // TODO(crawshaw): trigger PCNativeMainStart to return
    }()
    return int(C.PCNativeMainStart(0, nil))
}

//export lifecycleDead
func lifecycleDead() {
    theApp.SendLifecycle(lifecycle.StageDead)
}

//export lifecycleAlive
func lifecycleAlive() {
    theApp.SendLifecycle(lifecycle.StageAlive)
}

//export lifecycleVisible
func lifecycleVisible() {
    theApp.SendLifecycle(lifecycle.StageVisible)
}

//export lifecycleFocused
func lifecycleFocused() {
    theApp.SendLifecycle(lifecycle.StageFocused)
}

//export lifecycleAwaken
func lifecycleAwaken() {
    //TODO this is to be done later
    log.Debugf("lifecycleAwaken")
}

//export lifecycleSleep
func lifecycleSleep() {
    //TODO this is to be done later
    log.Debugf("lifecycleSleep")
}

//export crashEmergentExit
func crashEmergentExit() {
    theApp.SendCrash(crash.CrashEmergentExit)
}

//export engineDebugOutput
func engineDebugOutput(debug C.int) {
    if int(debug) == 0 {
        config.SetLogger(false)
    } else {
        config.SetLogger(true)
    }
}