// +build darwin
package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_PCFeedResponseForGet,-U,_PCFeedResponseForPost,-U,_PCFeedResponseForPut,-U,_PCFeedResponseForDelete

#include "PCResponseHandle.h"
*/
import "C"
import (
    "runtime"

    "github.com/stkim1/pc-core/route"
    "github.com/pkg/errors"
)

type feedResponse struct {
    method      string
    path        string
    response    string
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
            case feedResponse: {
                var (
                    cPath     *C.char = C.CString(feed.path)
                    cResponse *C.char = C.CString(feed.response)
                    // these strings will be freed on cocoa side to reduce performance degrade
                    //C.free(unsafe.Pointer(cPath))
                    //C.free(unsafe.Pointer(cPayload))
                )
                switch feed.method {
                    case route.RouteMethodGet: {
                        C.PCFeedResponseForGet(cPath, cResponse)
                    }
                    case route.RouteMethodPost: {
                        C.PCFeedResponseForPost(cPath, cResponse)
                    }
                    case route.RouteMethodPut: {
                        C.PCFeedResponseForPut(cPath, cResponse)
                    }
                    case route.RouteMethodDeleteh: {
                        C.PCFeedResponseForDelete(cPath, cResponse)
                    }
                }
            }
        }
    }
}

var theFeeder = &feeder{
    feedPipe: make(chan interface{}),
}

//export StartResponseFeed
func StartResponseFeed() {
    go theFeeder.feedLoop()
}

//export StopResponseFeed
func StopResponseFeed() {
    theFeeder.feedPipe <- stopFeed{}
}

func (f *feeder) FeedResponseForGet(path, payload string) error {
    if len(path) == 0 {
        return errors.Errorf("[ERR] invalid feed path")
    }
    f.feedPipe <- feedResponse {
        method:      route.RouteMethodGet,
        path:        path,
        response:    payload,

    }
    return nil
}

func (f *feeder) FeedResponseForPost(path, payload string) error {
    if len(path) == 0 {
        return errors.Errorf("[ERR] invalid feed path")
    }
    f.feedPipe <- feedResponse {
        method:      route.RouteMethodPost,
        path:        path,
        response:    payload,
    }
    return nil
}

func (f *feeder) FeedResponseForPut(path, payload string) error {
    if len(path) == 0 {
        return errors.Errorf("[ERR] invalid feed path")
    }
    f.feedPipe <- feedResponse {
        method:      route.RouteMethodPut,
        path:        path,
        response:    payload,
    }
    return nil
}

func (f *feeder) FeedResponseForDelete(path, payload string) error {
    if len(path) == 0 {
        return errors.Errorf("[ERR] invalid feed path")
    }
    f.feedPipe <- feedResponse {
        method:      route.RouteMethodDeleteh,
        path:        path,
        response:    payload,
    }
    return nil
}
