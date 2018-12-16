// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package rei2c

import (
	"errors"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/i2c"
)

// FIXME: Expose public symbols as relevant. Do not export more than needed!
// See https://periph.io/project/#requirements
// for the expectations.
//
// Use the following layout for drivers:
//  - exported support symbols
//  - Opts struct
//  - New func
//  - Dev struct and methods
//  - Private support code

// New opens a handle to the device. FIXME.
func New(i i2c.Bus) (*Dev, error) {
	d := &Dev{&i2c.Dev{Bus: i, Addr: 42}}
	// FIXME: Simulate a setup dance.
	var b [2]byte
	if err := d.c.Tx([]byte("in"), b[:]); err != nil {
		return nil, err
	}
	if b[0] != 'I' || b[1] != 'N' {
		return nil, errors.New("rei2c: unexpected reply")
	}
	return d, nil
}

// Dev is a handle to the device. FIXME.
type Dev struct {
	c conn.Conn
}

// Read is a method on your device. FIXME.
func (d *Dev) Read() string {
	var b [12]byte
	if err := d.c.Tx([]byte("what"), b[:]); err != nil {
		return err.Error()
	}
	return string(b[:])
}

// device registers
const (
	regGCONF    = 0x00 // General Configuration (1 byte)
	regGP1CONF  = 0x01 // GP 1 Configuration (1 byte)
	regGP2CONF  = 0x02 // GP 2 Configuration (1 byte)
	regGP3CONF  = 0x03 // GP 3 Configuration (1 byte)
	regINTCONF  = 0x04 // INT pin Configuration (1 byte)
	regESTATUS  = 0x05 // Encoder Status (1 byte)
	regI2STATUS = 0x06 // Secondary interrupt status (1 byte)
	regFSTATUS  = 0x07 // Fade process status (1 byte)
	regCVAL     = 0x08 // Counter Value (4 bytes)
	regCMAX     = 0x0C // Counter Max value (4 bytes)
	regCMIN     = 0x10 // Counter Min value (4 bytes)
	regISTEP    = 0x14 // Increment step value (4 bytes)
	regRLED     = 0x18 // LED red color intensity (1 byte)
	regGLED     = 0x19 // LED green color intensity (1 byte)
	regBLED     = 0x1A // LED blue color intensity (1 byte)
	regGP1REG   = 0x1B // I/O GP1 register (1 byte)
	regGP2REG   = 0x1C // I/O GP2 register (1 byte)
	regGP3REG   = 0x1D // I/O GP3 register (1 byte)
	regANTBOUNC = 0x1E // Anti-bouncing period (1 Byte)
	regDPPERIOD = 0x1F // Double push period (1 Byte)
	regFADERGB  = 0x20 // Fade timer RGB Encoder (1 Byte)
	regFADEGP   = 0x21 // Fade timer GP ports (1 Byte)
	regEEPROM   = 0x80 // EEPROM memory (128 bytes)
)
