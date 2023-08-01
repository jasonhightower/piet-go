package interpreter

import (
	"image/color"
)

const RED byte = 1
const YELLOW byte = 2
const GREEN byte = 3
const CYAN byte = 4
const BLUE byte = 5
const MAGENTA byte = 6

const LIGHT byte = 1
const NORMAL byte = 2
const DARK byte = 3

const low uint16 = 0x0000
const mid uint16 = 0xC0C0
const high uint16 = 0xFFFF

var l_red color.Color = color.RGBA64{R:high, G:mid, B:mid}
var n_red color.Color = color.RGBA64{R:high}
var d_red color.Color = color.RGBA64{R:mid}

var l_yellow color.Color = color.RGBA64{R:high, G:high, B:mid}
var n_yellow color.Color = color.RGBA64{R:high, G:high, B:low}
var d_yellow color.Color = color.RGBA64{R:mid, G:mid}

var l_green color.Color = color.RGBA64{R:mid, G:high, B:mid}
var n_green color.Color = color.RGBA64{G:high}
var d_green color.Color = color.RGBA64{G:mid}

var l_cyan color.Color = color.RGBA64{R:mid, G:high, B:high}
var n_cyan color.Color = color.RGBA64{G:high, B:high}
var d_cyan color.Color = color.RGBA64{G:mid, B:mid}

var l_blue color.Color = color.RGBA64{R:mid, G:mid, B:high}
var n_blue color.Color = color.RGBA64{B:high}
var d_blue color.Color = color.RGBA64{B:mid}

var l_magenta color.Color = color.RGBA64{R:high, G:mid, B:high}
var n_magenta color.Color = color.RGBA64{R:high, B:high}
var d_magenta color.Color = color.RGBA64{R:mid, B:mid}

func Hue(color color.Color) byte {
    r, g, b, _ := color.RGBA()
    if r > g && r > b {
        return RED
    } else if g > r && g > b {
        return GREEN
    } else if b > r && b > g {
        return BLUE
    } else if r == g && r > b {
        return YELLOW
    } else if g == b && g > r {
        return CYAN
    } else {
        return MAGENTA
    }
}

func Lightness(color color.Color) byte {
    r, g, b, _ := color.RGBA()
    if r != 0 && g != 0 && b != 0 {
        return LIGHT
    }
    if r < 0xffff && g < 0xffff && b < 0xffff {
        return DARK
    }
    return NORMAL
}

func IsWhite(color color.Color) bool {
    r, g, b, _ := color.RGBA()
    return r == 0xffff && g == 0xffff && b == 0xffff 
}

func IsBlack(color color.Color) bool {
    r, g, b, _ := color.RGBA()
    return r == 0 && g == 0 && b == 0 
}
