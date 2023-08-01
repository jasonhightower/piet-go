package interpreter

import (
	"container/list"
	"fmt"
	"image"
	"image/color"
)

const RIGHT byte = 0
const DOWN byte = 1
const LEFT byte = 2
const UP byte = 3

type PietInterpreter struct {
    stack Stack
    cc byte
    dp byte
    x int
    y int
    Csize int
}

func (pi *PietInterpreter) Execute(image image.Image) {
    pi.init()

    max_attempts := 8

    attempts := max_attempts
    for attempts > 0 {
        // store current position
        // find the edge in the direction of dp
        // find the exit point based on cc
        // can we move? if not that's an attempt
        // toggle cc
        // can we move? if not that's an attempt
        // toggle db
        // loop

        if pi.attempt_move() {
            pi.
        }

        if pi.x == 6 && pi.y == 6 { fmt.Printf("%d Attempts left - Checking (%d, %d)\n", attempts, pi.dp, pi.cc)}
        attempts--
        if pi.move(image) {
            attempts = max_attempts
        } else if attempts > 0 {
            if pi.cc == LEFT {
                pi.cc = RIGHT
            } else {
                pi.cc = LEFT
            }
            if pi.x == 6 && pi.y == 6 { fmt.Printf("%d Attempts left - Checking (%d, %d)\n", attempts, pi.dp, pi.cc)}
            attempts--
            if pi.move(image) {
                attempts = max_attempts
            } else {
                pi.dp = ((pi.dp + 1) % 4)
            }
        }
    }
}

func (pi *PietInterpreter) init() {
    pi.x = 0
    pi.y = 0
    pi.dp = RIGHT
    pi.cc = LEFT
}

func (pi *PietInterpreter) move(image image.Image) bool {
    debug := pi.x == 0 && pi.y == 0
    // find edges
    edge_x, edge_y := find_edge(pi.x, pi.y, pi.Csize, pi.dp, image)
    if debug {
        fmt.Printf("EDGE dp: %d, %d\n", edge_x, edge_y)
    }
    edge_x, edge_y = find_edge(edge_x, edge_y, pi.Csize, pi.cc, image)
    if debug {
        fmt.Printf("EDGE: %d, %d\n", edge_x, edge_y)
    }

    if pi.x == 6 && pi.y == 6 {
        fmt.Printf("Getting colors for: %d, %d\n", edge_x, edge_y)
    }

    // attempt to move
    switch pi.dp {
    case UP:
        edge_y--
    case DOWN:
        edge_y++
    case LEFT:
        edge_x--
    case RIGHT:
        edge_x++
    }

    /* TODO JH bring this back this needs to be a function call to find the number of codels in the current block
    size := int64((pi.x - edge_x) * (pi.y - edge_y))
    if size < 0 {
        size = -size
    } */

    if in_bounds(edge_x, edge_y, pi.Csize, image) && !is_black(edge_x, edge_y, pi.Csize, image) {
        if pi.x == 6 && pi.y == 6 {
            fmt.Printf("Getting colors for: %d, %d\n", edge_x, edge_y)
        }
        col_next := image.At(edge_x * pi.Csize, edge_y * pi.Csize)
        col_cur := image.At(pi.x * pi.Csize, pi.y * pi.Csize)

        cmd := pi.diff(col_cur, col_next)
        fmt.Printf("(%d, %d)", pi.x, pi.y)
        pi.cmd(cmd, 1)

        fmt.Printf(" | %s ", cmd_name(cmd))
        elem := pi.stack.data.Back()
        for elem != nil {
            fmt.Printf(" %d", elem.Value.(int64))
            elem = elem.Prev()
        }
        fmt.Println()

        pi.x = edge_x
        pi.y = edge_y
        return true
    } 
    return false
}

func matches(sx int, sy int, tx int, ty int, csize int, image image.Image) bool {
    s_col := image.At(sx * csize, sy * csize)
    t_col := image.At(tx * csize, ty * csize)
    return s_col == t_col
}

func is_black(x int, y int, csize int, image image.Image) bool {
    r, g, b, _ := image.At(x * csize, y * csize).RGBA()
    return r == 0 && g == 0 && b == 0
}

