// +build darwin
package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_PCEventFeedGet,-U,_PCEventFeedPost,-U,_PCEventFeedPut,-U,_PCEventFeedDelete

#include "PCEventHandle.h"
*/
import "C"
import (
    "runtime"

    "github.com/pkg/errors"
)

type feedMethodType int
const (
    feedMethodGet       feedMethodType = iota
    feedMethodPost
    feedMethodPut
    feedMethodDelete
)

type feedMessage struct {
    method    feedMethodType
    path      string
    payload   string
}

type stopFeed struct{}

type feeder struct {
    feedPipe    chan interface{}
}

// loop is to feed engine message to Cocoa side
//
// After Cocoa has captured the initial OS thread for processing Cocoa
// events in runApp, it starts loop on another goroutine. It is locked
// to an OS thread for delivering message to the same thread
func (h *feeder) feedLoop() {
    runtime.LockOSThread()
    for f := range h.feedPipe {
        switch feed := f.(type) {
            // stop feeding back to cocoa side
            case stopFeed: {
                runtime.UnlockOSThread()
                return
            }
            case feedMessage: {
                switch feed.method {
                    case feedMethodGet: {
                        var (
                            cPath = C.CString(feed.path)
                        )
                        C.PCEventFeedGet(cPath)
                        // these strings will be freed on cocoa side to reduce performance degrade
                        //C.free(unsafe.Pointer(cPath))
                    }
                    case feedMethodPost: {
                        var (
                            cPath    = C.CString(feed.path)
                            cPayload = C.CString(feed.payload)
                        )
                        C.PCEventFeedPost(cPath, cPayload)
                        // these strings will be freed on cocoa side to reduce performance degrade
                        //C.free(unsafe.Pointer(cPath))
                        //C.free(unsafe.Pointer(cPayload))
                    }
                }
            }
        }
    }
}

var theFeeder = &feeder{
    feedPipe: make(chan interface{}),
}

//export FeedStart
func FeedStart() {
    go theFeeder.feedLoop()
}

//export FeedStop
func FeedStop() {
    theFeeder.feedPipe <- stopFeed{}
}

func EventFeedGet(path string) error {
    if len(path) == 0 {
        return errors.Errorf("[ERR] invalid feed path")
    }
    theFeeder.feedPipe <- feedMessage {
        method:     feedMethodGet,
        path:       path,
    }
    return nil
}

func EventFeedPost(path, payload string) error {
    if len(path) == 0 {
        return errors.Errorf("[ERR] invalid feed path")
    }
    theFeeder.feedPipe <- feedMessage {
        method:     feedMethodPost,
        path:       path,
        payload:    payload,
    }
    return nil
}
