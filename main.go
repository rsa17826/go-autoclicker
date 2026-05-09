package main

import (
	"fmt"

	argparse "github.com/rsa17826/go-arg-lib"
	"github.com/rsa17826/go-input-lib"
)

func main() {
	var testDevice bool
	var waitForDevice bool
	var startAutoclicker []string
	var help bool
	argparse.ParseArgs([]argparse.ArgumentData{
		{Keys: []string{"testDevice"}, AfterCount: 0, Target: &testDevice},
		{Keys: []string{"startAutoclicker"}, AfterCount: 2, Target: &startAutoclicker},
		{Keys: []string{"waitForDevice"}, AfterCount: 2, Target: &waitForDevice},
		{Keys: []string{"help", "h"}, AfterCount: 0, Target: &help},
	})

	if help {
		showHelp()
	} else if testDevice {
		input.GetDeviceToUser()
	} else if startAutoclicker != nil {
		runAutoclicker(startAutoclicker[0], startAutoclicker[1], waitForDevice)
	} else {
		showHelp()
	}
}

func showHelp() {
	fmt.Println("--testDevice then press a key or button and it will return the id or name of the device")
	fmt.Println("--startAutoclicker takes a mouse and a keyboard")
	fmt.Println("--waitForDevice should be used if the device might not be active yet and should be waited for before starting")
}
