package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	pin := gpio.NewDirectPinDriver(r, "7")

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case val := <-ticker.C:
			pinVal, err := pin.DigitalRead()
			must(err)
			fmt.Printf("%+v\n", pinVal)
			fmt.Println(val)

		}
	}
}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
