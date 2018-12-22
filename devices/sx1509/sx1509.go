// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.
//-----------------------------------------------------------------------------

package sx1509

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/mmr"
)

//-----------------------------------------------------------------------------
// exported symbols

// sx1509 register addresses.
const (
	RegInputDisableB    = 0x00 // Input buffer disable register _ I/O[15_8] (Bank B) 0000 0000
	RegInputDisableA    = 0x01 // Input buffer disable register _ I/O[7_0] (Bank A) 0000 0000
	RegLongSlewB        = 0x02 // Output buffer long slew register _ I/O[15_8] (Bank B) 0000 0000
	RegLongSlewA        = 0x03 // Output buffer long slew register _ I/O[7_0] (Bank A) 0000 0000
	RegLowDriveB        = 0x04 // Output buffer low drive register _ I/O[15_8] (Bank B) 0000 0000
	RegLowDriveA        = 0x05 // Output buffer low drive register _ I/O[7_0] (Bank A) 0000 0000
	RegPullUpB          = 0x06 // Pull_up register _ I/O[15_8] (Bank B) 0000 0000
	RegPullUpA          = 0x07 // Pull_up register _ I/O[7_0] (Bank A) 0000 0000
	RegPullDownB        = 0x08 // Pull_down register _ I/O[15_8] (Bank B) 0000 0000
	RegPullDownA        = 0x09 // Pull_down register _ I/O[7_0] (Bank A) 0000 0000
	RegOpenDrainB       = 0x0A // Open drain register _ I/O[15_8] (Bank B) 0000 0000
	RegOpenDrainA       = 0x0B // Open drain register _ I/O[7_0] (Bank A) 0000 0000
	RegPolarityB        = 0x0C // Polarity register _ I/O[15_8] (Bank B) 0000 0000
	RegPolarityA        = 0x0D // Polarity register _ I/O[7_0] (Bank A) 0000 0000
	RegDirB             = 0x0E // Direction register _ I/O[15_8] (Bank B) 1111 1111
	RegDirA             = 0x0F // Direction register _ I/O[7_0] (Bank A) 1111 1111
	RegDataB            = 0x10 // Data register _ I/O[15_8] (Bank B) 1111 1111*
	RegDataA            = 0x11 // Data register _ I/O[7_0] (Bank A) 1111 1111*
	RegInterruptMaskB   = 0x12 // Interrupt mask register _ I/O[15_8] (Bank B) 1111 1111
	RegInterruptMaskA   = 0x13 // Interrupt mask register _ I/O[7_0] (Bank A) 1111 1111
	RegSenseHighB       = 0x14 // Sense register for I/O[15:12] 0000 0000
	RegSenseLowB        = 0x15 // Sense register for I/O[11:8] 0000 0000
	RegSenseHighA       = 0x16 // Sense register for I/O[7:4] 0000 0000
	RegSenseLowA        = 0x17 // Sense register for I/O[3:0] 0000 0000
	RegInterruptSourceB = 0x18 // Interrupt source register _ I/O[15_8] (Bank B) 0000 0000
	RegInterruptSourceA = 0x19 // Interrupt source register _ I/O[7_0] (Bank A) 0000 0000
	RegEventStatusB     = 0x1A // Event status register _ I/O[15_8] (Bank B) 0000 0000
	RegEventStatusA     = 0x1B // Event status register _ I/O[7_0] (Bank A) 0000 0000
	RegLevelShifter1    = 0x1C // Level shifter register 0000 0000
	RegLevelShifter2    = 0x1D // Level shifter register 0000 0000
	RegClock            = 0x1E // Clock management register 0000 0000
	RegMisc             = 0x1F // Miscellaneous device settings register 0000 0000
	RegLEDDriverEnableB = 0x20 // LED driver enable register _ I/O[15_8] (Bank B) 0000 0000
	RegLEDDriverEnableA = 0x21 // LED driver enable register _ I/O[7_0] (Bank A) 0000 0000
	RegDebounceConfig   = 0x22 // Debounce configuration register 0000 0000
	RegDebounceEnableB  = 0x23 // Debounce enable register _ I/O[15_8] (Bank B) 0000 0000
	RegDebounceEnableA  = 0x24 // Debounce enable register _ I/O[7_0] (Bank A) 0000 0000
	RegKeyConfig1       = 0x25 // Key scan configuration register 0000 0000
	RegKeyConfig2       = 0x26 // Key scan configuration register 0000 0000
	RegKeyData1         = 0x27 // Key value (column) 1111 1111
	RegKeyData2         = 0x28 // Key value (row) 1111 1111
	RegTOn0             = 0x29 // ON time register for I/O[0] 0000 0000
	RegIOn0             = 0x2A // ON intensity register for I/O[0] 1111 1111
	RegOff0             = 0x2B // OFF time/intensity register for I/O[0] 0000 0000
	RegTOn1             = 0x2C // ON time register for I/O[1] 0000 0000
	RegIOn1             = 0x2D // ON intensity register for I/O[1] 1111 1111
	RegOff1             = 0x2E // OFF time/intensity register for I/O[1] 0000 0000
	RegTOn2             = 0x2F // ON time register for I/O[2] 0000 0000
	RegIOn2             = 0x30 // ON intensity register for I/O[2] 1111 1111
	RegOff2             = 0x31 // OFF time/intensity register for I/O[2] 0000 0000
	RegTOn3             = 0x32 // ON time register for I/O[3] 0000 0000
	RegIOn3             = 0x33 // ON intensity register for I/O[3] 1111 1111
	RegOff3             = 0x34 // OFF time/intensity register for I/O[3] 0000 0000
	RegTOn4             = 0x35 // ON time register for I/O[4] 0000 0000
	RegIOn4             = 0x36 // ON intensity register for I/O[4] 1111 1111
	RegOff4             = 0x37 // OFF time/intensity register for I/O[4] 0000 0000
	RegTRise4           = 0x38 // Fade in register for I/O[4] 0000 0000
	RegTFall4           = 0x39 // Fade out register for I/O[4] 0000 0000
	RegTOn5             = 0x3A // ON time register for I/O[5] 0000 0000
	RegIOn5             = 0x3B // ON intensity register for I/O[5] 1111 1111
	RegOff5             = 0x3C // OFF time/intensity register for I/O[5] 0000 0000
	RegTRise5           = 0x3D // Fade in register for I/O[5] 0000 0000
	RegTFall5           = 0x3E // Fade out register for I/O[5] 0000 0000
	RegTOn6             = 0x3F // ON time register for I/O[6] 0000 0000
	RegIOn6             = 0x40 // ON intensity register for I/O[6] 1111 1111
	RegOff6             = 0x41 // OFF time/intensity register for I/O[6] 0000 0000
	RegTRise6           = 0x42 // Fade in register for I/O[6] 0000 0000
	RegTFall6           = 0x43 // Fade out register for I/O[6] 0000 0000
	RegTOn7             = 0x44 // ON time register for I/O[7] 0000 0000
	RegIOn7             = 0x45 // ON intensity register for I/O[7] 1111 1111
	RegOff7             = 0x46 // OFF time/intensity register for I/O[7] 0000 0000
	RegTRise7           = 0x47 // Fade in register for I/O[7] 0000 0000
	RegTFall7           = 0x48 // Fade out register for I/O[7] 0000 0000
	RegTOn8             = 0x49 // ON time register for I/O[8] 0000 0000
	RegIOn8             = 0x4A // ON intensity register for I/O[8] 1111 1111
	RegOff8             = 0x4B // OFF time/intensity register for I/O[8] 0000 0000
	RegTOn9             = 0x4C // ON time register for I/O[9] 0000 0000
	RegIOn9             = 0x4D // ON intensity register for I/O[9] 1111 1111
	RegOff9             = 0x4E // OFF time/intensity register for I/O[9] 0000 0000
	RegTOn10            = 0x4F // ON time register for I/O[10] 0000 0000
	RegIOn10            = 0x50 // ON intensity register for I/O[10] 1111 1111
	RegOff10            = 0x51 // OFF time/intensity register for I/O[10] 0000 0000
	RegTOn11            = 0x52 // ON time register for I/O[11] 0000 0000
	RegIOn11            = 0x53 // ON intensity register for I/O[11] 1111 1111
	RegOff11            = 0x54 // OFF time/intensity register for I/O[11] 0000 0000
	RegTOn12            = 0x55 // ON time register for I/O[12] 0000 0000
	RegIOn12            = 0x56 // ON intensity register for I/O[12] 1111 1111
	RegOff12            = 0x57 // OFF time/intensity register for I/O[12] 0000 0000
	RegTRise12          = 0x58 // Fade in register for I/O[12] 0000 0000
	RegTFall12          = 0x59 // Fade out register for I/O[12] 0000 0000
	RegTOn13            = 0x5A // ON time register for I/O[13] 0000 0000
	RegIOn13            = 0x5B // ON intensity register for I/O[13] 1111 1111
	RegOff13            = 0x5C // OFF time/intensity register for I/O[13] 0000 0000
	RegTRise13          = 0x5D // Fade in register for I/O[13] 0000 0000
	RegTFall13          = 0x5E // Fade out register for I/O[13] 0000 0000
	RegTOn14            = 0x5F // ON time register for I/O[14] 0000 0000
	RegIOn14            = 0x60 // ON intensity register for I/O[14] 1111 1111
	RegOff14            = 0x61 // OFF time/intensity register for I/O[14] 0000 0000
	RegTRise14          = 0x62 // Fade in register for I/O[14] 0000 0000
	RegTFall14          = 0x63 // Fade out register for I/O[14] 0000 0000
	RegTOn15            = 0x64 // ON time register for I/O[15] 0000 0000
	RegIOn15            = 0x65 // ON intensity register for I/O[15] 1111 1111
	RegOff15            = 0x66 // OFF time/intensity register for I/O[15] 0000 0000
	RegTRise15          = 0x67 // Fade in register for I/O[15] 0000 0000
	RegTFall15          = 0x68 // Fade out register for I/O[15] 0000 0000
	RegHighInputB       = 0x69 // High input enable register _ I/O[15_8] (Bank B) 0000 0000
	RegHighInputA       = 0x6A // High input enable register _ I/O[7_0] (Bank A) 0000 0000
	RegReset            = 0x7D // Software reset register 0000 0000
	RegTest1            = 0x7E // Test register 0000 0000
	RegTest2            = 0x7F // Test register 0000 0000
)

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
	val0, err := d.c.ReadUint8(RegInterruptMaskA)
	if err != nil {
		return nil, err
	}
	val1, err := d.c.ReadUint8(RegSenseHighB)
	if err != nil {
		return nil, err
	}
	if val0 != 0xff || val1 != 0 {
		return nil, errors.New("bad register values")
	}
	// apply user provided register initialisation
	for i := range opts.Init {
		err := d.c.WriteUint8(opts.Init[i].Reg, opts.Init[i].Val)
		if err != nil {
			return nil, errors.New("can't initialise register")
		}
	}
	return d, nil
}

