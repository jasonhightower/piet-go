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
    POther
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
    case POther:
        return "Other"
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
    return POther
}

func (c Pcolor) Diff(other Pcolor) Operand {
    if other == PBlack {
        return Break
    } else if c == PBlack {
        panic("Somehow we moved into a Black shape")
    } else if c == PWhite || other == PWhite {
        return Reset
    } else {
        hue := diffInSteps((byte(c) / 3), (byte(other) / 3), 6)
        lightness := diffInSteps(byte(c) % 3, byte(other) % 3, 3)
        return Operand(hue * 3 + lightness)
    }
}

func diffInSteps(val byte, other byte, max byte) byte {
    if other < val {
        return val + max - other
    }
    return  other - val
}

func (c Pcolor) Hue() (bool, byte) {
    if c == PWhite || c == PBlack || c == POther {
        return false, 0
    }
    return true, byte(c) / 3
}

func (c Pcolor) Lightness() (bool, byte) {
    if c == PWhite || c == PBlack || c == POther {
        return false, 0
    }
    return true, byte(c) % 3
}

func ParseGraph(shape *Shape, shapes [][]*Shape) {
    connected := make(map[*Shape]bool)
    ConnectShape(shapes[0][0], shapes, connected)
    log.Debugf("Connected %d shapes", len(connected))
}

type PInterpreter struct {
    sm *StackMachine
    cc CcDir
    dp DpDir
}

func NewInterpreter(capacity int) *PInterpreter {
    return &PInterpreter{
        sm : NewStackMachine(capacity)}
}

func (pi *PInterpreter) Interpret(shape *Shape) {
    const max_attempts = 4
    amount := max_attempts
    running := true
    for running && amount > 0 {
        amount -= 1
        if running, nextShape := pi.interpretNext(shape); running {
            if nextShape == nil {
                pi.cc = pi.cc.Toggle()
                if running, nextShape = pi.interpretNext(shape); running {
                    if nextShape == nil {
                        pi.dp = pi.dp.Rotate(1)

                        shape.color.Diff(
                    } else {
                        amount = max_attempts
                    }
                }
            } else {
                amount = max_attempts  
            }
        }
    }
}

func (ip *PInterpreter) exec(cur *Shape, next *Shape) {
    
}

func (ip *PInterpreter) interpretNext(shape *Shape) (bool, *Shape) {
    nextShape := shape.GetConnection(ip.dp, ip.cc)
    if nextShape == nil {

    }

    return false, nil
}


func (shape *Shape) GetConnection(dp DpDir, cc CcDir) *Shape {
    idx := byte(dp) * 2 + byte(cc)
    return shape.connections[idx]
}

func ConnectShape(shape *Shape, shapes [][]*Shape, connected map[*Shape]bool) {
    if contains, _ := connected[shape]; !contains {
        connected[shape] = true
        shape.connections[0] = FindShape(shape, shapes, DpRight, CcLeft)
        shape.connections[1] = FindShape(shape, shapes, DpRight, CcRight)
        shape.connections[2] = FindShape(shape, shapes, DpDown, CcLeft)
        shape.connections[3] = FindShape(shape, shapes, DpDown, CcRight)
        shape.connections[4] = FindShape(shape, shapes, DpLeft, CcLeft)
        shape.connections[5] = FindShape(shape, shapes, DpLeft, CcRight)
        shape.connections[6] = FindShape(shape, shapes, DpUp, CcLeft)
        shape.connections[7] = FindShape(shape, shapes, DpUp, CcRight)

        for _, connection := range shape.connections {
            if connection != nil {
                if contains, _ = connected[connection]; !contains {
                    ConnectShape(connection, shapes, connected)
                }
            }
        }
    }
}

func FindShape(shape *Shape, shapes [][]*Shape, dp DpDir, cc CcDir) *Shape {
    edge:= shape.FindEdge(dp, cc)
    next := edge.Add(dp.Direction())
    if next.X >= 0 && next.X < len(shapes) && next.Y >= 0 && next.Y < len(shapes[0]) {
        nextShape := shapes[next.X][next.Y]                        
        if nextShape.color != PBlack {
            return nextShape
        }
    } 
    return nil
}

func parseDir(shape *Shape, dp DpDir, cc CcDir, shapes [][]*Shape) Operand {
    edge:= shape.FindEdge(dp, cc)
    next := edge.Add(dp.Direction())
    if next.X >= 0 && next.X < len(shapes) && next.Y >= 0 && next.Y < len(shapes[0]) {
        nextShape := shapes[next.X][next.Y]                        
        return shape.color.Diff(nextShape.color)
    } else {
        return Break
    }
}

func Tokenize(img image.Image) [][]*Shape {
    shapes := []Shape{}

    // This could be a struct, ShapeImage perhaps?
    shapeRefs := make([][]*Shape, img.Bounds().Max.X)
    for x := 0; x < img.Bounds().Max.X; x++ {
        shapeRefs[x] = make([]*Shape, img.Bounds().Max.Y)
    }

    idx := 0
    for x := 0; x < img.Bounds().Max.X; x++ {
        for y := 0; y < img.Bounds().Max.Y; y++ {
            if shapeRefs[x][y] == nil {
                shape := new(Shape)
                shape.color = asPColor(img.At(x, y))
                fillShape(image.Point{X:x, Y:y}, img, shape, shapeRefs)
                shapes = append(shapes, *shape)
                log.Debugf("New shape at (%d, %d) - %s", x, y, shape.color)
            }
            idx++
        }
    }
    return shapeRefs
}

func fillShape(point image.Point, img image.Image, shape *Shape, shapeRefs [][]*Shape) {
    if !point.In(img.Bounds()) {
        return
    }
    if shapeRefs[point.X][point.Y] != nil {
        return
    }
    pointColor := asPColor(img.At(point.X, point.Y))
    if shape.color == pointColor {
        shape.Append(point)
        shapeRefs[point.X][point.Y] = shape
        fillShape(point.Add(DpUp.Direction()), img, shape, shapeRefs)
        fillShape(point.Add(DpRight.Direction()), img, shape, shapeRefs)
        fillShape(point.Add(DpDown.Direction()), img, shape, shapeRefs)
        fillShape(point.Add(DpLeft.Direction()), img, shape, shapeRefs)
    }     
    return
}


