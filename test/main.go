package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	s := &service{}
	start := time.Now() // 获取当前时间
	svrType := reflect.TypeOf(s)
	svrValue := reflect.ValueOf(s)
	fmt.Println(svrType)
	fmt.Println(svrValue)
	type req struct {
		a, b int
	}
	type Handler func(interface{}) interface{}
	m := make(map[string]Handler)
	for i := 0; i < svrType.NumMethod(); i++ {
		method := svrType.Method(i)
		m[method.Name] = func(in interface{}) interface{} {
			req := in.(req)
			res := method.Func.Call([]reflect.Value{svrValue, reflect.ValueOf(req.a), reflect.ValueOf(req.b)})
			return res[0]
		}
	}
	m0 := reflect.TypeOf(s).Method(0).Type
	fmt.Println(reflect.TypeOf(s).Method(0).Name)
	for i := 0; i < m0.NumIn(); i++ {
		fmt.Println("参数:", m0.In(i))
	}
	for i := 0; i < m0.NumOut(); i++ {
		fmt.Println("返回值:", m0.Out(i))
	}
	typeOfHero := reflect.TypeOf(*s)
	for i := 0; i < typeOfHero.NumField(); i++ {
		fmt.Printf("field' name is %s, type is %s, kind is %s\n", typeOfHero.Field(i).Name, typeOfHero.Field(i).Type, typeOfHero.Field(i).Type.Kind())
	}
	// 获取名称为 Name 的成员字段类型对象
	// nameField, _ := typeOfHero.FieldByName("name")
	// fmt.Printf("field' name is %s, type is %s, kind is %s\n", nameField.Name, nameField.Type, nameField.Type.Kind())
	f := m["Mut"](req{1, 2})
	fmt.Println("mut", f)
	elapsed := time.Since(start)
	fmt.Println("该函数执行完成耗时：", elapsed)
}

type service struct {
	name string
	age  int
}

func (s *service) Add(a, b int) int {
	return a + b
}
func (s *service) Mut(a, b int) int {
	return a * b
}
