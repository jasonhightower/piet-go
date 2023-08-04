package interpreter

import (
    "testing"
)


func TestCommand(t *testing.T) {

    pi := NewInterpreter(32)

    cmd := pi.diff(l_red, n_red)
    if cmd != Operand(1) {
        t.Errorf("Expected %d got %d", 2, cmd)
    }

    cmd = pi.diff(l_red, d_red)
    if cmd != Operand(2) {
        t.Errorf("Expected %d got %d", 2, cmd)
    }

    cmd = pi.diff(l_red, d_red)
    if cmd != Operand(2) {
        t.Errorf("Expected %d got %d", 2, cmd)
    }

    cmd = pi.diff(d_red, l_red)
    if cmd != Operand(1) {
        t.Errorf("Expected %d got %d", 1, cmd)
    }

}


func TestPushAndPop(t *testing.T) {
    values := []int32 {1, 2, 3} 

    pi := NewInterpreter(32)

    pi.sm.exec(Pop)

    for _, val := range values {
        pi.sm.exec(Push, val)
    }

    _, val := pi.sm.Peek()
    if val != values[len(values) - 1] {
        t.Errorf("Expected %d to be returned, got %d", values[len(values) - 1], val)
        return
    }
    pi.sm.exec(Pop)
    _, val = pi.sm.Peek()
    if val != values[len(values) - 2] {
        t.Errorf("Expected %d to be returned, got %d", values[len(values) - 2], val)
    }
}

func TestAdd(t *testing.T) {
    first, second := int32(5), int32(10)
    expected := first + second

    pi := NewInterpreter(32)
    pi.sm.exec(Push, 649201337)
    pi.sm.exec(Push, second)
    pi.sm.exec(Push, first)

    pi.sm.exec(Add)

    ok, result := pi.sm.Peek() 
    if !ok {
        t.Error("Unable to peek after add")
        return
    }
    if result != expected {
        t.Errorf("Addition was not executed correctly. Expected %d but got %d", expected, result)
    }
}

func TestSub(t *testing.T) {
    first, second := int32(7), int32(3)
    expected := second - first

    pi := NewInterpreter(32)
    pi.sm.exec(Push, 649201337)
    pi.sm.exec(Push,second)
    pi.sm.exec(Push,first)

    pi.sm.exec(Sub)

    ok, result := pi.sm.Peek() 
    if !ok {
        t.Error("Unable to peek after sub")
        return
    }
    if result != expected {
        t.Errorf("Subtraction was not executed correctly. Expected %d but got %d", expected, result)
    }
}

func TestMult(t *testing.T) {
    first, second := int32(7), int32(3)
    expected := second * first

    pi := NewInterpreter(32)
    pi.sm.exec(Push,649201337)
    pi.sm.exec(Push,second)
    pi.sm.exec(Push,first)

    pi.sm.exec(Mult)

    ok, result := pi.sm.Peek() 
    if !ok {
        t.Error("Unable to peek after add")
        return
    }
    if result != expected {
        t.Errorf("Multiplication was not executed correctly. Expected %d but got %d", expected, result)
    }

}


func TestDiv(t *testing.T) {
    first, second := int32(7), int32(3)
    expected := second / first

    pi := NewInterpreter(32)
    pi.sm.exec(Push,649201337)
    pi.sm.exec(Push,second)
    pi.sm.exec(Push,first)

    pi.sm.exec(Div)

    ok, result := pi.sm.Peek() 
    if !ok {
        t.Error("Unable to peek after divide")
        return
    }
    if result != expected {
        t.Errorf("Divide was not executed correctly. Expected %d but got %d", expected, result)
    }

}


func TestMod(t *testing.T) {
    first, second := int32(7), int32(3)
    expected := second % first

    pi := NewInterpreter(32)
    pi.sm.exec(Push,649201337)
    pi.sm.exec(Push,second)
    pi.sm.exec(Push,first)

    pi.sm.exec(Mod)

    ok, result := pi.sm.Peek() 
    if !ok {
        t.Error("Unable to peek after add")
        return
    }
    if result != expected {
        t.Errorf("Modulo was not executed correctly. Expected %d but got %d", expected, result)
    }
}

func TestNot(t *testing.T) {
    pi := NewInterpreter(32)
    pi.sm.exec(Push,649201337)
    pi.sm.exec(Push,0)
    pi.sm.exec(Not)

    expected := int32(1)
    _, result := pi.sm.Peek()
    if result != expected {
        t.Errorf("Not was not executed correctly. Expected %d but got %d", expected, result)
    }
    
    pi.sm.exec(Push,1)
    pi.sm.exec(Not)

    expected = int32(0)
    _, result = pi.sm.Peek()
    if result != expected {
        t.Errorf("Not was not executed correctly. Expected %d but got %d", expected, result)
    }
}

func TestGreater_True(t *testing.T) {
    pi := NewInterpreter(32)
    pi.sm.exec(Push,649201337)
    pi.sm.exec(Push,0)
    pi.sm.exec(Greater)

    expected := int32(1)
    _, result := pi.sm.Peek()
    if result != expected {
        t.Errorf("Greater was not executed correctly. Expected %d but got %d", expected, result)
    }
}


func TestGreater_False(t *testing.T) {
    pi := NewInterpreter(32)
    pi.sm.exec(Push,649201337)
    pi.sm.exec(Push,0)
    pi.sm.exec(Greater)

    expected := int32(1)
    _, result := pi.sm.Peek()
    if result != expected {
        t.Errorf("Greater was not executed correctly. Expected %d but got %d", expected, result)
    }

    pi.sm.exec(Push,2)
    pi.sm.exec(Greater)
    expected = int32(0)
    _, result = pi.sm.Peek()
    if result != expected {
        t.Errorf("Greater was not executed correctly. Expected %d but got %d", expected, result)
    }

}

func TestDuplicate(t *testing.T) {
    pi := NewInterpreter(32)
    pi.sm.exec(Push, 649201337)
    pi.sm.exec(Push, 3)
    pi.sm.exec(Dup)

    expected := int32(3)
    _, result := pi.sm.Peek()
    if result != expected {
        t.Errorf("Duplicate was not performed correctly. Expected %d but got %d", expected, result)
    }
    pi.sm.exec(Pop)
    _, result = pi.sm.Peek()
    if result != expected {
        t.Errorf("Value was not duplicated. Expected %d but got %d", expected, result)
    }
}

func TestEmptyStackPop(t *testing.T) {
    stack := Stack{}

    ok, val := stack.Pop()
    if ok {
        t.Error("Pop from empty stack should not have returned ok")
        return
    }
    if val != -1 {
        t.Error("Pop should return -1 when ok == false")
        return
    }
}
