package input

// StandardController
// http://wiki.nesdev.com/w/index.php/Standard_controller
type StandardController struct {
	state  uint8
	cur    uint8
	strobe bool
}

func (c *StandardController) Update(state uint8) {
	c.state = state
}

func (c *StandardController) Write(value uint8) {
	c.strobe = (value & 1) == 1
	c.cur = 1
}

func (c *StandardController) Read() uint8 {
	var v uint8
	if c.strobe {
		v = c.state & uint8(Standard_A)
	} else {
		input := c.state & c.cur
		c.cur = c.cur << 1
		if 0 < input {
			v = 1
		} else {
			v = 0
		}
	}
	return v | 0x40
}

type StandardControllerButton uint8

const (
	Standard_A      StandardControllerButton = 1 << 0
	Standard_B      StandardControllerButton = 1 << 1
	Standard_Select StandardControllerButton = 1 << 2
	Standard_Start  StandardControllerButton = 1 << 3
	Standard_Up     StandardControllerButton = 1 << 4
	Standard_Down   StandardControllerButton = 1 << 5
	Standard_Left   StandardControllerButton = 1 << 6
	Standard_Right  StandardControllerButton = 1 << 7
)
