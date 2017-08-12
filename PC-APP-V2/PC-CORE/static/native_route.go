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
    theApp.eventsIn <- route.RouteEventGet(C.GoString(path))
}

//export RouteEventPost
func RouteEventPost(path, payload *C.char) {
    theApp.eventsIn <- route.RouteEventPost(C.GoString(path), C.GoString(payload))
}

//export RouteEventPut
func RouteEventPut(path, payload *C.char) {
    theApp.eventsIn <- route.RouteEventPut(C.GoString(path), C.GoString(payload))
}

//export RouteEventDelete
func RouteEventDelete(path *C.char) {
    theApp.eventsIn <- route.RouteEventDelete(C.GoString(path))
}
