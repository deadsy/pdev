// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.
//-----------------------------------------------------------------------------

// Package rei2c is a driver for the I2C Rotary Encoder V2 Driver
//
// Author
// Jason Harris (https://github.com/deadsy)
//
// Datasheet
// https://github.com/Fattoresaimon/I2CEncoderV2
//
// Product Page
// https://www.kickstarter.com/projects/1351830006/i2c-encoder-v2
package rei2c

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/mmr"
)

//-----------------------------------------------------------------------------
// exported symbols

// rei2c register addresses
const (
	RegGCONF    = 0x00 // General Configuration (1 byte)
	RegGP1CONF  = 0x01 // GP 1 Configuration (1 byte)
	RegGP2CONF  = 0x02 // GP 2 Configuration (1 byte)
	RegGP3CONF  = 0x03 // GP 3 Configuration (1 byte)
	RegINTCONF  = 0x04 // INT pin Configuration (1 byte)
	RegESTATUS  = 0x05 // Encoder Status (1 byte)
	RegI2STATUS = 0x06 // Secondary interrupt status (1 byte)
	RegFSTATUS  = 0x07 // Fade process status (1 byte)
	RegCVAL     = 0x08 // Counter Value (4 bytes)
	RegCMAX     = 0x0C // Counter Max value (4 bytes)
	RegCMIN     = 0x10 // Counter Min value (4 bytes)
	RegISTEP    = 0x14 // Increment step value (4 bytes)
	RegRLED     = 0x18 // LED red color intensity (1 byte)
	RegGLED     = 0x19 // LED green color intensity (1 byte)
	RegBLED     = 0x1A // LED blue color intensity (1 byte)
	RegGP1REG   = 0x1B // I/O GP1 Register (1 byte)
	RegGP2REG   = 0x1C // I/O GP2 Register (1 byte)
	RegGP3REG   = 0x1D // I/O GP3 Register (1 byte)
	RegANTBOUNC = 0x1E // Anti-bouncing period (1 Byte)
	RegDPPERIOD = 0x1F // Double push period (1 Byte)
	RegFADERGB  = 0x20 // Fade timer RGB Encoder (1 Byte)
	RegFADEGP   = 0x21 // Fade timer GP ports (1 Byte)
	RegEEPROM   = 0x80 // EEPROM memory (128 bytes)
)

// gconf bits
const (
	gconfDTYPE = uint8(1 << 0) // Data type of the register: CVAL, CMAX, CMIN and ISTEP.
	gconfWRAPE = uint8(1 << 1) // Enable counter wrap.
	gconfDIRE  = uint8(1 << 2) // Direction of the encoder when increment.
	gconfIPUD  = uint8(1 << 3) // Interrupt Pull-UP disable.
	gconfRMOD  = uint8(1 << 4) // Reading Mode.
	gconfETYPE = uint8(1 << 5) // Set the encoder type (normal/illuminated)
	gconfMBANK = uint8(1 << 6) // Select the EEPROM memory bank. Each bank is 128 bytes.
	gconfRESET = uint8(1 << 7) // Reset the I2C Encoder V2
)

// status bits
const (
	statusPUSHR = uint8(1 << 0) // push button has been released
	statusPUSHP = uint8(1 << 1) // push button has been pressed
	statusPUSHD = uint8(1 << 2) // push button has been double pushed
	statusRINC  = uint8(1 << 3) // rotated in the increase direction
	statusRDEC  = uint8(1 << 4) // rotated in the decrease direction
	statusRMAX  = uint8(1 << 5) // maximum counter value has been reached
	statusRMIN  = uint8(1 << 6) // minimum counter value has been reached
	statusINT2  = uint8(1 << 7) // Secondary interrupt status
)

//-----------------------------------------------------------------------------

