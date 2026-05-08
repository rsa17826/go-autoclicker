package main

import (
	"fmt"
)

func main() {
	results := CheckArgs([]ArgumentData{
		{Keys: []string{"testDevice"}, AfterCount: 0, Default: false},
		{Keys: []string{"startAutoclicker"}, AfterCount: 2, Default: nil},
		{Keys: []string{"help", "h"}, AfterCount: 0, Default: false}, // Add this
	})

	if results[0].(bool) {
		getDeviceToUser()
	} else if startAutoclicker, ok := results[1].([]string); ok {
		runAutoclicker(startAutoclicker[0], startAutoclicker[1])
	} else {
		fmt.Println("No flags provided. Use --help for options.")
	}
}
