package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_osxmain -framework Cocoa

extern int osxmain(int argc, const char * argv[]);
*/
import "C"

func main() {
    // Perhaps the first thing main() function needs to do is initiate OSX main
    C.osxmain(0, nil)
}