package main

import (
	"fmt"
	"reflect"
)

type Point struct {
	X    int
	Y    int
	Name string
}

func main() {
	p := Point{X: 3, Y: 5, Name: "sx"}
	rv := reflect.ValueOf(p)
	fmt.Printf("rv.Type = %v\n", rv.Type())
	fmt.Printf("rv.Kind = %v\n", rv.Kind())
	fmt.Printf("rv.Interface = %v\n", rv.Interface())
	xv := rv.Field(0)
	fmt.Printf("rv.Filed = %v\n", xv)
	tp := rv.Type()
	name := tp.Field(0).Name
	fmt.Println(name)
	fmt.Printf("v.FieldByname(name)  = %v\n", rv.FieldByName(name).Interface())

	for i := 0; i < rv.NumField(); i++ {
		fmt.Printf("%#v\n", tp.Field(i))
		fv := rv.Field(i)
		fmt.Printf("%#v\n", fv)
		fmt.Printf("%#v\n", fv.Interface())

	}
}
