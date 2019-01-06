// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.
//-----------------------------------------------------------------------------

package rei2c

import (
	"errors"
)

//-----------------------------------------------------------------------------

// Opts specifies the rei2c configuration options.
type Opts struct {
	// I2CAddr is the I²C slave address to use.
	// Solderable links on the board allow the user to specify an arbitrary 7-bit address.
	I2CAddr uint16
	// RGB enables an illuminated RGB rotary encoder.
	RGB bool
	// DoublePush is the time (in 10ms increments) for the encoder double push. 0 is disabled.
	DoublePush uint8
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
