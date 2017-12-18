package irsdk

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/hidez8891/shm"
	"github.com/quimcalpe/iracing-sdk/lib/winevents"
)

// IRSDK is the main SDK object clients must use
type IRSDK struct {
	r             reader
	h             *header
	s             []string
	tVars         *TelemetryVars
	lastValidData int64
}

func (sdk *IRSDK) WaitForData(timeout time.Duration) bool {
	if !sdk.IsConnected() {
		initIRSDK(sdk)
	}
	if winevents.WaitForSingleObject(timeout) {
		return readVariableValues(sdk)
	}
	return false
}

func (sdk *IRSDK) GetVar(name string) (variable, error) {
	if !sessionStatusOK(sdk.h.status) {
		return variable{}, fmt.Errorf("Session is not active")
	}
	sdk.tVars.mux.Lock()
	if v, ok := sdk.tVars.vars[name]; ok {
		sdk.tVars.mux.Unlock()
		return v, nil
	}
	sdk.tVars.mux.Unlock()
	return variable{}, fmt.Errorf("Telemetry variable %q not found", name)
}

func (sdk *IRSDK) GetLastVersion() int {
	if !sessionStatusOK(sdk.h.status) {
		return -1
	}
	sdk.tVars.mux.Lock()
	last := sdk.tVars.lastVersion
	sdk.tVars.mux.Unlock()
	return last
}

func (sdk *IRSDK) GetSessionData(path string) (string, error) {
	if !sessionStatusOK(sdk.h.status) {
		return "", fmt.Errorf("Session not connected")
	}
	return getSessionDataPath(sdk.s, path)
}

func (sdk *IRSDK) IsConnected() bool {
	if sdk.h != nil {
		if sessionStatusOK(sdk.h.status) && (sdk.lastValidData+connTimeout > time.Now().Unix()) {
			return true
		}
	}

	return false
}

// ExportTo exports current memory data to a file
func (sdk *IRSDK) ExportIbtTo(fileName string) {
	rbuf := make([]byte, fileMapSize)
	_, err := sdk.r.ReadAt(rbuf, 0)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(fileName, rbuf, 0644)
}

// ExportTo exports current session yaml data to a file
func (sdk *IRSDK) ExportSessionTo(fileName string) {
	y := strings.Join(sdk.s, "\n")
	ioutil.WriteFile(fileName, []byte(y), 0644)
}

func (sdk *IRSDK) BroadcastMsg(msg Msg) {
	if msg.P2 == nil {
		msg.P2 = 0
	}
	winevents.BroadcastMsg(broadcastMsgName, msg.Cmd, msg.P1, msg.P2, msg.P3)
}

// Close clean up sdk resources
func (sdk *IRSDK) Close() {
	sdk.r.Close()
}

// Init creates a SDK instance to operate with
func Init(r reader) IRSDK {
	if r == nil {
		var err error
		r, err = shm.Open(fileMapName, fileMapSize)
		if err != nil {
			log.Fatal(err)
		}
	}

	sdk := IRSDK{r: r, lastValidData: 0}
	winevents.OpenEvent(dataValidEventName)
	initIRSDK(&sdk)
	return sdk
}

func initIRSDK(sdk *IRSDK) {
	h := readHeader(sdk.r)
	sdk.h = &h
	sdk.s = nil
	if sdk.tVars != nil {
		sdk.tVars.vars = nil
	}
	if sessionStatusOK(h.status) {
		sdk.s = readSessionData(sdk.r, &h)
		sdk.tVars = readVariableHeaders(sdk.r, &h)
		readVariableValues(sdk)
	}
}

func sessionStatusOK(status int) bool {
	return (status & stConnected) > 0
}
