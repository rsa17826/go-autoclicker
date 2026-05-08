package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

// InputEvent matches the 'input_event' struct in linux/input.h
type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

const EVIOCGNAME = 0x80ff4506

func getDeviceName(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return "Unknown (Permission Denied)"
	}
	defer f.Close()

	// Create a buffer to hold the name (up to 256 chars)
	name := make([]byte, 256)

	// Perform the ioctl syscall
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(EVIOCGNAME),
		uintptr(unsafe.Pointer(&name[0])),
	)

	if errno != 0 {
		return "Unknown Device"
	}

	// Trim the null characters from the buffer
	return string(bytes.Trim(name, "\x00"))
}

func getPersistentID(eventPath string) string {
	// eventPath is like "/dev/input/event4"
	absPath, _ := filepath.Abs(eventPath)

	matches, _ := filepath.Glob("/dev/input/by-id/*")
	for _, idPath := range matches {
		evalPath, _ := filepath.EvalSymlinks(idPath)
		if evalPath == absPath {
			return idPath // Found the persistent ID
		}
	}
	return "" // No persistent ID found (likely a virtual device)
}

func getDeviceFromIdOrName(input string) (*os.File, error) {
	var devicePath string

	if strings.HasPrefix(input, "id:") {
		// Use the persistent symlink in /dev/input/by-id/
		idPart := strings.TrimPrefix(input, "id:")
		devicePath = filepath.Join("/dev/input/by-id", idPart)

	} else if strings.HasPrefix(input, "name:") {
		// Scan all events to find the one with the matching Name
		targetName := strings.TrimPrefix(input, "name:")

		files, err := os.ReadDir("/dev/input/")
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			if strings.HasPrefix(f.Name(), "event") {
				path := filepath.Join("/dev/input/", f.Name())
				if getDeviceName(path) == targetName {
					devicePath = path
					break
				}
			}
		}
	}

	if devicePath == "" {
		return nil, fmt.Errorf("device not found for input: %s", input)
	}

	// Open the file and return the file pointer
	return os.Open(devicePath)
}
