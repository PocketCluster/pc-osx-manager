package main

/*
#cgo CFLAGS: -x objective-c
*/
import "C"
import (
    "github.com/stkim1/pc-core/event/route"
)

//export RouteEventGet
func RouteEventGet(path *C.char) {
    theApp.eventsIn <- route.RouteRequestEventGet(C.GoString(path))
}

//export RouteEventPost
func RouteEventPost(path, payload *C.char) {
    theApp.eventsIn <- route.RouteRequestEventPost(C.GoString(path), C.GoString(payload))
}

//export RouteEventPut
func RouteEventPut(path, payload *C.char) {
    theApp.eventsIn <- route.RouteRequestEventPut(C.GoString(path), C.GoString(payload))
}

//export RouteEventDelete
func RouteEventDelete(path *C.char) {
    theApp.eventsIn <- route.RouteRequestEventDelete(C.GoString(path))
}
