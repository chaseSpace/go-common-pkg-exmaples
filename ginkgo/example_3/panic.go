package example_3

import "log"

func DoPanic(willPanic bool) bool {
	if willPanic {
		log.Panicf("log-DoPanic")
	}
	log.Printf("log-NotPanic")
	return false
}
