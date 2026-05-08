package main

import (
	"encoding/binary"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

type ClickState struct {
	TurboActive bool
}

func runAutoclicker(mouseid string, kbdid string) {
	mousePath, err := getDeviceFromIdOrName(mouseid)
	if err != nil {
		println(err)
	}
	kbdPath, err := getDeviceFromIdOrName(kbdid)
	if err != nil {
		println(err)
	}
	mouse, _ := os.Open(mousePath)
	keyboard, _ := os.Open(kbdPath)
	// keyboard, _ := os.OpenFile(kbdPath, os.O_RDWR, 0666)
	vMouse, _ := createVirtualMouse()
	defer mouse.Close()
	defer keyboard.Close()
	defer vMouse.Close()

	// 1. Grab the physical mouse so the OS doesn't see double-input
	unix.IoctlSetInt(int(mouse.Fd()), EVIOCGRAB, 1)
	defer unix.IoctlSetInt(int(mouse.Fd()), EVIOCGRAB, 0)

	var turboEnabled bool
	var scrollEnabled bool
	var leftDown, rightDown bool

	// 1. THE WORKER: One goroutine that runs forever
	go func() {
		for {
			if turboEnabled && (leftDown || rightDown) {
				if leftDown {
					sendClick(vMouse, BTN_LEFT)
				}
				if rightDown {
					sendClick(vMouse, BTN_RIGHT)
				}
				time.Sleep(40 * time.Millisecond)
			} else {
				// If not clicking, sleep a bit so we don't pin the CPU
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// 2. THE KEYBOARD THREAD
	go func() {
		for {
			ev := readEvent(keyboard)
			if ev.Type == EV_KEY && ev.Code == KEY_SCROLL {
				if ev.Value == 1 { // Toggle on press
					turboEnabled = !turboEnabled
					scrollEnabled = !scrollEnabled
				}
			}
			if ev.Type == EV_KEY && ev.Code == KEY_Z {
				if !scrollEnabled {
					turboEnabled = ev.Value > 0
				}
			}
		}
	}()

	// 3. THE MAIN MOUSE LOOP (The "Pass-through")
	for {
		ev := readEvent(mouse)

		if ev.Type == EV_KEY {
			if ev.Code == BTN_LEFT {
				leftDown = (ev.Value == 1)
				if turboEnabled {
					continue
				} // Block physical click
			}
			if ev.Code == BTN_RIGHT {
				rightDown = (ev.Value == 1)
				if turboEnabled {
					continue
				} // Block physical click
			}
		}

		// This line now handles EVERYTHING else:
		// Mouse movement (REL_X/Y), Scrolling (REL_WHEEL), and Sync (EV_SYN)
		binary.Write(vMouse, binary.LittleEndian, ev)
	}
}

// Helper to send a virtual click (Down + Up + Sync)
func sendClick(v *os.File, code uint16) {
	binary.Write(v, binary.LittleEndian, input_event{Type: EV_KEY, Code: code, Value: 1})
	binary.Write(v, binary.LittleEndian, input_event{Type: EV_SYN, Code: 0, Value: 0})
	binary.Write(v, binary.LittleEndian, input_event{Type: EV_KEY, Code: code, Value: 0})
	binary.Write(v, binary.LittleEndian, input_event{Type: EV_SYN, Code: 0, Value: 0})
}

func readEvent(f *os.File) input_event {
	var ev input_event
	binary.Read(f, binary.LittleEndian, &ev)
	return ev
}
func monitorKeyboard(dev *os.File, state *ClickState) {
	for {
		ev := readEvent(dev) // Assume a helper that reads input_event struct

		// KEY_Z = 44, KEY_SCROLLLOCK = 70
		if ev.Type == EV_KEY {
			if ev.Code == 44 || ev.Code == 70 {
				state.TurboActive = (ev.Value > 0) // True if pressed or held
			}
		}
	}
}
func monitorMouse(realMouse *os.File, vMouse *os.File, state *ClickState) {
	// 1. Use unix.IoctlSetInt and unix.EVIOCGRAB
	unix.IoctlSetInt(int(realMouse.Fd()), EVIOCGRAB, 1)
	defer unix.IoctlSetInt(int(realMouse.Fd()), EVIOCGRAB, 0)

	leftPressed := false
	rightPressed := false

	for {
		ev := readEvent(realMouse)

		if ev.Type == EV_KEY {
			if ev.Code == BTN_LEFT {
				leftPressed = (ev.Value == 1)
			}
			if ev.Code == BTN_RIGHT {
				rightPressed = (ev.Value == 1)
			}

			if state.TurboActive && (leftPressed || rightPressed) {
				go func() {
					for (leftPressed || rightPressed) && state.TurboActive {
						if leftPressed {
							sendClick(vMouse, BTN_LEFT) // Use your helper function
						}
						if rightPressed {
							sendClick(vMouse, BTN_RIGHT) // Use your helper function
						}
						time.Sleep(50 * time.Millisecond)
					}
				}()
				continue // Don't pass through the original click
			}
		}

		binary.Write(vMouse, binary.LittleEndian, ev)
	}
}

// // Helper to toggle the physical LED
// func toggleLED(f *os.File, state int) {
// 	// EV_LED = 0x11, LED_SCROLL = 0x02
// 	ev := input_event{
// 		Type:  0x11,
// 		Code:  0x02,
// 		Value: int32(state),
// 	}
// 	binary.Write(f, binary.LittleEndian, ev)

// 	// Always send a SYN event after an update
// 	syn := input_event{Type: 0x00, Code: 0, Value: 0}
// 	binary.Write(f, binary.LittleEndian, syn)
// }
