package irsdk

import (
	"log"
)

type header struct {
	version  int
	status   int
	tickRate int // ticks per second (60 or 360 etc)

	// session information, updated periodicaly
	sessionInfoUpdate int // Incremented when session info changes
	sessionInfoLen    int // Length in bytes of session info string
	sessionInfoOffset int // Session info, encoded in YAML format

	// state data, output at tickRate Hz
	numVars      int // length of array pointed to by varHeaderOffset
	headerOffset int // offset to irsdk_varHeader[numVars] array, Describes the variables received in varBuf

	numBuf int
	bufLen int // length in bytes for one line
}

func readHeader(r reader) header {
	rbuf := make([]byte, 48)
	_, err := r.ReadAt(rbuf, 0)
	if err != nil {
		log.Fatal(err)
	}
	h := header{
		byte4ToInt(rbuf[0:4]),
		byte4ToInt(rbuf[4:8]),
		byte4ToInt(rbuf[8:12]),
		byte4ToInt(rbuf[12:16]),
		byte4ToInt(rbuf[16:20]),
		byte4ToInt(rbuf[20:24]),
		byte4ToInt(rbuf[24:28]),
		byte4ToInt(rbuf[28:32]),
		byte4ToInt(rbuf[32:36]),
		byte4ToInt(rbuf[36:40]),
	}
	//fmt.Printf("%v\n%+v\n", rbuf, h)
	return h
}
