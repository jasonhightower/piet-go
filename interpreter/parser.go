package interpreter

import (
	"fmt"
	"image"

	log "github.com/sirupsen/logrus"
)

type Parser struct {
    stack Stack
    tokens [][]*Shape
    pos image.Point
    bounds image.Rectangle
    dp DpDir
    cc CcDir
}

type Operation struct {
    Op Operand
    Val int32
}

func NewParser(tokens [][]*Shape, bounds image.Rectangle) *Parser {
    // TODO JH specify capacity
    return &Parser{tokens: tokens, bounds: bounds, stack: *NewStack(512)}
}

func (p *Parser) Parse() *[]Operation {
    log.Debug("Parsing")
    operations := []Operation{}        

    for true {
        log.Debugf("(%d, %d)", p.pos.X, p.pos.Y)
        p.move(&operations)
        lastOp := operations[len(operations) - 1]
        log.Debugf("(%d, %d) - %s - %s", p.pos.X, p.pos.Y, lastOp.Op, p.dp)
        switch lastOp.Op {
            case Exit: {
                return &operations
            }
            case Switch: {
                if ok, val := p.stack.Pop(); ok {
                    if val % 2 > 0 {
                        p.cc = p.cc.Toggle()
                    }
                }
            }
            case Push: {
                p.stack.Push(lastOp.Val)
            }
            case Pop: {
                p.stack.Pop()
            }
            case Add: {
                if ok, f, s := p.stack.PopPop(); ok {
                    p.stack.Push(f + s)
                }
            }
            case Sub: {
                if ok, f, s := p.stack.PopPop(); ok {
                    p.stack.Push(s - f)
                }
            }
            case Mult: {
                if ok, f, s := p.stack.PopPop(); ok {
                    p.stack.Push(s * f)
                }
            }
            case Div: {
                if ok, f, s := p.stack.PopPop(); ok {
                    p.stack.Push(s / f)
                }
            }
            case Mod: {
                if ok, f, s := p.stack.PopPop(); ok {
                    p.stack.Push(s % f)
                }
            }
            case Not: {
                if ok, val := p.stack.Pop(); ok {
                    if val == 0 {
                        p.stack.Push(1)
                    } else {
                        p.stack.Push(0)
                    }
                }
            }
            case Greater: {
                if ok, f, s := p.stack.PopPop(); ok {
                    if s > f {
                        p.stack.Push(1)
                    } else {
                        p.stack.Push(0)
                    }
                }
            }
            case Pointer: {
                if ok, val := p.stack.Pop(); ok {
                    p.dp = p.dp.Rotate(val)
                }
            }
            case Dup: {
                if ok, val := p.stack.Peek(); ok {
                    p.stack.Push(val)
                }
            }
            case Roll: {
                if ok, f, s := p.stack.PopPop(); ok {
                    p.stack.Roll(s, f)
                }
            }
            case NumIn: {
                // push "unknown"
            }
            case CharIn: {
                // push "unknown"
            }
            case NumOut: {
                if ok, val := p.stack.Pop(); ok {
                    fmt.Print(val)
                }
            }
            case CharOut: {
                if ok, val := p.stack.Pop(); ok {
                    fmt.Print(string(val))
                }
            }
            case Noop: {
            }
        }
        if lastOp.Op == Exit {
            return &operations
        } else if lastOp.Op == Switch {
            // check to see if this is deterministic ?? the stack is not altered by user input?? perhaps an optimization in a second/third pass
            // rewrite this as an if statement?
            // 4 potential outcomes based on % of current value on stack
            //  need to build instruction sets for each option
            // instruction sets need to be labelled based on x_y of next move
            // this needs to be re-written as a go to
            // check to see if instruction set already exists before building it
        } else if lastOp.Op == Pointer {
            // check to see if this is deterministic ?? the stack is not altered by user input?? perhaps an optimization in a second/third pass
            // check to see if this is deterministic ?? the stack is not altered by user input??
            // rewrite this as multiple if statements?
            // 2 potential outcomes based on % of current value on stack
            // need to build instruction sets for each option
            // instruction sets need to be labelled based on x_y of next move
            // check to see if instruction set already exists before building it
        } else if lastOp.Op == CharIn || lastOp.Op == NumIn {
            // TODO JH need to push an "Unknown" value onto the stack
            // any operations that operate on "unknowns" need to result in additional "Unknowns"        
        }
    }
    return nil
}

func (p *Parser) move(operations *[]Operation) {
    previous := p.tokens[p.pos.X][p.pos.Y]
    var nextPos image.Point
    if previous.color == PWhite {
        for i := 4; i >= 0; i-- {
            p.slideToEdge()        
            nextPos = p.pos.Add(p.dp.Direction())
            if nextPos.In(p.bounds) {
                next := p.tokens[nextPos.X][nextPos.Y]
                if next.color != PBlack {
                    op := previous.color.Diff(next.color)
                    // TODO JH need to check if op is Pointer or Switch and handle each case
                    // can perhaps check and handle that in the function that calls move??
                    p.pos = nextPos
                    *operations = append(*operations, Operation{Op: op, Val: previous.Size()})
                    return
                }
            }
            *operations = append(*operations, Operation{Op: Pointer, Val: 1})
            *operations = append(*operations, Operation{Op: Switch, Val: 1})
            p.dp = p.dp.Rotate(1)
            p.cc = p.cc.Toggle()
        }
    } else {
        for i := 4; i >= 0; i-- {
            p.moveToEdge()
            nextPos := p.pos.Add(p.dp.Direction())
            if nextPos.In(p.bounds) {
                next := p.tokens[nextPos.X][nextPos.Y]
                if next.color != PBlack {
                    op := previous.color.Diff(next.color)
                    // TODO JH need to check if op is Pointer or Switch and handle each case
                    // can perhaps check and handle that in the function that calls move??
                    p.pos = nextPos
                    *operations = append(*operations, Operation{Op: op, Val: previous.Size()})
                    return
                }
            }
            *operations = append(*operations, Operation{Op: Switch, Val: 1})
            p.cc = p.cc.Toggle()
            p.moveToEdge()
            nextPos = p.pos.Add(p.dp.Direction())
            if nextPos.In(p.bounds) {
                next := p.tokens[nextPos.X][nextPos.Y]
                if next.color != PBlack {
                    op := previous.color.Diff(next.color)
                    // TODO JH need to check if op is Pointer or Switch and handle each case
                    // can perhaps check and handle that in the function that calls move??
                    p.pos = nextPos
                    *operations = append(*operations, Operation{Op: op, Val: previous.Size()})
                    return
                }
            }
            *operations = append(*operations, Operation{Op: Pointer, Val: 1})
            p.dp = p.dp.Rotate(1)
        }
    }
    *operations = append(*operations, Operation{Op: Exit, Val: 0})
}

func (p *Parser) slideToEdge() {
    dp := p.dp
    direction := dp.Direction()
    current := p.tokens[p.pos.X][p.pos.Y]
    lookAhead := p.pos.Add(direction)
    for lookAhead.In(p.bounds) {
        token := p.tokens[lookAhead.X][lookAhead.Y]
        if token.color != current.color {
            return
        }         
        p.pos = lookAhead
    }
}

func (p *Parser) moveToEdge() {
    p.pos = p.tokens[p.pos.X][p.pos.Y].FindEdge(p.dp, p.cc)
}

