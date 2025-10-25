package main

import (
	"fmt"
	"os"
	"time"

	"github.com/leandronowras/device-api/internal/device"
)

func main() {
	d, err := device.New("iPhone", "Apple")
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not create device: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Device created successfully!")
	fmt.Println()
	fmt.Printf("Name:   %s\n", d.Name())
	fmt.Printf("Brand:  %s\n", d.Brand())
	fmt.Printf("ID:     %s\n", d.ID())
	fmt.Printf("State:  %s\n", d.State())
	fmt.Printf("Created: %s\n", d.CreationTime().Format(time.RFC822))
}
