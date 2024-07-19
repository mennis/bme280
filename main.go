package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
)

func main() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	refs := i2creg.All()
	if len(refs) == 0 {
		log.Print("No I²C buses available\n")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	// Take readings from the sensor every 2 seconds:
	for {
		select {
		case <-time.After(2 * time.Second):
			err := takeReading(bus)
			if err != nil {
				log.Printf("failed to take reading: %s", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func takeReading(bus i2c.BusCloser) error {
	// Open a handle to a bme280/bmp280 connected on the I²C bus using default
	// settings:
	dev, err := bmxx80.NewI2C(bus, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		return fmt.Errorf("failed to open BMxx80: %w", err)
	}
	defer dev.Halt()
	// Read temperature from the sensor:
	var env physic.Env
	if err = dev.Sense(&env); err != nil {
		return fmt.Errorf("failed to read environmental data: %w", err)
	}
	log.Printf("Temp: %8s\tHumidity: %9s\n", env.Temperature, env.Humidity)
	return nil
}

