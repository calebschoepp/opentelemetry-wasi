package sdk

import (
	"sdk.com/internal/docs/adder/add"
)

func addInternal(s add.Numbers) uint32 {
	return s.X + s.Y
}

func AddExternal(x, y uint32) uint32 {
	s := add.Numbers{
		X: x,
		Y: y,
	}

	return addInternal(s)
}
