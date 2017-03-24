// +build darwin
package hostapi

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_PCEventHandle

#include "PCEventHandle.h"

*/
import "C"
import (
    "encoding/json"
)

const (
    FeedType    string = "feed_type"
    FeedResult  string = "feed_ret"
    FeedMessage string = "feed_msg"
)

func SendFeedBack(message string) {
    data, err := json.Marshal(map[string]string{
        FeedType:   "api-feed",
        FeedResult: "api-success",
        FeedMessage: message,
    })
    if err == nil {
        C.PCEventHandle(C.CString(string(data)))
    }
}