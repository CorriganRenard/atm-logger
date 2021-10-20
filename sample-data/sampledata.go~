package sampledata

import (
	"log"
	"runtime"
)

//go:generate atm-logger
func CompareInt(a, b int) {
	if a < b {
		// RULE: a lessn than  b, do something
		saveAction(a, b)
	} else if a > b {
		// RULE: a greater than  b, do something else
		saveAction(a, b)
	} else {
		// RULE: a (%s) == b (%s), do something
		saveAction(a, b)
	}

}

func saveAction(args ...interface{}) uint64 {
	_, _, line, _ := runtime.Caller(1)
	log.Printf("line number: %v args: %#v", line, args)

	return 0
}
