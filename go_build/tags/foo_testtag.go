//+build !test_tag

package main

import "fmt"

const ENV = "生产"

var some_img_url = "b.com"

func foo() {
	fmt.Println("This is foo without test_tag")
}
