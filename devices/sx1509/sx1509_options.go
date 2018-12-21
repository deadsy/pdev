// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.
//-----------------------------------------------------------------------------

package sx1509

import (
	"errors"
)

//-----------------------------------------------------------------------------

// Opts specifies the sx1509 configuration options.
type Opts struct {
	// I2CAddr is the I²C slave address to use.
	// Its default value is 0x3e.
	// It can be set to other values (0x3f, 0x70, 0x71) depending on HW configuration.
	I2CAddr uint16
}

// DefaultOpts contains the default options to use.
var DefaultOpts = Opts{}

//-----------------------------------------------------------------------------

func (o *Opts) i2cAddr() (uint16, error) {
	switch o.I2CAddr {
	case 0:
		return 0x3e, nil // default
	case 0x3e, 0x3f, 0x70, 0x71:
		return o.I2CAddr, nil
	default:
		return 0, errors.New("i2c address not supported by device")
	}
}

//-----------------------------------------------------------------------------
