package irsdk

import "io"

const dataValidEventName string = "Local\\IRSDKDataValidEvent"
const fileMapName string = "Local\\IRSDKMemMapFileName"
const fileMapSize int32 = 1164 * 1024
const broadcastMsgName string = "IRSDK_BROADCASTMSG"
const connTimeout = 30

const (
	stConnected int = 1
)

type Msg struct {
	Cmd int
	P1  int
	P2  interface{}
	P3  int
}

type reader interface {
	io.Reader
	io.ReaderAt
	io.Closer
}

const (
	BroadcastCamSwitchPos            int = 0  // car position, group, camera
	BroadcastCamSwitchNum            int = 1  // driver #, group, camera
	BroadcastCamSetState             int = 2  // irsdk_CameraState, unused, unused
	BroadcastReplaySetPlaySpeed      int = 3  // speed, slowMotion, unused
	BroadcastReplaySetPlayPosition   int = 4  // irsdk_RpyPosMode, Frame Number (high, low)
	BroadcastReplaySearch            int = 5  // irsdk_RpySrchMode, unused, unused
	BroadcastReplaySetState          int = 6  // irsdk_RpyStateMode, unused, unused
	BroadcastReloadTextures          int = 7  // irsdk_ReloadTexturesMode, carIdx, unused
	BroadcastChatComand              int = 8  // irsdk_ChatCommandMode, subCommand, unused
	BroadcastPitCommand              int = 9  // irsdk_PitCommandMode, parameter
	BroadcastTelemCommand            int = 10 // irsdk_TelemCommandMode, unused, unused
	BroadcastFFBCommand              int = 11 // irsdk_FFBCommandMode, value (float, high, low)
	BroadcastReplaySearchSessionTime int = 12 // sessionNum, sessionTimeMS (high, low)
	BroadcastLast                    int = 13 // unused placeholder
)

const (
	ChatCommandMacro     int = 0 // pass in a number from 1-15 representing the chat macro to launch
	ChatCommandBeginChat int = 1 // Open up a new chat window
	ChatCommandReply     int = 2 // Reply to last private chat
	ChatCommandCancel    int = 3 // Close chat window
)

// this only works when the driver is in the car
const (
	PitCommandClear      int = 0  // Clear all pit checkboxes
	PitCommandWS         int = 1  // Clean the winshield, using one tear off
	PitCommandFuel       int = 2  // Add fuel, optionally specify the amount to add in liters or pass '0' to use existing amount
	PitCommandLF         int = 3  // Change the left front tire, optionally specifying the pressure in KPa or pass '0' to use existing pressure
	PitCommandRF         int = 4  // right front
	PitCommandLR         int = 5  // left rear
	PitCommandRR         int = 6  // right rear
	PitCommandClearTires int = 7  // Clear tire pit checkboxes
	PitCommandFR         int = 8  // Request a fast repair
	PitCommandClearWS    int = 9  // Uncheck Clean the winshield checkbox
	PitCommandClearFR    int = 10 // Uncheck request a fast repair
	PitCommandClearFuel  int = 11 // Uncheck add fuel
)

// You can call this any time, but telemtry only records when driver is in there car
const (
	TelemCommandStop    int = 0 // Turn telemetry recording off
	TelemCommandStart   int = 1 // Turn telemetry recording on
	TelemCommandRestart int = 2 // Write current file to disk and start a new one
)

const (
	RpyStateEraseTape int = 0 // clear any data in the replay tape
	RpyStateLast      int = 1 // unused place holder
)

const (
	ReloadTexturesAll    int = 0 // reload all textuers
	ReloadTexturesCarIdx int = 1 // reload only textures for the specific carIdx
)

// Search replay tape for events
const (
	RpySrchToStart      int = 0
	RpySrchToEnd        int = 1
	RpySrchPrevSession  int = 2
	RpySrchNextSession  int = 3
	RpySrchPrevLap      int = 4
	RpySrchNextLap      int = 5
	RpySrchPrevFrame    int = 6
	RpySrchNextFrame    int = 7
	RpySrchPrevIncident int = 8
	RpySrchNextIncident int = 9
	RpySrchLast         int = 10 // unused placeholder
)

const (
	RpyPosBegin   int = 0
	RpyPosCurrent int = 1
	RpyPosEnd     int = 2
	RpyPosLast    int = 3 // unused placeholder
)

// You can call this any time
const (
	FFBCommandMaxForce int = 0 // Set the maximum force when mapping steering torque force to direct input units (float in Nm)
	FFBCommandLast     int = 1 // unused placeholder
)

// irsdk_BroadcastCamSwitchPos or irsdk_BroadcastCamSwitchNum camera focus defines
// pass these in for the first parameter to select the 'focus at' types in the camera system.
const (
	csFocusAtIncident int = -3
	csFocusAtLeader   int = -2
	csFocusAtExiting  int = -1
	csFocusAtDriver   int = 0 // ctFocusAtDriver + car number...

)
