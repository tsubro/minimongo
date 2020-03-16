package test

import (
	"fmt"
	"reflect"
	// "unsafe"
	// t "minimongo/test1"
)

type testStruct struct {
	test int `minimongo:type`
}

func Testing() {

	o := testStruct{1}
	// Ref(&o)
	fmt.Println(o)
}

func Ref(o interface{}) {
	s := reflect.ValueOf(o).Elem()

	for i:=0; i<s.NumField(); i++ {
		f := s.Type().Field(i)
		tag := f.Tag

		fmt.Println(f.Name, tag)
	}

	// s.Field(0).SetInt(5)
	// rf := s.Field(0)
	// rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	// rf.SetInt(10)
}
