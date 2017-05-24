# Go Plug-ins & Vendored Dependencies ([\#20481](https://github.com/golang/go/issues/20481))
With the release of Go 1.8 came a feature long-sought by many developers --
support for modular plug-ins loadable at runtime. While Go plug-ins do have
some limitations today -- primarily being Linux only at this time -- they
are still incredibly useful.

Unless a project has vendored dependencies that is.

The utility of Go plug-ins is almost completely erased by fact that many
Go projects rely on vendored dependencies in order to ensure consistent
build results.

## The Problem
The problem is pretty straight-forward. When an application (`app`)
vendors a library (`lib`), the package path of the library is now
`path/to/app/vendor/path/to/lib`. However, the plug-in is likely
built against either `path/to/lib` or, if the plug-in vendors
dependencies as well, `path/to/plugin/vendor/path/to/lib`.

This of course makes total sense and behaves exactly as one would
expect with regards to Go packages. Despite the intent, these three
packages are *not* the same:

* `path/to/lib`
* `path/to/app/vendor/path/to/lib`
* `path/to/plugin/vendor/path/to/lib`

While the behavior is consistent with regards to Go packages, it
flies in the face of the utility provided by a combination of
vendored dependencies and the new Go plug-in model.

## Reproduction
This project makes it easy to reproduce the above issue.

### Requirements
To reproduce this issue Go 1.8.x and a Linux host are required:

```bash
$ go version
go version go1.8.1 linux/amd64
```

```bash
$ go env
GOARCH="amd64"
GOBIN=""
GOEXE=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOOS="linux"
GOPATH="/home/akutz/go"
GORACE=""
GOROOT="/home/akutz/.go/1.8.1"
GOTOOLDIR="/home/akutz/.go/1.8.1/pkg/tool/linux_amd64"
GCCGO="gccgo"
CC="gcc"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build699913681=/tmp/go-build -gno-record-gcc-switches"
CXX="g++"
CGO_ENABLED="1"
PKG_CONFIG="pkg-config"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
```

### Download
On a Linux host use `go get` to fetch this project:

```bash
$ go get github.com/akutz/gpd
```

### Run the program
The root of the project is a Go command-line program. Running it
will emit a message to the console:

```bash
$ go run main.go
Yes, we have no bananas,
We have no bananas today.
```

### Build the plug-in
If the program is run with a single argument it is treated as the
path to a Go plug-in. That plug-in is loaded and will emit a different
message to the console. First, build the plug-in:

```bash
$ go build -buildmode plugin -o mod.so ./mod
```

To verify that the produced file *is* a plug-in, use the `file` command:

```bash
$ file mod.so
mod.so: ELF 64-bit LSB shared object, x86-64, version 1 (SYSV), dynamically linked, BuildID[sha1]=8c78f9a393bd083bde91b2b34b8117592387f40e, not stripped
```

The file is reported as a *shared object*, verifying that it is indeed a
Go plug-in.

### Run the program with the plug-in
Run the program using the plug-in:

```bash
$ go run main.go mod.so
Yes there were thirty, thousand, pounds...
Of...bananas.
```

It works!

### Vendor the shared `dep` package
However, what happens when the program vendors the shared `dep` package?

```bash
$ mkdir -p vendor/github.com/akutz/gpd && cp -r dep vendor/github.com/akutz/gpd
$ go run main.go mod.so
error: failed to load plugin: plugin.Open: plugin was built with a different version of package github.com/akutz/gpd/lib
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x504c45]

goroutine 1 [running]:
github.com/akutz/gpd/lib.NewModule(0x535498, 0x6, 0x539cb5, 0x21)
	/home/akutz/go/src/github.com/akutz/gpd/lib/lib.go:28 +0x55
main.main()
	/home/akutz/go/src/github.com/akutz/gpd/main.go:32 +0x13a
exit status 2
```

The program fails!

This is because the `dep` package includes a type that
is used by both the shared `lib` package and the plug-in package, `mod`.

The plug-in linked against the `lib` package at `github.com/akutz/gpd/lib`
which itself linked against the `dep` package at `github.com/akutz/gpd/dep`.

However, vendoring the `dep` package for the program causes the `lib`
package as compiled into the program to link against
`github.com/akutz/gpd/vendor/github.com/akutz/gpd/dep`, resulting in
the program and the plug-in having two different versions of the `lib`
package!

### Vendor the shared `lib` package
However, what happens when the program vendors the shared `lib` package?

```bash
$ rm -fr vendor
$ mkdir -p vendor/github.com/akutz/gpd && cp -r lib vendor/github.com/akutz/gpd
$ go run main.go mod.so
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x504d65]

goroutine 1 [running]:
github.com/akutz/gpd/vendor/github.com/akutz/gpd/lib.NewModule(0x5355b8, 0x6, 0xc42000c2c0, 0x0)
	/home/akutz/go/src/github.com/akutz/gpd/vendor/github.com/akutz/gpd/lib/lib.go:28 +0x55
main.main()
	/home/akutz/go/src/github.com/akutz/gpd/main.go:32 +0x13a
exit status 2
```

The program fails! This is because the `lib` package contains a type
registry that can be used to both register types and construct new
instances of those types.

However, because the program's type registry is located in the package
`github.com/akutz/gpd/vendor/github.com/akutz/gpd/lib` and the plug-in
registered its type with `github.com/akutz/gpd/lib`, when the program
requests a new object for the type `mod_go`, a nil exception occurs
because the program and plug-in were accessing two different type
registries!

## The Hack
At the moment the only solution available is to create a build
toolchain using a list of transitive dependencies generated from
the application that is responsible for loading the plug-ins. This
list of dependencies can be used to create a custom `GOPATH` against
which any projects participating in the application must be built,
including the application itself, any shared libraries, and the
plug-ins.

## The Solution
Is there one? Two possible solutions are:

1. Allow a `src` directory at the root of a `vendor` directory so that
plug-ins can be built directly against a program's `vendor` directory.
Today that would require a bind mount.
2. Allow plug-ins to link directly against the Go program binary that
will load the programs.

Hopefully the Golang team can solve this issue as it really does prevent
Go plug-ins from being useful in a world where applications are often
required to vendor dependencies.
