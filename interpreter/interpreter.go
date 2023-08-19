package interpreter

import (
	"fmt"
	"image"
	clr "image/color"
	log "github.com/sirupsen/logrus"
)

type Pcolor byte
const (
    PLightRed Pcolor = iota
    PMidRed 
    PDarkRed 
    PLightYellow 
    PMidYellow 
    PDarkYellow 
    PLightGreen 
    PMidGreen 
    PDarkGreen 
    PLightCyan 
    PMidCyan 
    PDarkCyan 
    PLightBlue 
    PMidBlue 
    PDarkBlue 
    PLightMagenta 
    PMidMagenta 
    PDarkMagenta
    PWhite
    PBlack
)

type Operand byte
const (
    Push Operand = 1
    Pop Operand = 2
    Add Operand = 3
    Sub Operand = 4
    Mult Operand = 5
    Div Operand = 6
    Mod Operand = 7
    Not Operand = 8
    Greater Operand = 9
    Pointer Operand = 10
    Switch Operand = 11
    Dup Operand = 12
    Roll Operand = 13
    NumIn Operand = 14
    CharIn Operand = 15
    NumOut Operand = 16
    CharOut Operand = 17
    Break Operand = 18
    Noop Operand = 19
)

type DpDir byte 
const (
    DpRight DpDir = 0
    DpDown DpDir = 1
    DpLeft DpDir = 2
    DpUp DpDir = 3
)

type CcDir byte
const (
    CcLeft CcDir = iota
    CcRight 
)

func (c Pcolor) String() string {
    switch c {
    case PLightRed:
        return "LightRed"
    case PMidRed:
        return "Red"
    case PDarkRed:
        return "DarkRed"
    case PLightGreen:
        return "LightGreen"
    case PMidGreen:
        return "Green"
    case PDarkGreen:
        return "DarkGreen"
    case PLightBlue:
        return "LightBlue"
    case PMidBlue:
        return "Blue"
    case PDarkBlue:
        return "DarkBlue"
    case PLightMagenta:
        return "LightMagenta"
    case PMidMagenta:
        return "Magenta"
    case PDarkMagenta:
        return "DarkMagenta"
    case PLightCyan:
        return "LightCyan"
    case PMidCyan:
        return "Cyan"
    case PDarkCyan:
        return "DarkCyan"
    case PLightYellow:
        return "LightYellow"
    case PMidYellow:
        return "Yellow"
    case PDarkYellow:
        return "DarkYellow"
    case PWhite:
        return "White"
    case PBlack:
        return "Black"
    }
    panic(fmt.Sprintf("Unknown color %d", byte(c)))
}

const pzero uint8 = 0x00
const pmid uint8 = 0xC0
const phigh uint8 = 0xFF

var colorLut = map[clr.Color]Pcolor {
    clr.RGBA{A:phigh}: PBlack,
    clr.RGBA{A:phigh, R: phigh, G: phigh, B: phigh}: PWhite,

    clr.RGBA{A:phigh, R:phigh, G:pmid, B:pmid}: PLightRed,
    clr.RGBA{A:phigh, R:phigh}:                 PMidRed,
    clr.RGBA{A:phigh, R:pmid}:                  PDarkRed,

    clr.RGBA{A:phigh, G:phigh, B:pmid, R:pmid}: PLightGreen,
    clr.RGBA{A:phigh, G:phigh}:                 PMidGreen,
    clr.RGBA{A:phigh, G:pmid}:                  PDarkGreen,

    clr.RGBA{A:phigh, B:phigh, R:pmid, G:pmid}: PLightBlue,
    clr.RGBA{A:phigh, B:phigh}:                 PMidBlue,
    clr.RGBA{A:phigh, B:pmid}:                  PDarkBlue,

    clr.RGBA{A:phigh, R:phigh, G:pmid, B:phigh}:PLightMagenta,
    clr.RGBA{A:phigh, R:phigh, B:phigh}:        PMidMagenta,
    clr.RGBA{A:phigh, R:pmid, B:pmid}:          PDarkMagenta,

    clr.RGBA{A:phigh, R:phigh, G:phigh, B:pmid}:PLightYellow,
    clr.RGBA{A:phigh, R:phigh, G:phigh}:        PMidYellow,
    clr.RGBA{A:phigh, R:pmid, G:pmid}:          PDarkYellow,

    clr.RGBA{A:phigh, R:pmid, G:phigh, B:phigh}:PLightCyan,
    clr.RGBA{A:phigh, G:phigh, B:phigh}:        PMidCyan,
    clr.RGBA{A:phigh, G:pmid, B:pmid}:          PDarkCyan,
}

func asPColor(color clr.Color) Pcolor {
    if result, ok := colorLut[color]; ok {
        return result
    } else {
        log.Debugf("other - %s", color)
    }
    panic("Unrecognized color")
}

func diffInSteps(val byte, other byte, max byte) byte {
    if other < val {
        return other + max - val
    }
    return  other - val
}

func (c Pcolor) Hue() (bool, byte) {
    if c == PWhite || c == PBlack {
        return false, 0
    }
    return true, byte(c) / 3
}

func (c Pcolor) Lightness() (bool, byte) {
    if c == PWhite || c == PBlack {
        return false, 0
    }
    return true, byte(c) % 3
}

func NewInterpreter(capacity int) *PInterpreter {
    return &PInterpreter{
        stack : NewStack(capacity)}
}

