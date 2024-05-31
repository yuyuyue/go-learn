package pkg

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Sex  int
	Age  int `json:"age"`
}

func ReflectTest() {

	// 获取普通类型
	// a := 1
	// fmt.Print(reflect.ValueOf(a))
	// fmt.Print(reflect.TypeOf(a))

	// 获取结构体类型
	// typeMyStruct := reflect.TypeOf(Person{})
	// fmt.Println(typeMyStruct.Name()) // Person
	// fmt.Println(typeMyStruct.Kind()) // struct

	// 获取指针类型
	// personPointer := reflect.TypeOf(&Person{})
	// fmt.Println(personPointer.Elem().Name()) // Person
	// fmt.Println(personPointer.Elem().Kind()) // struct

	// personPointer := reflect.TypeOf(Person{})
	// for i := 0; i < personPointer.NumField(); i++ {
	// 	fieldType := personPointer.Field(i)
	// 	// fmt.Println(fieldType)                                 // struct
	// 	fmt.Println(personPointer.FieldByName(fieldType.Name)) // struct
	// }

	// 使用反射值对象修改变量的值
	// a := 2020
	// rf := reflect.ValueOf(&a)
	// rf.Elem().SetInt(200)
	// fmt.Println(a) // 200

	// 反射类型调用函数
	r_func := reflect.ValueOf(myfunc)
	//设置函数需要传入的参数也必须是反射类型
	params := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)}
	// 反射调用函数
	res := r_func.Call(params)
	//获取返回值
	fmt.Println(res[0].Int())
}

func myfunc(a, b int) int {
	return a + b
}
