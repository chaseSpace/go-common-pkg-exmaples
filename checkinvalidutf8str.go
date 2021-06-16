package main

import (
	"fmt"
	"reflect"
	"unicode/utf8"
)

type A struct {
	Name string
	C1   *C
	C2   *C
	C3   []*C
}

type C struct {
	Desc string
}

func main() {
	x := &A{
		Name: "xxx",
		C1:   &C{Desc: "ccc1"},
		C2:   &C{Desc: "ccc2"},
		C3: []*C{
			&C{Desc: "luck"},
			&C{Desc: "devil"},
			&C{Desc: string([]byte{0xFF})},
		},
	}
	CheckObjInvalidUTF8Str("", x)
}

// 使用反射找到一个对象中的包含非utf8字符的str字段以及内容
// NOTE: 没有解析map类型
func CheckObjInvalidUTF8Str(path string, obj interface{}) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		s := v.Interface().(string)
		valid := utf8.ValidString(s)
		fmt.Printf("CheckObjInvalidUTF8Str--str:%s valid:%v val:%s\n", path+v.Type().String(), valid, v.String())
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			v1 := v.Index(i)
			CheckObjInvalidUTF8Str(fmt.Sprintf(" %s-[slice:%s-field%d] ", path, v.Type(), i), v1.Interface())
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			v1 := v.Field(i)
			if !v1.CanSet() {
				continue
			}
			CheckObjInvalidUTF8Str(fmt.Sprintf("%s-[struct:%s-field%d] ", path, v.Type(), i), v1.Interface())
		}
	default:
		fmt.Printf("unknown obj kind:%s \n", v.Kind())
	}
}
