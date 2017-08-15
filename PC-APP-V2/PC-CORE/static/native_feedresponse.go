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

    "github.com/pkg/errors"
)

type feedMethodType int
const (
    feedReponseGet       feedMethodType = iota
    feedResponsePost
    feedResponsePut
    feedResponseDelete
)

type feedResponse struct {
    method      feedMethodType
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
                switch feed.method {
                    case feedReponseGet: {
                        var (
                            cPath = C.CString(feed.path)
                            cResponse = C.CString(feed.response)
                        )
                        C.PCFeedResponseForGet(cPath, cResponse)
                        // these strings will be freed on cocoa side to reduce performance degrade
                        //C.free(unsafe.Pointer(cPath))
                        //C.free(unsafe.Pointer(cPayload))
                    }
                    case feedResponsePost: {
                        var (
                            cPath = C.CString(feed.path)
                            cResponse = C.CString(feed.response)
                        )
                        C.PCFeedResponseForPost(cPath, cResponse)
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

//export StartResponseFeed
func StartResponseFeed() {
    go theFeeder.feedLoop()
}

//export StopResponseFeed
func StopResponseFeed() {
    theFeeder.feedPipe <- stopFeed{}
}

func FeedResponseForGet(path, payload string) error {
    if len(path) == 0 {
        return errors.Errorf("[ERR] invalid feed path")
    }
    theFeeder.feedPipe <- feedResponse {
        method:      feedReponseGet,
        path:        path,
        response:    payload,

    }
    return nil
}

func FeedResponseForPost(path, payload string) error {
    if len(path) == 0 {
        return errors.Errorf("[ERR] invalid feed path")
    }
    theFeeder.feedPipe <- feedResponse {
        method:      feedResponsePost,
        path:        path,
        response:    payload,
    }
    return nil
}