func (pi *PInterpreter) InterpretTokens(shape *Shape) {
    prevShape := shape
    ok := true
    for ok {
        if ok, nextShape := pi.move(prevShape); ok {
            operand := prevShape.Color().Diff(nextShape.Color())
            pi.exec(operand, prevShape.Size())
            prevShape = nextShape
        } else {
            return
        }
    }
}

func (ip *PInterpreter) move(shape *Shape) (bool, *Shape) {
    for attempts := 0; attempts < 4; attempts++ {
        if ok, nextShape := ip.attempt(shape); ok {
            return true, nextShape
        }
        ip.dp = ip.dp.Rotate(1)
    }
    return false, nil
}

func (ip *PInterpreter) attempt(shape *Shape) (bool, *Shape) {
    if ok, nextShape := ip.nextShape(shape); ok {
        return true, nextShape
    } else {
        ip.cc = ip.cc.Toggle()
        return ip.nextShape(shape)
    }
}

func (ip *PInterpreter) nextShape(shape *Shape) (bool, *Shape) {
    var nextShape *Shape
    var ok bool
    nextShape, ok = shape.Connection(ip.dp, ip.cc)
    if !ok || nextShape == nil || nextShape.Color() == PBlack {
        return false, nil
    }
    return true, nextShape
}

func (pi PInterpreter) peek() (bool, int32) {
    return pi.stack.Peek()
}

func (ip PInterpreter) debugOp(op Operand, output string, args ...int32) {
    state := output
    state += fmt.Sprintf("%s ", op)
    for i := 0; i < ip.stack.Len(); i++ {
        state += fmt.Sprintf("%8d", ip.stack.d[i])
    }
    log.Debugf(state)
}

func (pi *PInterpreter) debugOut(op Operand, args ...int32) string {
    state := ""
    if op == CharOut {
        if ok, val := pi.peek(); ok {
            state += fmt.Sprintf(" %s  ", string(val))
        }
    } else {
        state += "    "
    }
    return state
}

func (pi *PInterpreter) exec(op Operand, args ...int32) (bool, error) {
    output := pi.debugOut(op, args...)
    switch op {
    case Push:
        if len(args) == 1 {
            pi.stack.Push(args[0])
        }
    case Pop:
        pi.stack.Pop()
    case Add:
        if ok, f, s := pi.stack.PopPop(); ok {
            pi.stack.Push(f + s)
        }
    case Sub:
        if ok, f, s := pi.stack.PopPop(); ok {
            pi.stack.Push(s - f)
        }
    case Mult:
        if ok, f, s := pi.stack.PopPop(); ok {
            pi.stack.Push(f * s)
        }
    case Div:
        if ok, f, s := pi.stack.PopPop(); ok {
            pi.stack.Push(s / f)
        }
    case Mod:
        if ok, f, s := pi.stack.PopPop(); ok {
            pi.stack.Push(s % f)
        }
    case Not:
        if ok, val := pi.stack.Pop(); ok {
            if val == 0 {
                pi.stack.Push(1)
            } else {
                pi.stack.Push(0)
            }
        }
    case Greater:
        if ok, f, s := pi.stack.PopPop(); ok {
            if s > f {
                pi.stack.Push(1)
            } else {
                pi.stack.Push(0)
            }
        }
    case Pointer:
        if ok, val := pi.stack.Pop(); ok {
            pi.dp = pi.dp.Rotate(val)
        }
    case Switch:
        if ok, val := pi.stack.Pop(); ok {
            if val % 2 == 1 {
                pi.cc = pi.cc.Toggle()
            }
        }
    case Dup:
        if ok, val := pi.stack.Peek(); ok {
            pi.stack.Push(val)
        }
    case Roll:
        if ok, rolls, depth := pi.stack.PopPop(); ok {
            pi.stack.Roll(depth, rolls)
        }
    case NumIn:
    case CharIn:
    case NumOut:
    case CharOut:
        ok, val := pi.stack.Pop()
        if ok {
            fmt.Print(string(val))
        } 
    case Noop:
    default:
        panic(fmt.Sprintf("Unknown instruction %v", op))
    }
    pi.debugOp(op, output, args...)
    return false, fmt.Errorf("Unsupported instruction %v", op)
}

func (o Operand) String() string {
    switch o {
    case Push:
        return "psh"
    case Pop:
        return "pop"
    case Add:
        return "add"
    case Sub:
        return "sub"
    case Mult:
        return "mul"
    case Div:
        return "div"
    case Mod:
        return "mod"
    case Not:
        return "not"
    case Greater:
        return "greater"
    case Pointer:
        return "ptr"
    case Switch:
        return "swi"
    case Dup:
        return "dup"
    case Roll:
        return "rol"
    case NumIn:
        return " in"
    case CharIn:
        return " in"
    case NumOut:
        return "out"
    case CharOut:
        return "out"
    case Noop:
        return "nop"
    }
    return "unknown"
}

func (d DpDir) Direction() image.Point {
    return dir_lut[d]
}

func (d DpDir) String() string {
    switch d {
    case DpRight:
        return "right"
    case DpDown:
        return "down"
    case DpLeft:
        return "left"
    case DpUp:
        return "up"
    }
    return "unknown"
}

func (dp DpDir) Rotate(amount int32) DpDir {
    return DpDir(abs(byte(int32(dp) + amount) % 4))
}

func (c CcDir) Toggle() CcDir {
    if c == CcLeft {
        return CcRight
    }
    return CcLeft
}

func (c CcDir) String() string {
    switch c {
    case CcLeft:
        return "left"
    case CcRight:
        return "right"
    }
    return "unknown"
}

func (c CcDir) Direction(dp DpDir) image.Point {
    return dir_cc_lut[dp][c]
}
