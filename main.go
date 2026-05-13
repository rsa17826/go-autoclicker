package main

import (
	argparse "github.com/rsa17826/go-arg-lib"
)

func main() {
	// var testDevice bool
	// var waitForDevice bool
	var startAutoclicker bool
	args := []argparse.ArgumentData{
		// {Keys: []string{"testDevice"}, AfterCount: 0, Target: &testDevice, Description: "interactive device selection"},
		{Keys: []string{"startAutoclicker"}, AfterCount: 0, Target: &startAutoclicker, Description: "starts the autoclicker"},
		// {Keys: []string{"startAutoclicker"}, AfterCount: 2, Target: &startAutoclicker, Description: "starts the autoclicker", ExampleArgs: []string{"keyboard", "mouse"}},
		// {Keys: []string{"waitForDevice"}, AfterCount: 0, Target: &waitForDevice},
	}
	argparse.ParseArgs(args)

	if startAutoclicker {
		runAutoclicker()
	} else {
		argparse.PrintHelp(args, nil)
	}
}