//-----------------------------------------------------------------------------

const SX1509_MAX_ROWS = 8 // maximum key scan rows
const SX1509_MAX_COLS = 8 // maximum key scan columns
const SX1509_DEBOUNCE_COUNT = 2

// Dev is the device object.
type Dev struct {
	c    mmr.Dev8
	opts Opts

	sample [SX1509_DEBOUNCE_COUNT]uint64 // debounce buffer for key samples
	keys   uint64                        // current debounced key state
	idx    uint                          // buffer index
	row    uint                          // current scan row

}

func (d *Dev) String() string {
	return fmt.Sprintf("sx1509{%s}", d.c.Conn)
}

// Halt the device.
func (d *Dev) Halt() error {
	return nil
}

// Poll the device.
func (d *Dev) Poll() {
	// read the column bits
	col, _ := d.c.ReadUint8(RegDataB)
	// add it to the sample buffer
	d.sample[d.idx] &= ^(uint64(0xff) << (d.row << 3))
	d.sample[d.idx] |= uint64(col^0xff) << (d.row << 3)
	// work out the current key state
	var keys uint64
	for i := 0; i < SX1509_DEBOUNCE_COUNT; i++ {
		keys |= d.sample[i]
	}
	// has it changed?
	if keys != d.keys {
		keyEvents(keys & ^d.keys, "dn")
		keyEvents(^keys&d.keys, "up")
		d.keys = keys
	}
	// increment/wrap the row index
	d.row++
	if d.row == SX1509_MAX_ROWS {
		// back to the 0th row
		d.row = 0
		// increment/wrap the debounce buffer index
		d.idx++
		if d.idx == SX1509_DEBOUNCE_COUNT {
			d.idx = 0
		}
	}
	// write the row selection bits
	d.c.WriteUint8(RegDataA, ^(1 << d.row))
}

// reset the device.
func (d *Dev) reset() error {
	if err := d.c.WriteUint8(RegReset, 0x12); err != nil {
		return err
	}
	if err := d.c.WriteUint8(RegReset, 0x34); err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------
// Private support code

// successively convert the multiple key bits to 0..63
func getKey(val *uint64) int {
	if *val == 0 {
		return -1
	}
	key := bits.TrailingZeros64(*val)
	*val &= ^(1 << uint(key))
	return key
}

// generate key events
func keyEvents(bits uint64, event string) {
	for key := getKey(&bits); key >= 0; key = getKey(&bits) {
		fmt.Printf("key %d event %s\n", key, event)
	}
}

//-----------------------------------------------------------------------------
