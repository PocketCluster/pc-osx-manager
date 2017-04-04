// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package main

// Simple on-screen app debugging for OS X. Not an officially supported
// development target for apps, as screens with mice are very different
// than screens with touch panels.

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
    "github.com/stkim1/pc-core/event/lifecycle"
    "github.com/stkim1/pc-core/event/crash"
)

var initThreadID uint64

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
}

// this was app package of main()
func mainLifeCycle(f func(App)) int {
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
    theApp.sendLifecycle(lifecycle.StageDead)
}

//export lifecycleAlive
func lifecycleAlive() {
    theApp.sendLifecycle(lifecycle.StageAlive)
}

//export lifecycleVisible
func lifecycleVisible() {
    theApp.sendLifecycle(lifecycle.StageVisible)
}

//export lifecycleFocused
func lifecycleFocused() {
    theApp.sendLifecycle(lifecycle.StageFocused)
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
    theApp.sendCrash(crash.CrashEmergentExit)
}

//export engineDebugOutput
func engineDebugOutput(debug C.int) {
    if int(debug) == 0 {
        setLogger(false)
    } else {
        setLogger(true)
    }
}