package pkg

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// 泛型函数声明：T为类型形参
func maxGenerics[T ordered](sl []T) T {
	if len(sl) == 0 {
		panic("slice is empty")
	}
	max := sl[0]
	for _, v := range sl[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func fooo[T any](a T) T {
	var zero T
	return zero
}

type maxableSlice[T ordered] struct {
	elems []T
}

func doSomething[T1, T2 any](t1 T1, t2 T2) T1 {
	var a T1 // 声明变量
	var b T2
	a, b = t1, t2 // 同类型赋值
	_ = b
	f := func(t T1) {
	}
	f(a)    // 传给其他函数
	p := &a // 取变量地址
	_ = p
	var i interface{} = a // 转换或赋值给interface{}类型变量
	_ = i
	c := new(T1) // 传递给预定义函数
	_ = c
	f(a)                    // 将变量传给其他函数
	sl := make([]T1, 0, 10) // 作为复合类型中的元素类型
	_ = sl
	j, ok := i.(T1) // 用在类型断言中
	_ = ok
	_ = j
	switch i.(type) { // 作为type switch中的case类型
	case T1:
	case T2:
	}
	return a // 从函数返回
}

type Stringer interface {
	ordered
	comparable
	String() string
}

func StringifyWithoutZero[T Stringer](s []T, max T) (ret []string) {
	var zero T
	for _, v := range s {
		if v == zero {
			continue
		}
		if v == zero || v >= max {
			continue
		}
		ret = append(ret, v.String())
	}
	return ret
}

type MyString string

func (s MyString) String() string {
	return string(s)
}

type BasicInterface interface { // 基本接口类型
	M1()
}
type NonBasicInterface interface { // 非基本接口类型
	BasicInterface
	~int | ~string // 包含类型元素
}

// func (MyString) M1() {
// }
func DoubleDefined[S ~[]E, E constraints.Integer](s S) S {}

func foooo[T NonBasicInterface](a T) { // 非基本接口类型作为约束
}

func barr[T BasicInterface](a T) { // 基本接口类型作为约束
}

func GenericTest() {

	// 调用泛型函数：int为类型实参
	// m := maxGenerics[int]([]int{1, 2, -4, -6, 7, 0})
	// f := fooo(4)

	// var sl = maxableSlice[int]{
	// 	elems: []int{1, 2, -4, -6, 7, 0}, // 编译器错误：cannot use generic type maxableSlice[T ordered] without instantiation
	// }
	sl := StringifyWithoutZero([]MyString{"I", "", "love", "", "golang"}, "1")
	var s = MyString("hello")
	var bi BasicInterface = s     // 基本接口类型支持常规用法
	var nbi NonBasicInterface = s // 非基本接口不支持常规用法，导致编译器错误：cannot use type NonBasicInterface outside a type constraint: interface contains type constraints
	bi.M1()
	nbi.M1()
	foooo(s)
	barr(s)
	fmt.Println(sl)
}
