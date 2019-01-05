// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.
//-----------------------------------------------------------------------------

package rei2c

import (
	"errors"
)

//-----------------------------------------------------------------------------

// RegInit specifies the initial value for a device register.
type RegInit struct {
	Reg uint8
	Val uint8
}

// Opts specifies the sx1509 configuration options.
type Opts struct {
	// I2CAddr is the I²C slave address to use.
	// Solderable links on the board allow the user to specify an arbitrary 7-bit address.
	I2CAddr uint16
	// Init is an optional set of initial register values.
	Init []RegInit
}

// DefaultOpts contains the default options to use.
var DefaultOpts = Opts{}

const DefaultAddress = 0x55

//-----------------------------------------------------------------------------

func (o *Opts) i2cAddr() (uint16, error) {
	if o.I2CAddr == 0 {
		return DefaultAddress, nil // default
	}
	// reserved address
	if (o.I2CAddr <= 7) || (o.I2CAddr&0x78 == 0x78) {
		return 0, errors.New("i2c address not supported by device")
	}
	return o.I2CAddr, nil
}

//-----------------------------------------------------------------------------
