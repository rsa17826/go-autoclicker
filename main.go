package main

import (
	"fmt"

	argparse "github.com/rsa17826/go-arg-lib"
	"github.com/rsa17826/go-input-lib"
)

func main() {
	var testDevice bool
	var startAutoclicker []string
	var help bool
	argparse.ParseArgs([]argparse.ArgumentData{
		{Keys: []string{"testDevice"}, AfterCount: 0, Target: &testDevice},
		{Keys: []string{"startAutoclicker"}, AfterCount: 2, Target: &startAutoclicker},
		{Keys: []string{"help", "h"}, AfterCount: 0, Target: &help},
	})

	if testDevice {
		input.GetDeviceToUser()
	} else if startAutoclicker != nil {
		runAutoclicker(startAutoclicker[0], startAutoclicker[1])
	} else {
		fmt.Println("No flags provided. Use --help for options.")
	}
}
