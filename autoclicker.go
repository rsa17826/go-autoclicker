package main

import "fmt"

func runAutoclicker(mouse string, keyboard string) {
	file, err := getDeviceFromIdOrName(mouse)
	if err != nil {
		fmt.Printf("Error opening device: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Printf("Successfully hooked into: %s\n", file.Name())
}
