#Golang iRacing SDK

Golang implementation of iRacing SDK

##Install

1. Execute `go get gituhub.com/quimcalpe/iracing-sdk`
2. Add "gituhub.com/quimcalpe/iracing-sdk" to your imports

##Usage

Simplest example:
```go
package main

import (
    "fmt"
    "gituhub.com/quimcalpe/iracing-sdk"
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

    "gituhub.com/quimcalpe/iracing-sdk"
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

##Examples

* [Export](examples/export) Telemetry Data and Session yaml to files

* Broadcast [Commands](examples/commands) to iRacing

* Simple [Dashboard](examples/dashboard) for external monitors or phones