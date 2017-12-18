package irsdk

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type varBuffer struct {
	tickCount int // used to detect changes in data
	bufOffset int // offset from header
}

type variable struct {
	varType     int // irsdk_VarType
	offset      int // offset fron start of buffer row
	count       int // number of entrys (array) so length in bytes would be irsdk_VarTypeBytes[type] * count
	countAsTime bool
	Name        string
	Desc        string
	Unit        string
	Value       interface{}
	rawBytes    []byte
}

func (v variable) String() string {
	var ret string
	switch v.varType {
	case 0:
		ret = fmt.Sprintf("%c", v.Value)
	case 1:
		ret = fmt.Sprintf("%v", v.Value)
	case 2:
		ret = fmt.Sprintf("%d", v.Value)
	case 3:
		ret = fmt.Sprintf("%s", v.Value)
	case 4:
		ret = fmt.Sprintf("%f", v.Value)
	case 5:
		ret = fmt.Sprintf("%f", v.Value)
	default:
		ret = fmt.Sprintf("Unknown (%d)", v.varType)
	}
	return ret
}

// TelemetryVars holds all variables we can read from telemetry live
type TelemetryVars struct {
	lastVersion int
	vars        map[string]variable
	mux         sync.Mutex
}

func findLatestBuffer(r reader, h *header) varBuffer {
	var vb varBuffer
	foundTickCount := 0
	for i := 0; i < h.numBuf; i++ {
		rbuf := make([]byte, 16)
		_, err := r.ReadAt(rbuf, int64(48+i*16))
		if err != nil {
			log.Fatal(err)
		}
		currentVb := varBuffer{
			byte4ToInt(rbuf[0:4]),
			byte4ToInt(rbuf[4:8]),
		}
		//fmt.Printf("BUFF?: %+v\n", currentVb)
		if foundTickCount < currentVb.tickCount {
			foundTickCount = currentVb.tickCount
			vb = currentVb
		}
	}
	//fmt.Printf("BUFF: %+v\n", vb)
	return vb
}

func readVariableHeaders(r reader, h *header) *TelemetryVars {
	vars := TelemetryVars{vars: make(map[string]variable, h.numVars)}
	for i := 0; i < h.numVars; i++ {
		rbuf := make([]byte, 144)
		_, err := r.ReadAt(rbuf, int64(h.headerOffset+i*144))
		if err != nil {
			log.Fatal(err)
		}
		v := variable{
			byte4ToInt(rbuf[0:4]),
			byte4ToInt(rbuf[4:8]),
			byte4ToInt(rbuf[8:12]),
			int(rbuf[12]) > 0,
			bytesToString(rbuf[16:48]),
			bytesToString(rbuf[48:112]),
			bytesToString(rbuf[112:144]),
			nil,
			nil,
		}
		vars.vars[v.Name] = v
	}
	return &vars
}

func readVariableValues(sdk *IRSDK) bool {
	newData := false
	if sessionStatusOK(sdk.h.status) {
		// find latest buffer for variables
		vb := findLatestBuffer(sdk.r, sdk.h)
		sdk.tVars.mux.Lock()
		if sdk.tVars.lastVersion < vb.tickCount {
			newData = true
			sdk.tVars.lastVersion = vb.tickCount
			sdk.lastValidData = time.Now().Unix()
			for varName, v := range sdk.tVars.vars {
				var rbuf []byte
				switch v.varType {
				case 0:
					rbuf = make([]byte, 1)
					_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
					if err != nil {
						log.Fatal(err)
					}
					v.Value = string(rbuf[0])
				case 1:
					rbuf = make([]byte, 1)
					_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
					if err != nil {
						log.Fatal(err)
					}
					v.Value = int(rbuf[0]) > 0
				case 2:
					rbuf = make([]byte, 4)
					_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
					if err != nil {
						log.Fatal(err)
					}
					v.Value = byte4ToInt(rbuf)
				case 3:
					rbuf = make([]byte, 4)
					_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
					if err != nil {
						log.Fatal(err)
					}
					v.Value = byte4toBitField(rbuf)
				case 4:
					rbuf = make([]byte, 4)
					_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
					if err != nil {
						log.Fatal(err)
					}
					v.Value = byte4ToFloat(rbuf)
				case 5:
					rbuf = make([]byte, 8)
					_, err := sdk.r.ReadAt(rbuf, int64(vb.bufOffset+v.offset))
					if err != nil {
						log.Fatal(err)
					}
					v.Value = byte8ToFloat(rbuf)
				}
				v.rawBytes = rbuf
				sdk.tVars.vars[varName] = v
			}
		}
		sdk.tVars.mux.Unlock()
	}

	return newData
}
