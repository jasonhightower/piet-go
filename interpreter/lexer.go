package interpreter

import (
	"fmt"
	"image"
	clr "image/color"
	log "github.com/sirupsen/logrus"
)

type PInterpreter struct {
    stack *Stack
    cc CcDir
    dp DpDir
    pos image.Point
}

type Stack struct {
    d []int32
    len int
    head int32
}

type CodelImage struct {
    image image.Image
    csize int
    bounds image.Rectangle
}

type Shape struct {
    points []image.Point
    color Pcolor
    connections map[byte]*Shape
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

func ParseTokens(shapes [][]*Shape) {
    connected := make(map[*Shape]bool)
    connectShape(shapes[0][0], shapes, connected)
}

// FIXME should this be part of Pcolor or something else?
func (c Pcolor) Diff(other Pcolor) Operand {
    if c == PWhite || other == PWhite {
        return Noop
    } else if other == PBlack {
        return Break
    } else if c == PBlack {
        panic("Somehow we moved into a Black shape")
    } else {
        hue := diffInSteps((byte(c) / 3), (byte(other) / 3), 6)
        lightness := diffInSteps(byte(c) % 3, byte(other) % 3, 3)
        return Operand(hue * 3 + lightness)
    }
}

func connectShape(shape *Shape, shapes [][]*Shape, connected map[*Shape]bool) {
    if contains, _ := connected[shape]; !contains {
        connected[shape] = true
        establishConnection(shape, shapes, DpRight, CcLeft)
        establishConnection(shape, shapes, DpRight, CcRight)
        establishConnection(shape, shapes, DpDown, CcLeft)
        establishConnection(shape, shapes, DpDown, CcRight)
        establishConnection(shape, shapes, DpLeft, CcLeft)
        establishConnection(shape, shapes, DpLeft, CcRight)
        establishConnection(shape, shapes, DpUp, CcLeft)
        establishConnection(shape, shapes, DpUp, CcRight)

        for _, connection := range shape.connections {
            if contains, _ = connected[connection]; !contains {
                connectShape(connection, shapes, connected)
            }
        }
    }
}

func establishConnection(shape *Shape, shapes[][]*Shape, dp DpDir, cc CcDir) {
    if shape.color == PWhite {
        return
    }
    if target, ok := findShape(shape, shapes, dp, cc); ok {
        idx := byte(dp) * 2 + byte(cc)
        if shape.connections == nil {
            shape.connections = make(map[byte]*Shape)
        }
        shape.connections[idx] = target
    }
}

func findShape(shape *Shape, shapes [][]*Shape, dp DpDir, cc CcDir) (*Shape, bool) {
    edge:= shape.FindEdge(dp, cc)
    next := edge.Add(dp.Direction())
    if next.X >= 0 && next.X < len(shapes) && next.Y >= 0 && next.Y < len(shapes[0]) {
        nextShape := shapes[next.X][next.Y]                        
        if nextShape.color == PWhite {
            // FIXME DO SOMETHING!
        } else if nextShape.color != PBlack {
            return nextShape, true
        }
    } 
    return nil, false
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


var dir_lut = [4]image.Point{
    {X:1,Y:0},  // right
    {X:0,Y:1},  // down 
    {X:-1,Y:0}, // left
    {X:0,Y:-1}} // up
var dir_cc_lut = [4][2]image.Point{
    {{X:0, Y:-1}, {X:0, Y:1}}, // right
    {{X:1, Y:0}, {X:-1, Y:0}}, // down
    {{X:0, Y:1}, {X:0, Y:-1}}, // left
    {{X:-1, Y:0}, {X:1, Y:0}}} // up


func NewCodelImage(csize int, img image.Image) *CodelImage {
    imageBounds := img.Bounds()

    ibounds := image.Rectangle{
        Min: image.Point{X:0, Y:0},
        Max: image.Point{imageBounds.Max.X / csize, imageBounds.Max.Y / csize}}

    return &CodelImage{csize: csize, image: img, bounds: ibounds}
}

func (c CodelImage) ColorModel() clr.Model {
    return c.image.ColorModel()
}

func (c CodelImage) Bounds() image.Rectangle {
    return c.bounds
}

func (c CodelImage) At(x int, y int) clr.Color {
    return c.image.At(x * c.csize, y * c.csize)
}

func (s Shape) Color() Pcolor {
    return s.color
}

func (s Shape) Connection(dp DpDir, cc CcDir) (*Shape, bool) {
    if s.connections == nil {
        return nil, false
    }
    var shape *Shape
    var ok bool
    if shape, ok = s.connections[byte(dp) * 2 + byte(cc)]; ok {
        return shape, true 
    }
    return nil, false
}

func NewShape(point image.Point) Shape {
    return Shape{
        points: []image.Point{point}}
}

func (shape *Shape) Size() int32 {
    return int32(len(shape.points))
}

func (shape *Shape) Contains(point image.Point) bool {
    // this is dumb, make it faster later
    for i := 0; i < len(shape.points); i++ {
        if shape.points[i] == point {
            return true
        }
    }
    return false
}

func (shape *Shape) Append(point image.Point) {
    if shape.Contains(point) {
        return
    }
    shape.points = append(shape.points, point)
}

func min(r int, s int) int {
    if r < s {
        return r
    }
    return s
}

func max(r int, s int) int {
    if r > s {
        return r
    }
    return s
}

func (shape *Shape) FindEdge(direction DpDir, cc CcDir) image.Point {
    // find edge in dp direction
    cur_edge := shape.points[0]
    for i := 1; i < len(shape.points); i++ {
        switch direction {
        case DpUp:
            if cur_edge.Y > shape.points[i].Y {
                cur_edge = shape.points[i]
            }         
        case DpDown:
            if cur_edge.Y < shape.points[i].Y {
                cur_edge = shape.points[i]
            }         
        case DpRight:
            if cur_edge.X < shape.points[i].X {
                cur_edge = shape.points[i]
            }         
        case DpLeft:
            if cur_edge.X > shape.points[i].X {
                cur_edge = shape.points[i]
            }         
        }
    }

    switch direction {
    case DpUp:
        for i := 0; i < len(shape.points); i++ {
            if shape.points[i].Y == cur_edge.Y{
               if cc == CcLeft {
                    if cur_edge.X > shape.points[i].X {
                        cur_edge = shape.points[i]
                    }
               } else {
                   if cur_edge.X < shape.points[i].X {
                       cur_edge = shape.points[i]
                   }
               }
            }
        }
    case DpDown:
        for i := 0; i < len(shape.points); i++ {
            if shape.points[i].Y == cur_edge.Y {
               if cc == CcRight {
                    if cur_edge.X > shape.points[i].X {
                        cur_edge = shape.points[i]
                    }
               } else {
                   if cur_edge.X < shape.points[i].X {
                       cur_edge = shape.points[i]
                   }
               }
            }
        }
 
    case DpLeft:
        for i := 0; i < len(shape.points); i++ {
            if shape.points[i].X == cur_edge.X {
               if cc == CcRight {
                    if cur_edge.Y > shape.points[i].Y {
                        cur_edge = shape.points[i]
                    }
               } else {
                   if cur_edge.Y < shape.points[i].Y {
                       cur_edge = shape.points[i]
                   }
               }
            }
        }
 
    case DpRight:
        for i := 0; i < len(shape.points); i++ {
            if shape.points[i].X == cur_edge.X {
               if cc == CcLeft {
                    if cur_edge.Y > shape.points[i].Y {
                        cur_edge = shape.points[i]
                    }
               } else {
                   if cur_edge.Y < shape.points[i].Y {
                       cur_edge = shape.points[i]
                   }
               }
            }
        }
    }
    return cur_edge
}

func find_shape(x int, y int, pi *PInterpreter, img image.Image, color clr.Color, shape *Shape, seen map[int]bool) {
    pos := y * img.Bounds().Max.X + x
    if seen[pos] {
        return
    }
    if !in_bounds(x, y, img) {
        seen[pos] = true     
        return
    }
    cur_color := img.At(x, y)
    if color == cur_color {
        shape.Append(image.Point{X: x, Y: y})
        seen[pos] = true
        find_shape(x - 1, y, pi, img, color, shape, seen)
        find_shape(x + 1, y, pi, img, color, shape, seen)
        find_shape(x, y - 1, pi, img, color, shape, seen)
        find_shape(x, y + 1, pi, img, color, shape, seen)
    }
    return
}

func find_next_move(shape *Shape, direction DpDir, cc CcDir) (int, int) {
    var edge image.Point
    if shape.Size() == 1 {
        edge = shape.points[0]
    } else {
        edge = shape.FindEdge(direction, cc)
    }

    switch direction {
    case DpUp:
        return edge.X, edge.Y - 1
    case DpDown:
        return edge.X, edge.Y + 1
    case DpLeft:
        return edge.X - 1, edge.Y
    case DpRight:
        return edge.X + 1, edge.Y
    }
    return 0, 0
}

func (pi *PInterpreter) Interpret(image image.Image) {
    max_attempts := 8

    running := true
    for running {

        // store current position/shape
        shape := Shape{} 
        seen := map[int]bool{}
        col_cur := image.At(pi.pos.X, pi.pos.Y)
        find_shape(pi.pos.X, pi.pos.Y, pi, image, col_cur, &shape, seen)

        attempts := max_attempts
        valid_move := false
        for !valid_move && attempts > 0 {

            x, y := find_next_move(&shape, pi.dp, pi.cc)
            if in_bounds(x, y, image) && !is_black(x, y, image) {
                col_next := image.At(x , y)
                operand := diff(col_cur, col_next)
                pi.exec(operand, shape.Size())
                pi.pos.X, pi.pos.Y = x, y

                valid_move = true
            } else {
                attempts--
                pi.cc = pi.cc.Toggle()

                x, y = find_next_move(&shape, pi.dp, pi.cc)
                if in_bounds(x, y, image) && !is_black(x, y, image) {
                    col_next := image.At(x, y)
                    operand := diff(col_cur, col_next)
                    pi.exec(operand, shape.Size())
                    pi.pos.X, pi.pos.Y = x, y

                    valid_move = true
                } else {
                    attempts--
                    pi.dp = pi.dp.Rotate(1)
                }
            }
        }
        if attempts == 0 {
            running = false
        }
    }
}

func matches(sx int, sy int, tx int, ty int, image image.Image) bool {
    s_col := image.At(sx, sy)
    t_col := image.At(tx, ty)
    return s_col == t_col
}

func is_black(x int, y int, image image.Image) bool {
    r, g, b, _ := image.At(x, y).RGBA()
    return r == 0 && g == 0 && b == 0
}

func in_bounds(x int, y int, image image.Image) bool {
    max_x := image.Bounds().Dx()
    max_y := image.Bounds().Dy()
    return x >= 0 && y >=0 && (x) < max_x && (y) < max_y  
}

func find_edge(x int, y int, direction image.Point, image image.Image) (int, int) {
    cx, cy := x, y
    for in_bounds(cx + direction.X, cy + direction.Y, image) && matches(x, y, cx + direction.X, cy + direction.Y, image) {
        cx += direction.X
        cy += direction.Y
    }
    return cx, cy
}

func find_next_edge(pi *PInterpreter, image image.Image, tries int) (bool, int, int, int, int) {
    x, y := pi.pos.X, pi.pos.Y

    for can_move(x, y, pi.dp, image.Bounds()) {
        cur_col := image.At(x, y)
        direction := pi.dp.Direction()
        next_col := image.At(x, y)
        if next_col != cur_col {
            return true, direction.X, direction.Y, x, y
        }
    }

    if tries > 0 && x == pi.pos.X && y == pi.pos.Y {
        pi.dp = (pi.dp + 1) % 4
        return find_next_edge(pi, image, tries - 1)
    }
    return false, -1, -1, -1, -1
}

func can_move(x int, y int, direction DpDir, bounds image.Rectangle) bool {
    p := direction.Direction().Add(image.Point{X:x,Y:y})
    return p.In(bounds)
}

func diff(cur clr.Color, next clr.Color) Operand {
    h_steps := steps(Hue(cur), Hue(next), 6)
    l_steps := steps(Lightness(cur), Lightness(next), 3)

    return Operand(h_steps * 3 + l_steps)
}

func steps(cur byte, next byte, max byte) byte {
    if cur > next {
        next += max
    }
    return next - cur
}

func abs(val byte) byte {
    if val < 0 {
        return -val
    }
    return val
}

func NewStack(capacity int) *Stack {
    return &Stack{
        d: make([]int32, capacity),
        head: -1}
}

func (s *Stack) Roll(depth int32, rolls int32) {
    if s.Len() <= 1 || int(depth) > s.Len() {
        return
    }
    ip := s.Len() - (int(depth) % s.Len())
    i := ip
    j := ip + (int(rolls) % int(depth))
    tmpi := s.d[i]
    var tmpj int32

    for n :=0; n < int(depth);  {
        tmpj = s.d[j]
        s.d[j] = tmpi
        i++
        n++

        if n == int(depth) {
            return
        }
        tmpi = s.d[i]
        s.d[i] = tmpj
        j++
        if j > int(s.head) {
            j = ip
        }
        n++
    }
}

func (stack *Stack) Len() int {
    return int(stack.head) + 1
}

func (stack *Stack) Push(num int32) {
    if int(stack.head) + 1 == len(stack.d) {
        panic(fmt.Sprintf("Stack overflow %d - %d", stack.len + 1, len(stack.d)))
    }
    stack.head += 1
    stack.d[stack.head] = num
}

func (stack *Stack) Pop() (bool, int32) {
    if stack.Len() == 0 {
        return false, -1
    }
    elem := stack.d[stack.head]
    stack.d[stack.head] = 0
    stack.head -= 1
    return true, elem 
}

func (stack *Stack) Peek() (bool, int32) {
    if stack.Len() > 0 {
        return true, stack.d[stack.head]
    }
    return false, -1
}

func (stack *Stack) PopPop() (bool, int32, int32) {
    if stack.Len() > 1 {
        _, first := stack.Pop()
        _, second := stack.Pop()
        return true, first, second
    }
    return false, -1, -1
}


