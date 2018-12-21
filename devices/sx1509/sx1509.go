// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.
//-----------------------------------------------------------------------------

package sx1509

import (
	"encoding/binary"
	"errors"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/mmr"
)

//-----------------------------------------------------------------------------
// exported support symbols

//-----------------------------------------------------------------------------

// New returns the Dev object for an sx1509 on an I2C bus.
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

//-----------------------------------------------------------------------------
// Dev struct and methods

// Dev is a handle to the device.
type Dev struct {
	c    mmr.Dev8
	opts Opts
}

// Read is a method on your device.
func (d *Dev) Read() string {
	return ""
}

//-----------------------------------------------------------------------------
// Private support code

const (
	regInputDisableB    = 0x00 // Input buffer disable register _ I/O[15_8] (Bank B) 0000 0000
	regInputDisableA    = 0x01 // Input buffer disable register _ I/O[7_0] (Bank A) 0000 0000
	regLongSlewB        = 0x02 // Output buffer long slew register _ I/O[15_8] (Bank B) 0000 0000
	regLongSlewA        = 0x03 // Output buffer long slew register _ I/O[7_0] (Bank A) 0000 0000
	regLowDriveB        = 0x04 // Output buffer low drive register _ I/O[15_8] (Bank B) 0000 0000
	regLowDriveA        = 0x05 // Output buffer low drive register _ I/O[7_0] (Bank A) 0000 0000
	regPullUpB          = 0x06 // Pull_up register _ I/O[15_8] (Bank B) 0000 0000
	regPullUpA          = 0x07 // Pull_up register _ I/O[7_0] (Bank A) 0000 0000
	regPullDownB        = 0x08 // Pull_down register _ I/O[15_8] (Bank B) 0000 0000
	regPullDownA        = 0x09 // Pull_down register _ I/O[7_0] (Bank A) 0000 0000
	regOpenDrainB       = 0x0A // Open drain register _ I/O[15_8] (Bank B) 0000 0000
	regOpenDrainA       = 0x0B // Open drain register _ I/O[7_0] (Bank A) 0000 0000
	regPolarityB        = 0x0C // Polarity register _ I/O[15_8] (Bank B) 0000 0000
	regPolarityA        = 0x0D // Polarity register _ I/O[7_0] (Bank A) 0000 0000
	regDirB             = 0x0E // Direction register _ I/O[15_8] (Bank B) 1111 1111
	regDirA             = 0x0F // Direction register _ I/O[7_0] (Bank A) 1111 1111
	regDataB            = 0x10 // Data register _ I/O[15_8] (Bank B) 1111 1111*
	regDataA            = 0x11 // Data register _ I/O[7_0] (Bank A) 1111 1111*
	regInterruptMaskB   = 0x12 // Interrupt mask register _ I/O[15_8] (Bank B) 1111 1111
	regInterruptMaskA   = 0x13 // Interrupt mask register _ I/O[7_0] (Bank A) 1111 1111
	regSenseHighB       = 0x14 // Sense register for I/O[15:12] 0000 0000
	regSenseLowB        = 0x15 // Sense register for I/O[11:8] 0000 0000
	regSenseHighA       = 0x16 // Sense register for I/O[7:4] 0000 0000
	regSenseLowA        = 0x17 // Sense register for I/O[3:0] 0000 0000
	regInterruptSourceB = 0x18 // Interrupt source register _ I/O[15_8] (Bank B) 0000 0000
	regInterruptSourceA = 0x19 // Interrupt source register _ I/O[7_0] (Bank A) 0000 0000
	regEventStatusB     = 0x1A // Event status register _ I/O[15_8] (Bank B) 0000 0000
	regEventStatusA     = 0x1B // Event status register _ I/O[7_0] (Bank A) 0000 0000
	regLevelShifter1    = 0x1C // Level shifter register 0000 0000
	regLevelShifter2    = 0x1D // Level shifter register 0000 0000
	regClock            = 0x1E // Clock management register 0000 0000
	regMisc             = 0x1F // Miscellaneous device settings register 0000 0000
	regLEDDriverEnableB = 0x20 // LED driver enable register _ I/O[15_8] (Bank B) 0000 0000
	regLEDDriverEnableA = 0x21 // LED driver enable register _ I/O[7_0] (Bank A) 0000 0000
	regDebounceConfig   = 0x22 // Debounce configuration register 0000 0000
	regDebounceEnableB  = 0x23 // Debounce enable register _ I/O[15_8] (Bank B) 0000 0000
	regDebounceEnableA  = 0x24 // Debounce enable register _ I/O[7_0] (Bank A) 0000 0000
	regKeyConfig1       = 0x25 // Key scan configuration register 0000 0000
	regKeyConfig2       = 0x26 // Key scan configuration register 0000 0000
	regKeyData1         = 0x27 // Key value (column) 1111 1111
	regKeyData2         = 0x28 // Key value (row) 1111 1111
	regTOn0             = 0x29 // ON time register for I/O[0] 0000 0000
	regIOn0             = 0x2A // ON intensity register for I/O[0] 1111 1111
	regOff0             = 0x2B // OFF time/intensity register for I/O[0] 0000 0000
	regTOn1             = 0x2C // ON time register for I/O[1] 0000 0000
	regIOn1             = 0x2D // ON intensity register for I/O[1] 1111 1111
	regOff1             = 0x2E // OFF time/intensity register for I/O[1] 0000 0000
	regTOn2             = 0x2F // ON time register for I/O[2] 0000 0000
	regIOn2             = 0x30 // ON intensity register for I/O[2] 1111 1111
	regOff2             = 0x31 // OFF time/intensity register for I/O[2] 0000 0000
	regTOn3             = 0x32 // ON time register for I/O[3] 0000 0000
	regIOn3             = 0x33 // ON intensity register for I/O[3] 1111 1111
	regOff3             = 0x34 // OFF time/intensity register for I/O[3] 0000 0000
	regTOn4             = 0x35 // ON time register for I/O[4] 0000 0000
	regIOn4             = 0x36 // ON intensity register for I/O[4] 1111 1111
	regOff4             = 0x37 // OFF time/intensity register for I/O[4] 0000 0000
	regTRise4           = 0x38 // Fade in register for I/O[4] 0000 0000
	regTFall4           = 0x39 // Fade out register for I/O[4] 0000 0000
	regTOn5             = 0x3A // ON time register for I/O[5] 0000 0000
	regIOn5             = 0x3B // ON intensity register for I/O[5] 1111 1111
	regOff5             = 0x3C // OFF time/intensity register for I/O[5] 0000 0000
	regTRise5           = 0x3D // Fade in register for I/O[5] 0000 0000
	regTFall5           = 0x3E // Fade out register for I/O[5] 0000 0000
	regTOn6             = 0x3F // ON time register for I/O[6] 0000 0000
	regIOn6             = 0x40 // ON intensity register for I/O[6] 1111 1111
	regOff6             = 0x41 // OFF time/intensity register for I/O[6] 0000 0000
	regTRise6           = 0x42 // Fade in register for I/O[6] 0000 0000
	regTFall6           = 0x43 // Fade out register for I/O[6] 0000 0000
	regTOn7             = 0x44 // ON time register for I/O[7] 0000 0000
	regIOn7             = 0x45 // ON intensity register for I/O[7] 1111 1111
	regOff7             = 0x46 // OFF time/intensity register for I/O[7] 0000 0000
	regTRise7           = 0x47 // Fade in register for I/O[7] 0000 0000
	regTFall7           = 0x48 // Fade out register for I/O[7] 0000 0000
	regTOn8             = 0x49 // ON time register for I/O[8] 0000 0000
	regIOn8             = 0x4A // ON intensity register for I/O[8] 1111 1111
	regOff8             = 0x4B // OFF time/intensity register for I/O[8] 0000 0000
	regTOn9             = 0x4C // ON time register for I/O[9] 0000 0000
	regIOn9             = 0x4D // ON intensity register for I/O[9] 1111 1111
	regOff9             = 0x4E // OFF time/intensity register for I/O[9] 0000 0000
	regTOn10            = 0x4F // ON time register for I/O[10] 0000 0000
	regIOn10            = 0x50 // ON intensity register for I/O[10] 1111 1111
	regOff10            = 0x51 // OFF time/intensity register for I/O[10] 0000 0000
	regTOn11            = 0x52 // ON time register for I/O[11] 0000 0000
	regIOn11            = 0x53 // ON intensity register for I/O[11] 1111 1111
	regOff11            = 0x54 // OFF time/intensity register for I/O[11] 0000 0000
	regTOn12            = 0x55 // ON time register for I/O[12] 0000 0000
	regIOn12            = 0x56 // ON intensity register for I/O[12] 1111 1111
	regOff12            = 0x57 // OFF time/intensity register for I/O[12] 0000 0000
	regTRise12          = 0x58 // Fade in register for I/O[12] 0000 0000
	regTFall12          = 0x59 // Fade out register for I/O[12] 0000 0000
	regTOn13            = 0x5A // ON time register for I/O[13] 0000 0000
	regIOn13            = 0x5B // ON intensity register for I/O[13] 1111 1111
	regOff13            = 0x5C // OFF time/intensity register for I/O[13] 0000 0000
	regTRise13          = 0x5D // Fade in register for I/O[13] 0000 0000
	regTFall13          = 0x5E // Fade out register for I/O[13] 0000 0000
	regTOn14            = 0x5F // ON time register for I/O[14] 0000 0000
	regIOn14            = 0x60 // ON intensity register for I/O[14] 1111 1111
	regOff14            = 0x61 // OFF time/intensity register for I/O[14] 0000 0000
	regTRise14          = 0x62 // Fade in register for I/O[14] 0000 0000
	regTFall14          = 0x63 // Fade out register for I/O[14] 0000 0000
	regTOn15            = 0x64 // ON time register for I/O[15] 0000 0000
	regIOn15            = 0x65 // ON intensity register for I/O[15] 1111 1111
	regOff15            = 0x66 // OFF time/intensity register for I/O[15] 0000 0000
	regTRise15          = 0x67 // Fade in register for I/O[15] 0000 0000
	regTFall15          = 0x68 // Fade out register for I/O[15] 0000 0000
	regHighInputB       = 0x69 // High input enable register _ I/O[15_8] (Bank B) 0000 0000
	regHighInputA       = 0x6A // High input enable register _ I/O[7_0] (Bank A) 0000 0000
	regReset            = 0x7D // Software reset register 0000 0000
	regTest1            = 0x7E // Test register 0000 0000
	regTest2            = 0x7F // Test register 0000 0000
)

//-----------------------------------------------------------------------------

// reset the device.
func (d *Dev) reset() error {
	if err := d.c.WriteUint8(regReset, 0x12); err != nil {
		return err
	}
	if err := d.c.WriteUint8(regReset, 0x34); err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------

func makeDev(c conn.Conn, opts *Opts) (*Dev, error) {
	d := &Dev{
		opts: *opts,
		c:    mmr.Dev8{Conn: c, Order: binary.LittleEndian},
	}

	// reset the device
	if err := d.reset(); err != nil {
		return nil, err
	}

	// check the expected default values for some registers
	val0, err := d.c.ReadUint8(regInterruptMaskA)
	if err != nil {
		return nil, err
	}
	val1, err := d.c.ReadUint8(regSenseHighB)
	if err != nil {
		return nil, err
	}
	if val0 != 0xff || val1 != 0 {
		return nil, errors.New("bad register values")
	}

	return d, nil
}

//-----------------------------------------------------------------------------
