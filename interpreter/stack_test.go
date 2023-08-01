package interpreter

import (
    "testing"
)

func TestStackPushAndPop(t *testing.T) {
    values := []int64 { 1, 2, 3}

    stack := Stack{}
    
    for _, val := range values {
        stack.Push(val)
    }

    for i := len(values) -1; i >= 0; i-- {
        ok, ret := stack.Pop()
        if !ok {
            t.Error("stack.Pop: Unable to pop successfully")
            return
        }
        if ret != values[i] {
            t.Errorf("stack.Pop: Pop order incorrect. Expected %d, got %d\n", ret, values[i]) 
            return
        }
    }
}

func TestEmptyStackPeek(t *testing.T) {
    values := []int64 {1, 2, 3}
    stack := Stack{}

    for _, val := range values {
        stack.Push(val)
    }

    ok, val := stack.Peek()
    if !ok {
        t.Error("stack.Peek: Unable to peek successfully")
        return
    }
    expected := values[len(values) -1]
    if val != expected {
        t.Errorf("stack.Peek: Peek order incorrect. Expected %d, got %d\n", val, expected) 
        return
    }
}
