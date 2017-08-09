# vbox

[![GoDoc](https://godoc.org/github.com/pwnall/vbox?status.svg)](https://godoc.org/github.com/pwnall/vbox)

This is a VirtualBox API client for Go, heavily inspired by
[vboxgo](https://github.com/th4t/vboxgo/).

The package's API closely mirrors the
[VirtualBox COM/XPCOM API](https://www.virtualbox.org/sdkref/), so it breaks a
lot of Go guidelines. Therefore, most users will prefer using a higher-level
library, or building their own abstractions over this library.

This is the author's first piece of Go code, so feedback is welcome.


## Usage

The package should build out of a clean checkout. The `samples` directory
contains reasonable starting code for a new library user.

The package generally follows the VirtualBox XPCOM API, meaning that it is
rather cumbersome.


## Prerequisites

Building this package requires a [cgo](https://golang.org/cmd/cgo)-enabled Go
installation. Most notably, the cgo requirement seems to preclude
cross-compilation.

The package dynamically loads `VBoxXPCOMC`, a library that implements the
VirtualBox XPCOM/COM API. The library is included with standard VirtualBox
installations.

On systems where VirtualBox is installed at a non-standard location, the
`VBOX_APP_HOME` environment variable must be set to point to the installation
location. The following example accomplishes that on 64-bit Fedora.

```bash
export VBOX_APP_HOME=/usr/lib64/virtualbox
```


## Testing

Go's standard process for running tests should work, provided that VirtualBox
is installed in a standard path, or that `VBOX_APP_HOME` is set up.

```bash
go test
```

The tests will be really slow the first time around, because they have to
download the [Lubuntu](http://lubuntu.net/) 15.04 x86 ISO. The massive delay
can be avoided by downloading the ISO manually, possibly from a local copy.

```bash
wget http://cdimage.ubuntu.com/lubuntu/releases/15.04/release/lubuntu-15.04-desktop-i386.iso
mkdir -p test_tmp
mv lubuntu-15.04-desktop-i386.iso test_tmp/lubuntu-15.04.iso
```

The dependency on a 700MB image was not taken lightly. We previously tried
using [TinyCore](http://en.wikipedia.org/wiki/Tiny_Core_Linux) and
[Damn Small Linux](http://www.damnsmalllinux.org/). Unfortunately, both
distributions have a broken mouse setup, which causes failures in the mouse
automation tests (`mouse_test.go`). Suggestions for eliminating the dependency
on Lubuntu are welcome.


## Debugging

When debugging failing tests, it is useful to start the `VBoxSVC` process in a
console, and inspect its console output. VirtualBox and the API client start
the process automatically, but it dies after 5 seconds of inactivity. So,
keeping the VirtualBox UI closed for 5 seconds should get rid of the existing
process.

The following environment variables enable logging in Release builds of
VboxSVC, which are included in the downloadable packages and most
distributions. The variables were listed off of the
[https://www.virtualbox.org/wiki/VBoxMainLogging](VirtualBox wiki).

```bash
export VBOXSVC_RELEASE_LOG=main.e.l.f+gui.e.l.f
export VBOXSVC_RELEASE_LOG_FLAGS="time tid thread"
export VBOXSVC_RELEASE_LOG_DEST=stdout
```


## Vendored VirtualBox SDK

The package contains a subset of the VirtualBox SDK, under
`VirtualBoxSDK`. The original version was obtained by unzipping the
SDK package on the
[VirtualBox downloads page](https://www.virtualbox.org/wiki/Downloads).

The vendored version removes all the files that are not related to the C
bindings, which was necessary to keep the repository small.


## Copyright and Licensing

The licensing situation of this package is complicated due to issues outside of
the author's control. Briefly, you can most likely consider the package to be
MIT-licensed.

All files outside of `third_party/` are (C) Victor Costan 2015, and made
available under the MIT license, which is contained in the `LICENSE` file.

The vendored VirtualBox has most files (under `bindings/c/glue`) licensed under
the MIT license. However, one header file (under `bindings/c/include`) is
licensed under the LGPL v2.

According to the VirtualBox API developers, the header files do not generate
code, so including them should not activate LGPL's viral infection clause.
The details are in the forum posts surrounding
[this post](https://forums.virtualbox.org/viewtopic.php?f=34&t=65063#p323121).
