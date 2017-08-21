package main

/*
#cgo CFLAGS: -x objective-c
*/
import "C"
import (
    "github.com/stkim1/pc-core/event/route"
)

//export RouteRequestGet
func RouteRequestGet(path *C.char) {
    theApp.eventsIn <- route.RouteRequestGet(C.GoString(path))
}

//export RouteRequestPost
func RouteRequestPost(path, request *C.char) {
    theApp.eventsIn <- route.RouteRequestPost(C.GoString(path), C.GoString(request))
}

//export RouteRequestPut
func RouteRequestPut(path, request *C.char) {
    theApp.eventsIn <- route.RouteRequestPut(C.GoString(path), C.GoString(request))
}

//export RouteRequestDelete
func RouteRequestDelete(path *C.char) {
    theApp.eventsIn <- route.RouteRequestDelete(C.GoString(path))
}
