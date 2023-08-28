package main

import (
	//    bytecode "github.com/jasonhightower/bytecode"
	"fmt"

	"github.com/jasonhightower/jcr"
)

type ConstantPoolBuilder struct {
    cp *jcr.ConstantPool
    utf8s []jcr.CpIndex
    nameTypes []jcr.CpIndex
    classes []jcr.CpIndex
}

func (cpb *ConstantPoolBuilder) AddString(stringIndex jcr.CpIndex) jcr.CpIndex {
    for i, c := range cpb.cp.Constants {
        if c.Type() == jcr.TString {
            if c.(jcr.ConstString).StringIndex  == stringIndex {
                return jcr.CpIndex(i + 1)
            }
        }
    }
    return cpb.cp.Add(jcr.ConstString{StringIndex:stringIndex})
}

func (cpb *ConstantPoolBuilder) AddUtf8(utf8 string) jcr.CpIndex {
    for i, c := range cpb.cp.Constants {
        if c.Type() == jcr.TUtf8 {
            if  c.(jcr.ConstUtf8).String() == utf8 {
                return jcr.CpIndex(i + 1)
            }
        }
    }
    return cpb.cp.Add(jcr.ConstUtf8{Data: []byte(utf8), Length: uint16(len(utf8))})
}
    
func (cpb *ConstantPoolBuilder) AddClass(name jcr.CpIndex) jcr.CpIndex {
    for i, c := range cpb.cp.Constants {
        if c.Type() == jcr.TClass {
            if  c.(jcr.ConstClass).NameIndex == name {
                fmt.Println("Found Classref")
                return jcr.CpIndex(i + 1)
            }
        }
    }
    fmt.Println("Adding new Classref")
    return cpb.cp.Add(jcr.ConstClass{NameIndex: name})
}

func (cpb *ConstantPoolBuilder) AddMethod(class jcr.CpIndex, nameType jcr.CpIndex) jcr.CpIndex {
    for i, c := range cpb.cp.Constants {
        if c.Type() == jcr.TMethodRef {
            m := c.(jcr.ConstMethod)
            if m.ClassIndex == class && m.NameAndTypeIndex == nameType {
                return jcr.CpIndex(i + 1)
            }
        }
    }
    return cpb.cp.Add(jcr.ConstMethod{ClassIndex: class, NameAndTypeIndex: nameType})
}

func (cp *ConstantPoolBuilder) find(cType jcr.ConstantType, matchFunc func(any, any) bool) jcr.CpIndex {

    return jcr.CpIndex(0)
}

func BuildClass(name string) *jcr.Class {

    class := jcr.Class{Major: 64, Minor: 0}
    class.Flags = jcr.FLAG_PUBLIC

    cpBuilder := ConstantPoolBuilder{ cp: &jcr.ConstantPool{}}

    superNameIndex := cpBuilder.AddUtf8("java/lang/Object")
    class.SuperIndex = cpBuilder.AddClass(superNameIndex)

    classNameIndex := cpBuilder.AddUtf8(name)
    class.ThisIndex = cpBuilder.AddClass(classNameIndex)



    class.ConstantPool = cpBuilder.cp

    return &class
}


