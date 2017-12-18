package main

import "github.com/quimcalpe/iracing-sdk"

func main() {
	var sdk irsdk.IRSDK
	sdk = irsdk.Init(nil)
	defer sdk.Close()
	sdk.ExportIbtTo("data.ibt")
	sdk.ExportSessionTo("data.yml")
}
