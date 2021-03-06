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

	"github.com/deadsy/pdev/devices/rei2c"
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
	devAddr := flag.Uint("adr", rei2c.DefaultAddress, "I²C device address")
	busSpeed := flag.Int("hz", 0, "I²C bus speed")

	flag.Parse()

	opts := rei2c.DefaultOpts

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

	opts.RGB = true
	opts.DoublePush = 50

	dev, err := rei2c.New(i2cBus, &opts)
	if err != nil {
		return fmt.Errorf("couldn't open rei2c: %s", err)
	}

	fmt.Printf("%s\n", dev)

	dev.WrLED(rei2c.RGB{0, 0, 255})
	x, err := dev.RdLED()
	fmt.Printf("%d %d %d %v\n", x.R, x.G, x.B, err)

	dev.WrCntMin(0)
	dev.WrCntMax(255)
	dev.WrCntVal(0)
	dev.WrCntStep(1)

	for {
		dev.Poll()
		time.Sleep(50 * time.Millisecond)
	}

	return dev.Halt()
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "rei2c: %s\n", err)
		os.Exit(1)
	}
}

//-----------------------------------------------------------------------------
