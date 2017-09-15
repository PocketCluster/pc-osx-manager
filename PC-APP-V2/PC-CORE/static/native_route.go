package main

/*
#cgo CFLAGS: -x objective-c
*/
import "C"
import (
    "github.com/stkim1/pc-core/route"
)

//export RouteRequestGet
func RouteRequestGet(path *C.char) {
    theApp.Send(route.RouteRequestGet(C.GoString(path)))
}

//export RouteRequestPost
func RouteRequestPost(path, request *C.char) {
    theApp.Send(route.RouteRequestPost(C.GoString(path), C.GoString(request)))
}

//export RouteRequestPut
func RouteRequestPut(path, request *C.char) {
    theApp.Send(route.RouteRequestPut(C.GoString(path), C.GoString(request)))
}

//export RouteRequestDelete
func RouteRequestDelete(path *C.char) {
    theApp.Send(route.RouteRequestDelete(C.GoString(path)))
}