func in_bounds(x int, y int, csize int, image image.Image) bool {
    max_x := image.Bounds().Dx()
    max_y := image.Bounds().Dy()
    return x >= 0 && y >=0 && (x * csize) < max_x && (y * csize) < max_y  
}

func find_edge(x int, y int, csize int, direction byte, image image.Image) (int, int) {
    cx, cy := x, y
    mod_x, mod_y := 0, 0
    switch direction {
    case RIGHT:
        mod_x = 1
    case LEFT:
        mod_x = -1
    case UP:
        mod_y = -1
    case DOWN:
        mod_y = 1
    }
    for in_bounds(cx + mod_x, cy + mod_y, csize, image) && matches(x, y, cx + mod_x, cy + mod_y, csize, image) {
        cx += mod_x
        cy += mod_y
    }
    return cx, cy
}

func find_next_edge(pi *PietInterpreter, image image.Image, tries int) (bool, int, int, int, int) {
    x, y := pi.x, pi.y

    for can_move(x, y, pi.dp, pi.Csize, image.Bounds()) {
        cx := x
        cy := y
        cur_col := image.At(x * pi.Csize, y * pi.Csize)
        switch pi.dp {
        case RIGHT:
            x++
        case LEFT:
            x--
        case UP:
            y--
        case DOWN:
            y++
        }

        next_col := image.At(x * pi.Csize, y * pi.Csize)
        if next_col != cur_col {
            return true, cx, cy, x, y
        }
    }

    if tries > 0 && x == pi.x && y == pi.y {
        pi.dp = (pi.dp + 1) % 4
        return find_next_edge(pi, image, tries - 1)
    }

    fmt.Printf("Could not move to %d, %d - %d, %d - %d\n", x, y, image.Bounds().Size().X, image.Bounds().Size().Y, pi.dp)
    return false, -1, -1, -1, -1
}

func (pi *PietInterpreter) cmd(cmd byte, blocksize int64) {
    switch cmd {
        case 1: 
            pi.Push(blocksize)
        case 2: 
            pi.Pop()
        case 3:
            pi.Add()
        case 4:
            pi.Sub()
        case 5:
            pi.Mult()
        case 6:
            pi.Div()
        case 7:
            pi.Mod()
        case 8:
            pi.Not()
        case 9:
            pi.Greater()
        case 10:
            pi.Pointer()
        case 11:
            pi.Switch()
        case 12:
            pi.Dup()
        case 13:
            pi.Roll()
        case 14:
            pi.NumIn()
        case 15:
            pi.CharIn()
        case 16:
            pi.NumOut()
        case 17:
            pi.CharOut()
    }
}

func cmd_name(cmd byte) string {
    switch cmd {
        case 1: 
            return "psh"
        case 2: 
            return "pop"
        case 3:
            return "add"
        case 4:
            return "sub"
        case 5:
            return "mul"
        case 6:
            return "div"
        case 7:
            return "mod"
        case 8:
            return "not"
        case 9:
            return "gtr"
        case 10:
            return "pnt"
            // TODO JH pointer
        case 11:
            return "swi"
            // TODO JH switch
        case 12:
            return "dup"
        case 13:
            return "rol"
            // TODO JH roll
        case 14:
            // TODO JH in - number
            return "->n"
        case 15:
            // TODO JH in - char
            return "->c"
        case 16:
            // TODO JH out - number
            return "n->"
        case 17:
            // TODO JH out - char
            return "c->"
    }
    panic("Unknown command: " + string(cmd))
}

func can_move(x int, y int, direction byte, csize int, bounds image.Rectangle) bool {
    switch direction {
    case RIGHT:
        return (x + 1) * csize < bounds.Size().X
    case LEFT:
        return x - 1 >= 0
    case UP:
        return y -1 >= 0
    case DOWN:
        return (y + 1) * csize < bounds.Size().Y
    }
    // TODO this should be a panic ... 
    return false
}




