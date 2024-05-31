package pkg

import (
	"fmt"
	"runtime"
)

func TraceTest() {
	defer Trace()()
	foo()
}

func Trace() func() {
	pc, _, _, ok := runtime.Caller(1)

	if !ok {
		panic("not find call")
	}
	fn := runtime.FuncForPC(pc)
	name := fn.Name()

	fmt.Println("enter", name)
	return func() {
		fmt.Println("exit", name)
	}
}

func foo() {
	defer Trace()()
	bar()
}
func bar() {
	defer Trace()()
}
