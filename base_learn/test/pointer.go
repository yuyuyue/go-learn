package pkg

import (
	"fmt"
	"reflect"
)

type T1 int

type t2 struct {
	n int
	m int
}

type I interface {
	M1()
}

type S1 struct {
	T1
	*t2
	I
	a int
	b string
}
type S2 struct {
	T1 T1
	t2 *t2
	I  I
	a  int
	b  string
}

func dumpMethodSet(i interface{}) {
	dynTyp := reflect.TypeOf(i)
	if dynTyp == nil {
		fmt.Printf("there is no dynamic type\n")
		return
	}
	n := dynTyp.NumMethod()
	if n == 0 {
		fmt.Printf("%s's method set is empty!\n", dynTyp)
		return
	}
	fmt.Printf("%s's method set:\n", dynTyp)
	for j := 0; j < n; j++ {
		fmt.Println("-", dynTyp.Method(j).Name)
	}
	fmt.Printf("\n")
}

func PointerTest() {
	var s1 S1
	dumpMethodSet(s1)
	// var s2 S2
	var i I
	dumpMethodSet(i)
}
