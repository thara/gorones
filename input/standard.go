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
		v = c.state & StandardA
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

type StandardControllerButton = uint8

const (
	StandardA      StandardControllerButton = 1 << 0
	StandardB      StandardControllerButton = 1 << 1
	StandardSelect StandardControllerButton = 1 << 2
	StandardStart  StandardControllerButton = 1 << 3
	StandardUp     StandardControllerButton = 1 << 4
	StandardDown   StandardControllerButton = 1 << 5
	StandardLeft   StandardControllerButton = 1 << 6
	StandardRight  StandardControllerButton = 1 << 7
)
