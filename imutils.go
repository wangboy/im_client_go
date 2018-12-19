package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func callReflect(any interface{}, name string, args ... interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	if v := reflect.ValueOf(any).MethodByName(name); v.String() == "<invalid Value>" {
		return nil
	} else {
		return v.Call(inputs)
	}
}

type Test struct{}

func (t *Test) PrintInfo(i int, s string) string {
	fmt.Println("call method PrintInfo i", i, ",s :", s)
	return s + strconv.Itoa(i)
}
func (t *Test) ShowMsg() string {
	fmt.Println("\nshow msg input 'call reflect'")
	return "ShowMsg"
}
func callReflectMethod() {
	fmt.Printf("\n callReflectMethod PrintInfo :%s", callReflect(&Test{}, "PrintInfo", 10, "TestMethod")[0].String())
	fmt.Printf("\n callReflectMethod ShowMsg  %s", callReflect(&Test{}, "ShowMsg")[0].String())

	//<invalid Value>
	callReflect(&Test{}, "ShowMs")
	if result := callReflect(&Test{}, "ShowMs"); result != nil {
		fmt.Printf("\n callReflectMethod ShowMs %s", result[0].String())
	} else {
		fmt.Println("\n callReflectMethod ShowMs didn't run ")
	}
	fmt.Println("\n reflect all ")
}

var typeRegistry = make(map[string]reflect.Type)

func registerType(elem interface{}) {
	t := reflect.TypeOf(elem).Elem()
	typeRegistry[t.Name()] = t
}

func newStruct(name string) (interface{}, bool) {
	elem, ok := typeRegistry[name]
	if !ok {
		return nil, false
	}
	return reflect.New(elem).Elem().Interface(), true
}

func init() {
	registerType((*test)(nil))
}

type TestInterface interface {
	doTest(i int32) string
}

type test struct {
	Name string
	Sex  int
}

func (t *test) testFunc() {
	fmt.Println(" call testFunc ")
}

func (t *test) doTest(i int32) string {
	fmt.Println(" input is ", i, reflect.TypeOf(t))
	return "input is " + string(i)
}

func main2() {
	//callReflectMethod()
	//
	//to := reflect.TypeOf(Test{})
	//fmt.Println(to)
	//
	//newZero := reflect.Zero(to)
	//fmt.Println(reflect.TypeOf(newZero), newZero)

	fmt.Println("=====================")

	structName := "test"

	s, ok := newStruct(structName)
	if !ok {
		return
	}

	fmt.Println(" s type ", reflect.TypeOf(s))

	fmt.Println(s, reflect.TypeOf(s))
	r, e := s.(TestInterface)
	fmt.Println(" cast s ", r, e)

	t, ok := s.(test)
	if !ok {
		return
	}
	t.Name = "i am test"
	fmt.Println(t, reflect.TypeOf(t))

	t.doTest(7777777)

	testIn := t
	testIn.doTest(888888)


	testInterface,ok := s.(TestInterface)
	if !ok {
		fmt.Println(" cast not ok", testInterface , ok)
		return
	}
	testInterface.doTest(123)

}
