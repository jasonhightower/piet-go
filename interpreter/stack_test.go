package interpreter

import (
	"testing"
)

func TestStackPushAndPop(t *testing.T) {
    values := []int32 { 3, 2, 1}

    stack := NewStack(16)

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

func TestRollStack(t *testing.T) {
    values := [7]int32 {25, 13, 11, -7, 1, 2, 3}
    expected := [7]int32 {25, 13, 1, 2, 3, 11, -7}
    stack := NewStack(16)

    for _, val := range values {
        stack.Push(val)
    }

    stack.Roll(5, 3)

    for i := len(values) - 1; i >= 0; i--{
        if ok, val := stack.Pop(); !ok || val != expected[i] {
            if !ok {
                t.Errorf("stack.Pop failed unexpectedly")
                return
            } 
            t.Errorf("Expected %d, Got %d", expected[i], val)
            return
        }
    }
}

func TestEmptyStackPeek(t *testing.T) {
    values := []int32 {1, 2, 3}
    stack := NewStack(16)

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
