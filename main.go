package main

import (
	"fmt"
)

func main() {
	var testDevice bool
	var startAutoclicker []string
	var help bool
	ParseArgs([]ArgumentData{
		{Keys: []string{"testDevice"}, AfterCount: 0, Target: &testDevice},
		{Keys: []string{"startAutoclicker"}, AfterCount: 2, Target: &startAutoclicker},
		{Keys: []string{"help", "h"}, AfterCount: 0, Target: &help}, // Add this
	})

	if testDevice {
		getDeviceToUser()
	} else if startAutoclicker != nil {
		runAutoclicker(startAutoclicker[0], startAutoclicker[1])
	} else {
		fmt.Println("No flags provided. Use --help for options.")
	}
}
