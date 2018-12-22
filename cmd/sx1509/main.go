// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.
//-----------------------------------------------------------------------------

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/deadsy/pdev/devices/sx1509"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/pin"
	"periph.io/x/periph/conn/pin/pinreg"
	"periph.io/x/periph/host"
)

//-----------------------------------------------------------------------------

func printPin(fn string, p pin.Pin) {
	name, pos := pinreg.Position(p)
	if name != "" {
		log.Printf("  %-4s: %-10s found on header %s, #%d", fn, p, name, pos)
	} else {
		log.Printf("  %-4s: %-10s", fn, p)
	}
}

//-----------------------------------------------------------------------------

func mainImpl() error {

	busId := flag.String("bus", "", "I²C bus")
	devAddr := flag.Uint("adr", 0x3e, "I²C device address")
	busSpeed := flag.Int("hz", 0, "I²C bus speed")

	flag.Parse()

	opts := sx1509.DefaultOpts

	if *devAddr != 0 {
		if *devAddr < 0 || *devAddr > 0x7f {
			return errors.New("invalid i2c device address")
		}
		opts.I2CAddr = uint16(*devAddr)
	}

	if _, err := host.Init(); err != nil {
		return err
	}

	i2cBus, err := i2creg.Open(*busId)
	if err != nil {
		return fmt.Errorf("couldn't open i2c bus: %s", err)
	}
	defer i2cBus.Close()

	if p, ok := i2cBus.(i2c.Pins); ok {
		printPin("SCL", p.SCL())
		printPin("SDA", p.SDA())
	}

	if *busSpeed != 0 {
		if err := i2cBus.SetSpeed(physic.Frequency(*busSpeed) * physic.Hertz); err != nil {
			return fmt.Errorf("couldn't set i2c bus speed: %s", err)
		}
	}

	// TODO device interrupt pin

	opts.Init = []sx1509.RegInit{
		{sx1509.RegClock, 0x50},
		{sx1509.RegMisc, 0x10},
		{sx1509.RegDirA, 0x00},
		{sx1509.RegOpenDrainA, 0xff},
		{sx1509.RegPullUpB, 0xff},
	}

	dev, err := sx1509.New(i2cBus, &opts)
	if err != nil {
		return fmt.Errorf("couldn't open sx1509: %s", err)
	}

	fmt.Printf("%s\n", dev)

	for {
		dev.Poll()
		time.Sleep(4 * time.Millisecond)
	}

	return dev.Halt()
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "sx1509: %s\n", err)
		os.Exit(1)
	}
}

//-----------------------------------------------------------------------------
