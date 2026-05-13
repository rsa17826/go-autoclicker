package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rsa17826/go-input-lib"
)

type ClickState struct {
	TurboActive bool
}

func runAutoclicker() {
	vMouse, err := input.CreateVirtualMouse()
	if err != nil {
		panic(err)
	}
	defer vMouse.Close()

	var turboEnabled bool
	var scrollEnabled bool
	var leftDown, rightDown bool

	go func() {
		for {
			if turboEnabled && (leftDown || rightDown) {
				if leftDown {
					sendClick(vMouse, input.BTN_LEFT)
				}
				if rightDown {
					sendClick(vMouse, input.BTN_RIGHT)
				}
				time.Sleep(40 * time.Millisecond)
			} else {
				// If not clicking, sleep a bit so we don't pin the CPU
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	for {
		conn, err := net.Dial("unix", "/tmp/kbd_manager.sock")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		fmt.Fprintln(conn, "FILTER")

		for {
			var ev WireEvent
			err := binary.Read(conn, binary.LittleEndian, &ev)
			if err != nil {
				fmt.Println("read error:", err)
				return
			}
			if ev.Type == input.EV_KEY {
				if ev.Code == input.KEY_SCROLL {
					if ev.Value == 1 {
						turboEnabled = !turboEnabled
						scrollEnabled = !scrollEnabled
					}
				}
				if ev.Code == input.KEY_Z {
					if !scrollEnabled {
						turboEnabled = ev.Value > 0
					}
				}
			}
			if ev.Type == input.EV_KEY {
				if ev.Code == input.BTN_LEFT {
					leftDown = (ev.Value == 1)
					if turboEnabled {
						// block
					}
				}
				if ev.Code == input.BTN_RIGHT {
					rightDown = (ev.Value == 1)
					if turboEnabled {
						// block
					}
				}
			}
		}
	}
}

// Helper to send a virtual click (Down + Up + Sync)
func sendClick(v *os.File, code uint16) {
	binary.Write(v, binary.LittleEndian, input.InputEvent{Type: input.EV_KEY, Code: code, Value: 1})
	binary.Write(v, binary.LittleEndian, input.InputEvent{Type: input.EV_SYN, Code: 0, Value: 0})
	binary.Write(v, binary.LittleEndian, input.InputEvent{Type: input.EV_KEY, Code: code, Value: 0})
	binary.Write(v, binary.LittleEndian, input.InputEvent{Type: input.EV_SYN, Code: 0, Value: 0})
}

type WireEvent struct {
	Sec   int64
	Usec  int64
	Type  uint16
	Code  uint16
	Value int32
}

// // Helper to toggle the physical LED
// func toggleLED(f *os.File, state int) {
// 	// EV_LED = 0x11, LED_SCROLL = 0x02
// 	ev := input.InputEvent{
// 		Type:  0x11,
// 		Code:  0x02,
// 		Value: int32(state),
// 	}
// 	binary.Write(f, binary.LittleEndian, ev)

// 	// Always send a SYN event after an update
// 	syn := input.InputEvent{Type: 0x00, Code: 0, Value: 0}
// 	binary.Write(f, binary.LittleEndian, syn)
// }
