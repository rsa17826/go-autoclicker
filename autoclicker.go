package main

import (
	"encoding/binary"
	"fmt"
	"net"
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
					vMouse.SendEvent(input.EV_KEY, input.BTN_LEFT, 1)
					vMouse.Sync()
					vMouse.SendEvent(input.EV_KEY, input.BTN_LEFT, 0)
					vMouse.Sync()
				}
				if rightDown {
					vMouse.SendEvent(input.EV_KEY, input.BTN_RIGHT, 1)
					vMouse.Sync()
					vMouse.SendEvent(input.EV_KEY, input.BTN_RIGHT, 0)
					vMouse.Sync()
				}
				time.Sleep(40 * time.Millisecond)
			} else {
				// If not clicking, sleep a bit so we don't pin the CPU
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	conn, err := net.Dial("unix", "/tmp/kbd_manager.sock")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Fprintln(conn, "FILTER")

	for {
		var shouldBlock bool
		var ev WireEvent
		err := binary.Read(conn, binary.LittleEndian, &ev)
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		println(ev.Type)
		if ev.Type == input.EV_KEY {
			if ev.Code == input.KEY_SCROLLLOCK {
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
			if ev.Code == input.BTN_LEFT {
				leftDown = (ev.Value == 1)
				if turboEnabled {
					shouldBlock = true
				}
			}
			if ev.Code == input.BTN_RIGHT {
				rightDown = (ev.Value == 1)
				if turboEnabled {
					shouldBlock = true
				}
			}
		}
		if shouldBlock {
			conn.Write([]byte{'1'}) // Tell server NOT to pass this to the virtual device
		} else {
			conn.Write([]byte{'0'}) // Tell server to pass it through normally
		}
	}
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
