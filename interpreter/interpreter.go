package interpreter

import (
	"fmt"
	"image"
	"image/color"
    log "github.com/sirupsen/logrus"
)

// =============== Stack Machine =================

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

type Operand byte

const (
    Push Operand = iota + 1
    Pop 
    Add 
    Sub 
    Mult 
    Div 
    Mod 
    Not 
    Greater 
    Pointer 
    Switch 
    Dup 
    Roll 
    NumIn 
    CharIn 
    NumOut 
    CharOut
)

type StackMachine struct {
    stack Stack
    dp DpDir
    cc CcDir
}

func NewStackMachine(capacity int) StackMachine {
    return StackMachine{
        stack: NewStack(capacity),
        dp: DpRight,
        cc: CcLeft}
}

// TODO JH relocate
func (s *Stack) Roll(depth int32, rolls int32) {
    ip := s.Len() - int(depth)
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

type CodelImage struct {
    image image.Image
    csize int
    bounds image.Rectangle
}

func NewCodelImage(csize int, img image.Image) *CodelImage {
    imageBounds := img.Bounds()

    ibounds := image.Rectangle{
        Min: image.Point{X:0, Y:0},
        Max: image.Point{imageBounds.Max.X / csize, imageBounds.Max.Y / csize}}

    return &CodelImage{csize: csize, image: img, bounds: ibounds}
}

func (c CodelImage) ColorModel() color.Model {
    return c.image.ColorModel()
}

func (c CodelImage) Bounds() image.Rectangle {
    return c.bounds
}

func (c CodelImage) At(x int, y int) color.Color {
    return c.image.At(x * c.csize, y * c.csize)
}

func (sm StackMachine) Peek() (bool, int32) {
    return sm.stack.Peek()
}

func (sm *StackMachine) exec(op Operand, args ...int32) (bool, error) {
    log.Debugf("Executing: %s\n", op)
    switch op {
    case Push:
        if len(args) == 1 {
            sm.stack.Push(args[0])
        }
    case Pop:
        sm.stack.Pop()
    case Add:
        if ok, f, s := sm.stack.PopPop(); ok {
            sm.stack.Push(f + s)
        }
    case Sub:
        if ok, f, s := sm.stack.PopPop(); ok {
            sm.stack.Push(s - f)
        }
    case Mult:
        if ok, f, s := sm.stack.PopPop(); ok {
            sm.stack.Push(f * s)
        }
    case Div:
        if ok, f, s := sm.stack.PopPop(); ok {
            sm.stack.Push(s / f)
        }
    case Mod:
        if ok, f, s := sm.stack.PopPop(); ok {
            sm.stack.Push(s % f)
        }
    case Not:
        if ok, val := sm.stack.Pop(); ok {
            if val == 0 {
                sm.stack.Push(1)
            } else {
                sm.stack.Push(0)
            }
        }
    case Greater:
        if ok, f, s := sm.stack.PopPop(); ok {
            if s > f {
                sm.stack.Push(1)
            } else {
                sm.stack.Push(0)
            }
        }
    case Pointer: // no-op
        if ok, val := sm.stack.Pop(); ok {
//            fmt.Printf("Rotating pointer (%d) %s ", val, sm.dp) 
            sm.RotateDp(val)
//            fmt.Printf("to %s\n", sm.dp)
        }
    case Switch:  // no-op
        if ok, val := sm.stack.Pop(); ok {
            if val % 2 == 1 {
                sm.ToggleCC()
            }
        }
    case Dup:
        if ok, val := sm.stack.Peek(); ok {
            sm.stack.Push(val)
        }
    case Roll:
        if ok, rolls, depth := sm.stack.PopPop(); ok {
            sm.stack.Roll(depth, rolls)
        }
    case NumIn:
    case CharIn:
    case NumOut:
    case CharOut:
        ok, val := sm.stack.Pop()
        if ok {
            fmt.Print(string(val))
        }
    default:
        panic(fmt.Sprintf("Unknown instruction %v", op))
    }
    return false, fmt.Errorf("Unsupported instruction %v", op)
}

func (o Operand) String() string {
    switch o {
    case Push:
        return "push"
    case Pop:
        return "pop"
    case Add:
        return "add"
    case Sub:
        return "sub"
    case Mult:
        return "mult"
    case Div:
        return "div"
    case Mod:
        return "mod"
    case Not:
        return "not"
    case Greater:
        return "greater"
    case Pointer:
        return "pointer"
    case Switch:
        return "switch"
    case Dup:
        return "dup"
    case Roll:
        return "roll"
    case NumIn:
        return "numin"
    case CharIn:
        return "charin"
    case NumOut:
        return "numout"
    case CharOut:
        return "charout"
    }
    return "unknown"
}

// ====================== END STACK MACHINE ========================



// ====================== BEGIN INTERPRETER ========================

type DpDir byte 

const (
    DpRight DpDir = iota
    DpDown    
    DpLeft
    DpUp
)

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

// TODO JH relocate me
func (sm *StackMachine) RotateDp(amount int32) {
    // TODO JH clean this up
    sm.dp = DpDir(abs(byte(int32(sm.dp) + amount) % 4))
}

type CcDir byte

const (
    CcLeft CcDir = iota
    CcRight 
)

func (c CcDir) String() string {
    switch c {
    case CcLeft:
        return "left"
    case CcRight:
        return "right"
    }
    return "unknown"
}


// TODO JH this should be a LUT
func (c CcDir) Direction(dp DpDir) image.Point {
    return dir_cc_lut[dp][c]
}

// TODO JH relocate me
func (sm *StackMachine) ToggleCC() {
    if sm.cc == CcLeft {
        sm.cc = CcRight
    } else {
        sm.cc = CcLeft
    }
}

type PietInterpreter struct {
    sm StackMachine
    pos image.Point
}

func NewInterpreter(capacity int) *PietInterpreter {
    return &PietInterpreter{
        sm: NewStackMachine(capacity)}
}

type Shape struct {
    points []image.Point
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

// HIT
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

// HIT
func find_shape(x int, y int, pi *PietInterpreter, img image.Image, color color.Color, shape *Shape, seen map[int]bool) {
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

// HIT
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

func (pi *PietInterpreter) Execute(image image.Image) {
    // TODO JH should use a constructor instead of this
    pi.init()

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

            x, y := find_next_move(&shape, pi.sm.dp, pi.sm.cc)
//            fmt.Printf("(%d, %d) -> (%d, %d) ", pi.pos.X, pi.pos.Y, x, y)
            if in_bounds(x, y, image) && !is_black(x, y, image) {
                col_next := image.At(x , y)
                operand := pi.diff(col_cur, col_next)
                pi.sm.exec(operand, shape.Size())
                pi.pos.X, pi.pos.Y = x, y

                valid_move = true
            } else {
                attempts--
                pi.sm.ToggleCC()

                x, y = find_next_move(&shape, pi.sm.dp, pi.sm.cc)
//                fmt.Printf("(%d, %d) -> (%d, %d) ", pi.pos.X, pi.pos.Y, x, y)
                if in_bounds(x, y, image) && !is_black(x, y, image) {
                    col_next := image.At(x, y)
                    operand := pi.diff(col_cur, col_next)
                    pi.sm.exec(operand, shape.Size())
                    pi.pos.X, pi.pos.Y = x, y

                    valid_move = true
                } else {
                    attempts--
                    pi.sm.RotateDp(1)
                }
            }
        }
        if attempts == 0 {
            running = false
        }
    }
//    fmt.Println("Done")
}

func (pi *PietInterpreter) init() {
    pi.pos = image.Point{X:0, Y:0}
    pi.sm.dp = DpRight
    pi.sm.cc = CcLeft
}

// HIT
func (pi *PietInterpreter) move(image image.Image) bool {
    // find edges
    edge_x, edge_y := find_edge(pi.pos.X, pi.pos.Y, pi.sm.dp.Direction(), image)
    edge_x, edge_y = find_edge(edge_x, edge_y, pi.sm.cc.Direction(pi.sm.dp), image)

    // attempt to move
    dir := pi.sm.dp.Direction()
    edge_x += dir.X
    edge_y += dir.Y

    if in_bounds(edge_x, edge_y, image) && !is_black(edge_x, edge_y, image) {
        col_next := image.At(edge_x, edge_y)
        col_cur := image.At(pi.pos.X, pi.pos.Y)

        operand := pi.diff(col_cur, col_next)
        // TODO JH need to adjust 1 here
        pi.sm.exec(operand, 1)

        pi.pos.X = edge_x
        pi.pos.Y = edge_y
        return true
    } 
    return false
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

func find_next_edge(pi *PietInterpreter, image image.Image, tries int) (bool, int, int, int, int) {
    x, y := pi.pos.X, pi.pos.Y

    for can_move(x, y, pi.sm.dp, image.Bounds()) {
        cur_col := image.At(x, y)
        direction := pi.sm.dp.Direction()
        next_col := image.At(x, y)
        if next_col != cur_col {
            return true, direction.X, direction.Y, x, y
        }
    }

    if tries > 0 && x == pi.pos.X && y == pi.pos.Y {
        pi.sm.dp = (pi.sm.dp + 1) % 4
        return find_next_edge(pi, image, tries - 1)
    }
    return false, -1, -1, -1, -1
}

func can_move(x int, y int, direction DpDir, bounds image.Rectangle) bool {
    p := direction.Direction().Add(image.Point{X:x,Y:y})
    return p.In(bounds)
}

func (pi *PietInterpreter) diff(cur color.Color, next color.Color) Operand {
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

type Stack struct {
    d []int32
    len int
    head int32
}

func NewStack(capacity int) Stack {
    return Stack{
        d: make([]int32, capacity),
        head: -1}
}

func (stack *Stack) Len() int {
    return stack.len
}

func (stack *Stack) Push(num int32) {
    if stack.len + 1 >= len(stack.d) {
        panic(fmt.Sprintf("Stack overflow %d - %d", stack.len + 1, len(stack.d)))
    }
    stack.head += 1
    stack.d[stack.head] = num
    stack.len += 1
}

func (stack *Stack) Pop() (bool, int32) {
    if stack.Len() == 0 {
        return false, -1
    }
    elem := stack.d[stack.head]
    stack.d[stack.head] = 0
    stack.head -= 1
    stack.len -= 1
    return true, elem 
}


func (stack *Stack) merge( merge_func func(int32, int32) int32) {
    ok, first, second := stack.PopPop()
    if ok {
        result := merge_func(first, second)
        stack.Push(result)
    }
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
