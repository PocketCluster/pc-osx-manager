// +build darwin
package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_PCEventHandle

#include "PCEventHandle.h"

*/
import "C"
import (
    "encoding/json"
    "runtime"
)

const (
    FeedType    string = "feed_type"
    FeedResult  string = "feed_ret"
    FeedMessage string = "feed_msg"
)

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
        case string:
            // TODO : sort out f with detailed category
            data, err := json.Marshal(map[string]string{
                FeedType:   "api-feed",
                FeedResult: "api-success",
                FeedMessage: feed,
            })
            if err == nil {
                C.PCEventHandle(C.CString(string(data)))
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

func FeedSend(message string) {
    theFeeder.feedPipe <- message
}
