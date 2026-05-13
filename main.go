package main

import (
	argparse "github.com/rsa17826/go-arg-lib"
)

func main() {
	// var testDevice bool
	// var waitForDevice bool
	var startAutoclicker bool
	var downFor int
	var upFor int
	args := []argparse.ArgumentData{
		// {Keys: []string{"testDevice"}, AfterCount: 0, Target: &testDevice, Description: "interactive device selection"},
		{Keys: []string{"startAutoclicker"}, AfterCount: 0, Target: &startAutoclicker, Description: "starts the autoclicker"},
		{Keys: []string{"downFor"}, AfterCount: 1, Target: &downFor, Description: "how long each press lasts", ExampleArgs: []string{"ms"}, Default: []any{10}},
		{Keys: []string{"upFor"}, AfterCount: 1, Target: &upFor, Description: "how long between each press", ExampleArgs: []string{"ms"}, Default: []any{10}},
		// {Keys: []string{"startAutoclicker"}, AfterCount: 2, Target: &startAutoclicker, Description: "starts the autoclicker", ExampleArgs: []string{"keyboard", "mouse"}},
		// {Keys: []string{"waitForDevice"}, AfterCount: 0, Target: &waitForDevice},
	}
	argparse.ParseArgs(args)

	if startAutoclicker {
		runAutoclicker(downFor, upFor)
	} else {
		argparse.PrintHelp(args, nil)
	}
}