func (pi *PietInterpreter) diff(cur color.Color, next color.Color) byte {
    h_steps := steps(Hue(cur), Hue(next), 6)
    l_steps := steps(Lightness(cur), Lightness(next), 3)

    return h_steps * 3 + l_steps
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

func (pi *PietInterpreter) Add() {
    pi.stack.merge(_add_val)
}

func (pi *PietInterpreter) Sub() {
    pi.stack.merge(_subtract)
}

func (pi *PietInterpreter) Mult() {
    pi.stack.merge(_multiply)
}

func (pi *PietInterpreter) Div() {
    pi.stack.merge(_divide)
}

func (pi *PietInterpreter) Mod() {
    pi.stack.merge(_modulo)
}

func (pi *PietInterpreter) Greater() {
    pi.stack.merge(_greater)
}

func (pi *PietInterpreter) Dup() {
    ok, val := pi.stack.Peek()
    if ok {
        pi.stack.Push(val)
    }
}

func (pi *PietInterpreter) Pop() {
    pi.stack.Pop()
}

func (pi *PietInterpreter) Peek() (bool, int64) {
    return pi.stack.Peek()
}

func (pi *PietInterpreter) Push(val int64) {
    pi.stack.Push(val)
}

func (pi *PietInterpreter) Not() {
    ok, top := pi.stack.Pop()
    if ok {
        if top == 0 {
            pi.stack.Push(1)
        } else {
            pi.stack.Push( 0)
        }
    }
}

func (pi *PietInterpreter) Pointer() {
    ok, count := pi.stack.Pop()
    if ok {
        pi.dp = byte((pi.dp + byte(count)) % 4)
        /*
        switch pi.dp {
        case DOWN:
            fmt.Println("pointing DOWN")
        case LEFT:
            fmt.Println("pointing LEFT")
        case UP:
            fmt.Println("pointing UP")
        case RIGHT:
            fmt.Println("pointing RIGHT")
        default:
            fmt.Printf("pointing <unknown> %d\n", pi.dp)
        }
        */
    }
}

func (pi *PietInterpreter) Switch() {
    ok, times := pi.stack.Pop()
    if ok && times % 2 != 0 {
        if pi.cc == LEFT {
            pi.cc = RIGHT
        } else {
            pi.cc = LEFT
        }
    }
}

func (pi *PietInterpreter) Roll() {
    ok, rolls, depth := pi.stack.PopPop()
    if ok {
        insert_point := pi.stack.data.Front()
        for i := 1; i < int(depth); i++ {
            insert_point = insert_point.Next()
        }

        for i := 0; i < int(rolls); i++ {
            elem := pi.stack.data.Front()
            pi.stack.data.MoveAfter(elem, insert_point)
            insert_point = elem
        }
    }
}

func (pi *PietInterpreter) CharIn() {
   // toto later 
}

func (pi *PietInterpreter) NumIn() {
    // todo later
}

func (pi *PietInterpreter) CharOut() {
    ok, val := pi.stack.Pop()
    if ok {
        fmt.Print(string(val))
    }
}

func (pi *PietInterpreter) NumOut() {
    ok, val := pi.Peek()
    if ok {
        pi.Pop()
        fmt.Print(val)
    }
}


type Stack struct {
    data list.List
}

func (stack *Stack) Len() int {
    return stack.data.Len()
}

func (stack *Stack) Push(num int64) {
    stack.data.PushFront(num)
}

func (stack *Stack) Pop() (bool, int64) {
    if stack.Len() == 0 {
        return false, -1
    }
    elem := stack.data.Front()
    return true, stack.data.Remove(elem).(int64)
}


func (stack *Stack) merge( merge_func func(int64, int64) int64) {
    ok, first, second := stack.PopPop()
    if ok {
        result := merge_func(first, second)
        stack.Push(result)
    }
}

func (stack *Stack) Peek() (bool, int64) {
    if stack.Len() > 0 {
        return true, stack.data.Front().Value.(int64)
    }
    return false, -1
}

func (stack *Stack) PopPop() (bool, int64, int64) {
    if stack.Len() > 1 {
        _, first := stack.Pop()
        _, second := stack.Pop()
        return true, first, second
    }
    return false, -1, -1
}

func _greater(f int64, s int64) int64 {
    if s > f {
        return 1
    }
    return 0
}

func _modulo(f int64, s int64) int64 {
    return s % f
}

func _multiply(f int64, s int64) int64 {
    return f * s
}

func _divide(f int64, s int64) int64 {
    return s/f
}

func _subtract(f int64, s int64) int64 {
    return s - f
}

func _add_val(f int64, s int64) int64 {
    return f + s
}

