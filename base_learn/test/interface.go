package pkg

type MyInterface interface {
	M1()
}
type T int

func (T) M1() {
	println("T's M1")
}

// 如果类型实现了接口的所有方法，该类型变量可以右值赋值给接口类型 I 的变量
// 如果一个类型 T 的方法集合是某接口类型 I 的方法集合的等价集合或超集，我们就说类型 T 实现了接口类型 I，那么类型 T 的变量就可以作为合法的右值赋值给接口类型 I 的变量
// 类型只要实现接口的方法，就实现了接口，这是go约定的契约。并且开发中遵守小契约
func InterfaceTest() {
	// var t T
	// var i interface{} = t
	// v1, ok := i.(MyInterface)
	// if !ok {
	// 	panic("the value of i is not MyInterface")
	// }
	// v1.M1()
	// fmt.Printf("the type of v1 is %T\n", v1) // the type of v1 is main.T

	// i = int64(13)
	// _, ok2 := i.(MyInterface)
	// if !ok2 {
	// 	panic("the value of i is not MyInterface")
	// }
	var eif1 interface{} = (*int)(nil) // 空接口类型
	// eif1 = (*T)(nil)
	// 空接口 = nil
	println("eif1:", eif1, (*T)(nil)) // (0x0,0x0)
}
