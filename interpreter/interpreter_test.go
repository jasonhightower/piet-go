package interpreter

import (
    "testing"
)


func TestCommand(t *testing.T) {

    pi := PietInterpreter{}

    cmd := pi.diff(l_red, n_red)
    if cmd != 1 {
        t.Errorf("Expected %d got %d", 2, cmd)
    }

    cmd = pi.diff(l_red, d_red)
    if cmd != 2 {
        t.Errorf("Expected %d got %d", 2, cmd)
    }

    cmd = pi.diff(l_red, d_red)
    if cmd != 2 {
        t.Errorf("Expected %d got %d", 2, cmd)
    }

    cmd = pi.diff(d_red, l_red)
    if cmd != 1 {
        t.Errorf("Expected %d got %d", 1, cmd)
    }

}


func TestPushAndPop(t *testing.T) {
    values := []int64 {1, 2, 3} 

    pi := PietInterpreter{}

    pi.Pop()

    for _, val := range values {
        pi.Push(val)
    }

    _, val := pi.Peek()
    if val != values[len(values) - 1] {
        t.Errorf("Expected %d to be returned, got %d", values[len(values) - 1], val)
        return
    }
    pi.Pop()
    _, val = pi.Peek()
    if val != values[len(values) - 2] {
        t.Errorf("Expected %d to be returned, got %d", values[len(values) - 2], val)
    }
}

func TestAdd(t *testing.T) {
    first, second := int64(5), int64(10)
    expected := first + second

    pi := PietInterpreter{}
    pi.Push(649201337)
    pi.Push(second)
    pi.Push(first)

    pi.Add()

    ok, result := pi.Peek() 
    if !ok {
        t.Error("Unable to peek after add")
        return
    }
    if result != expected {
        t.Errorf("Addition was not executed correctly. Expected %d but got %d", expected, result)
    }
}

func TestSub(t *testing.T) {
    first, second := int64(7), int64(3)
    expected := second - first

    pi := PietInterpreter{}
    pi.Push(649201337)
    pi.Push(second)
    pi.Push(first)

    pi.Sub()

    ok, result := pi.Peek() 
    if !ok {
        t.Error("Unable to peek after sub")
        return
    }
    if result != expected {
        t.Errorf("Subtraction was not executed correctly. Expected %d but got %d", expected, result)
    }
}

func TestMult(t *testing.T) {
    first, second := int64(7), int64(3)
    expected := second * first

    pi := PietInterpreter{}
    pi.Push(649201337)
    pi.Push(second)
    pi.Push(first)

    pi.Mult()

    ok, result := pi.Peek() 
    if !ok {
        t.Error("Unable to peek after add")
        return
    }
    if result != expected {
        t.Errorf("Multiplication was not executed correctly. Expected %d but got %d", expected, result)
    }

}


func TestDiv(t *testing.T) {
    first, second := int64(7), int64(3)
    expected := second / first

    pi := PietInterpreter{}
    pi.Push(649201337)
    pi.Push(second)
    pi.Push(first)

    pi.Div()

    ok, result := pi.Peek() 
    if !ok {
        t.Error("Unable to peek after divide")
        return
    }
    if result != expected {
        t.Errorf("Divide was not executed correctly. Expected %d but got %d", expected, result)
    }

}


func TestMod(t *testing.T) {
    first, second := int64(7), int64(3)
    expected := second % first

    pi := PietInterpreter{}
    pi.Push(649201337)
    pi.Push(second)
    pi.Push(first)

    pi.Mod()

    ok, result := pi.Peek() 
    if !ok {
        t.Error("Unable to peek after add")
        return
    }
    if result != expected {
        t.Errorf("Modulo was not executed correctly. Expected %d but got %d", expected, result)
    }
}

func TestNot(t *testing.T) {
    pi := PietInterpreter{}
    pi.Push(649201337)
    pi.Push(0)
    pi.Not()

    expected := int64(1)
    _, result := pi.Peek()
    if result != expected {
        t.Errorf("Not was not executed correctly. Expected %d but got %d", expected, result)
    }
    
    pi.Push(1)
    pi.Not()

    expected = int64(0)
    _, result = pi.Peek()
    if result != expected {
        t.Errorf("Not was not executed correctly. Expected %d but got %d", expected, result)
    }
}

func TestGreater_True(t *testing.T) {
    pi := PietInterpreter{}
    pi.Push(649201337)
    pi.Push(0)
    pi.Greater()

    expected := int64(1)
    _, result := pi.Peek()
    if result != expected {
        t.Errorf("Greater was not executed correctly. Expected %d but got %d", expected, result)
    }
}


func TestGreater_False(t *testing.T) {
    pi := PietInterpreter{}
    pi.Push(649201337)
    pi.Push(0)
    pi.Greater()

    expected := int64(1)
    _, result := pi.Peek()
    if result != expected {
        t.Errorf("Greater was not executed correctly. Expected %d but got %d", expected, result)
    }

    pi.Push(2)
    pi.Greater()
    expected = int64(0)
    _, result = pi.Peek()
    if result != expected {
        t.Errorf("Greater was not executed correctly. Expected %d but got %d", expected, result)
    }

}

func TestDuplicate(t *testing.T) {
    pi := PietInterpreter{}
    pi.Push(649201337)
    pi.Push(3)
    pi.Dup()

    expected := int64(3)
    _, result := pi.Peek()
    if result != expected {
        t.Errorf("Duplicate was not performed correctly. Expected %d but got %d", expected, result)
    }
    pi.Pop()
    _, result = pi.Peek()
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
