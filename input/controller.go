package input

// https://www.nesdev.org/wiki/Input_devices

// Controller represents IO for general-purpose controller ports from NES
type Controller interface {
	Write(value uint8)
	Read() uint8
}
