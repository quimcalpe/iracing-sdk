package main

//#include<conio.h>
import "C"

import (
	"fmt"
	"iracing-sdk"
	"os"
)

func main() {
	var sdk irsdk.IRSDK
	sdk = irsdk.Init(nil)
	defer sdk.Close()

	fmt.Println("Available commands:")
	fmt.Println(" c -> Open chat")
	fmt.Println(" p -> Clear tire pit checkboxes")
	fmt.Println(" f -> Change FFB mode")
	fmt.Println("Press Esc to exit")

	trueFFBState := false
	for {
		c := int(C.getch())
		switch c {
		case 27: // esc
			os.Exit(0)
		case 99: // c
			sdk.BroadcastMsg(irsdk.Msg{
				Cmd: irsdk.BroadcastChatComand,
				P1:  irsdk.ChatCommandBeginChat,
			})
			fmt.Println("* Send request to start a chat")
		case 112: // p
			sdk.BroadcastMsg(irsdk.Msg{
				Cmd: irsdk.BroadcastPitCommand,
				P1:  irsdk.PitCommandClearTires,
			})
			fmt.Println("* Send request to clear tire checkboxes")
		case 102: // f
			var force float64
			if trueFFBState {
				force = -1.0
			} else {
				force = 20.9998
			}
			trueFFBState = !trueFFBState
			sdk.BroadcastMsg(irsdk.Msg{
				Cmd: irsdk.BroadcastFFBCommand,
				P1:  irsdk.FFBCommandMaxForce,
				P2:  force,
			})
			if force < 0 {
				fmt.Println("* Set wheel to user controlled FFB")
			} else {
				fmt.Printf("* Set wheel to %f Nm\n", force)
			}
		default:
			fmt.Println("Unknown command")
		}
	}

}