// New returns the Dev object for an rei2c on an I2C bus.
func New(b i2c.Bus, opts *Opts) (*Dev, error) {
	if opts == nil {
		opts = &DefaultOpts
	}
	addr, err := opts.i2cAddr()
	if err != nil {
		return nil, err
	}
	d, err := makeDev(&i2c.Dev{Bus: b, Addr: addr}, opts)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func makeDev(c conn.Conn, opts *Opts) (*Dev, error) {
	d := &Dev{
		opts: *opts,
		c:    mmr.Dev8{Conn: c, Order: binary.BigEndian},
	}
	// reset the device
	if err := d.reset(); err != nil {
		return nil, err
	}

	// check the expected default values for some registers
	val0, err := d.c.ReadUint8(RegGP1CONF)
	if err != nil {
		return nil, err
	}
	val1, err := d.c.ReadUint8(RegANTBOUNC)
	if err != nil {
		return nil, err
	}
	if val0 != 0 || val1 != 25 {
		return nil, errors.New("bad register values")
	}

	// setup the general configuration
	var gconf uint8
	if opts.RGB {
		// enable an illuminated RGB encoder.
		gconf |= gconfETYPE
	}
	err = d.wrGCONF(gconf)
	if err != nil {
		return nil, err
	}

	// setup the double push time
	if opts.DoublePush != 0 {
		err := d.c.WriteUint8(RegDPPERIOD, opts.DoublePush)
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}

//-----------------------------------------------------------------------------

// Dev is the device object.
type Dev struct {
	c    mmr.Dev8
	opts Opts

	gconf uint8
}

func (d *Dev) String() string {
	return fmt.Sprintf("rei2c{%s}", d.c.Conn)
}

// Halt the device.
func (d *Dev) Halt() error {
	return nil
}

// Poll the device.
func (d *Dev) Poll() {
	status, _ := d.rdESTATUS()
	if status&(statusRINC|statusRDEC) != 0 {
		n, _ := d.RdCntVal()
		fmt.Printf("count %d\n", n)
	}
	if status&statusPUSHR != 0 {
		fmt.Printf("pushr\n")
	}
	if status&statusPUSHP != 0 {
		fmt.Printf("pushp\n")
	}
	if status&statusPUSHD != 0 {
		fmt.Printf("pushd\n")
	}
}

//-----------------------------------------------------------------------------

/*

// RdMem reads from the rei2c EEPROM.
func (d *Dev) RdMem(addr uint8) (uint8, error) {
	var err error
	if addr <= 0x7f {
		// switch to bank 0
		if d.gconf&gconfMBANK != 0 {
			err = d.wrGCONF(d.gconf & ^gconfMBANK)
		}
		addr += RegEEPROM
	} else {
		// switch to bank 1
		if d.gconf&gconfMBANK == 0 {
			err = d.wrGCONF(d.gconf | gconfMBANK)
		}
	}
	if err != nil {
		return 0, err
	}
	val, err := d.c.ReadUint8(addr)
	if err != nil {
		return 0, err
	}
	time.Sleep(1 * time.Millisecond)
	return val, nil
}

*/

func (d *Dev) setBank(gconf, adr uint8) (uint8, uint8, error) {

	x := gconf

	if adr <= 0x7f {
		// switch to bank 0
		if x&gconfMBANK != 0 {
			x &= ^gconfMBANK
		}
		adr += RegEEPROM
	} else {
		// switch to bank 1
		if x&gconfMBANK == 0 {
			x |= gconfMBANK
		}
	}

	var err error
	if x != gconf {
		err = d.c.WriteUint8(RegGCONF, x)
	}

	return x, adr, err
}

// RdMem reads from the rei2c EEPROM.
func (d *Dev) RdMem(base uint8, n int) ([]uint8, error) {
	if n == 0 {
		return nil, nil
	}

	gconf, err := d.c.ReadUint8(RegGCONF)
	if err != nil {
		return nil, err
	}

	mem := make([]uint8, n)
	for i := range mem {
		adr := base + uint8(i)
		x, adr, err := d.setBank(gconf, adr)
		if err != nil {
			return nil, err
		}

		gconf = x

		val, err := d.c.ReadUint8(adr)
		if err != nil {
			return nil, err
		}
		mem[i] = val
		time.Sleep(1 * time.Millisecond)
	}
	return mem, nil
}

//-----------------------------------------------------------------------------
// LED control

type RGB struct {
	R, G, B uint8
}

func (d *Dev) RdLED() (RGB, error) {
	r, _ := d.c.ReadUint8(RegRLED)
	g, _ := d.c.ReadUint8(RegGLED)
	b, _ := d.c.ReadUint8(RegBLED)
	return RGB{r, g, b}, nil
	/*
		var rgb RGB
		err := d.c.ReadStruct(RegRLED, &rgb)
		if err != nil {
			return RGB{}, err
		}
		return rgb, nil
	*/
}

func (d *Dev) WrLED(rgb RGB) error {
	d.c.WriteUint8(RegRLED, rgb.R)
	d.c.WriteUint8(RegGLED, rgb.G)
	d.c.WriteUint8(RegBLED, rgb.B)
	return nil
	//return d.c.WriteStruct(RegRLED, &rgb)
}

//-----------------------------------------------------------------------------
// counter control

func (d *Dev) RdCntMin() (uint32, error) {
	return d.c.ReadUint32(RegCMIN)
}

func (d *Dev) RdCntMax() (uint32, error) {
	return d.c.ReadUint32(RegCMAX)
}

func (d *Dev) RdCntVal() (uint32, error) {
	n, err := d.c.ReadUint8(RegCVAL + 3)
	return uint32(n), err
	//return d.c.ReadUint32(RegCVAL)
}

func (d *Dev) RdCntStep() (uint32, error) {
	return d.c.ReadUint32(RegISTEP)
}

func (d *Dev) WrCntMin(n uint32) error {
	return d.c.WriteUint32(RegCMIN, n)
}

func (d *Dev) WrCntMax(n uint32) error {
	return d.c.WriteUint32(RegCMAX, n)
}

func (d *Dev) WrCntVal(n uint32) error {
	return d.c.WriteUint32(RegCVAL, n)
}

func (d *Dev) WrCntStep(n uint32) error {
	return d.c.WriteUint32(RegISTEP, n)
}

//-----------------------------------------------------------------------------

// rdESTATUS read the encoder status register
func (d *Dev) rdESTATUS() (uint8, error) {
	return d.c.ReadUint8(RegESTATUS)
}

// wrGCONF write the general configuration register.
func (d *Dev) wrGCONF(val uint8) error {
	err := d.c.WriteUint8(RegGCONF, val)
	if err != nil {
		return err
	}
	d.gconf = val
	return nil
}

// reset the device.
func (d *Dev) reset() error {
	err := d.wrGCONF(gconfRESET)
	if err != nil {
		return err
	}
	d.gconf = 0
	time.Sleep(10 * time.Millisecond)
	return nil
}

//-----------------------------------------------------------------------------
