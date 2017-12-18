package winevents

/*
#include <windows.h>

HANDLE _OpenEvent(CHAR *event_name) {
	return OpenEvent(SYNCHRONIZE, FALSE, event_name);
}

DWORD _WaitForSingleObject(HANDLE h, DWORD timeout) {
	return WaitForSingleObject(h,timeout);
}

void _SendNotifyMessage(UINT msgID, int msg, int var1, int var2, int var3) {
 	SendNotifyMessage(HWND_BROADCAST, msgID, MAKELONG(msg, var1), MAKELONG(var2, var3));
}
*/
import "C"
import (
	"log"
	"time"
)

var eventHandle C.HANDLE

func OpenEvent(eventName string) {
	evt := (*C.CHAR)(C.CString(eventName))
	eventHandle = C._OpenEvent(evt)
}

func WaitForSingleObject(timeout time.Duration) bool {
	t0 := time.Now().UnixNano()
	timeoutInt := int(timeout / time.Millisecond)
	r := C._WaitForSingleObject(eventHandle, C.DWORD(timeoutInt))
	if C.GetLastError() != 0 {
		remainingTimeout := timeoutInt - int((time.Now().UnixNano()-t0)/1000000)
		if remainingTimeout > 0 {
			time.Sleep(time.Duration(remainingTimeout) * time.Millisecond)
		}
		return false
	}
	return r == 0
}

func BroadcastMsg(msgName string, msg int, p1 int, p2 interface{}, p3 int) bool {
	var p2Int int
	switch v := p2.(type) {
	case int, int8, int16, int32, int64:
		p2Int = v.(int)
	case float32, float64:
		p2Int = (int)(v.(float64) * 65536.0)
	default:
		log.Fatal("Second param must be an int or a float")
	}
	msgNameChar := (*C.CHAR)(C.CString(msgName))
	msgID := C.RegisterWindowMessage(msgNameChar)
	if msgID < 0 {
		return false
	}
	C._SendNotifyMessage(msgID, C.int(msg), C.int(p1), C.int(p2Int), C.int(p3))
	return true
}
