package main

import "fmt"

//go:generate go build -ldflags='-X=main.VersionString=1.0' main.go

var VersionString = "unset"

func main() {
	fmt.Println("Version:", VersionString)
}
