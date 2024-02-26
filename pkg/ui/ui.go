package ui

import (
	omap "github.com/wk8/go-ordered-map"
)

type KeyboardButton struct{
	Text string
}

type KeyboardMenu struct{
	Buttons omap.OrderedMap
}