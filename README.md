# Golang iRacing SDK

Golang implementation of iRacing SDK

## Install

You need a gcc compiler to build the SDK, Golang does not have (as far as I know) unsafe low level access to memory map files and windows broadcast events, so CGO is used to bridge this functions with C native ones.
The easiest way is to install MiniGw for 64 bits: https://sourceforge.net/projects/mingw-w64/

With a gcc compiler in place, you can follow the standard path get to external libs in Go
1. Execute `go get github.com/quimcalpe/iracing-sdk`

## Usage

Simplest example:
```go
package main

import (
    "fmt"
    "github.com/quimcalpe/iracing-sdk"
)

func main() {
	var sdk irsdk.IRSDK
	sdk = irsdk.Init(nil)
	defer sdk.Close()
    speed, _ := sdk.GetVar("Speed")
    fmt.Printf("Speed: %s", speed)
}
```

Get data in a loop live
```go
package main

import (
    "fmt"
    "log"

    "github.com/quimcalpe/iracing-sdk"
)

func main() {
	var sdk irsdk.IRSDK
	sdk = irsdk.Init(nil)
	defer sdk.Close()

    for {
        sdk.WaitForData(100 * time.Millisecond)
        speed, err := sdk.GetVar("Speed")
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Speed: %s", speed)
    }
}
```

Work with an offline ibt file
```go
reader, err := os.Open("data.ibt")
if err != nil {
    log.Fatal(err)
}
sdk = irsdk.Init(reader)
...
```

## Examples

* [Export](examples/export) Telemetry Data and Session yaml to files

* Broadcast [Commands](examples/commands) to iRacing

* Simple [Dashboard](examples/dashboard) for external monitors or phones