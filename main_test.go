package main

import (
	"testing"

	sampledata "github.com/corriganrenard/atm-logger/sample-data"
)

func TestFoo(t *testing.T) {
	t.Log("test")
	sampledata.CompareInt(3, 3, 4, 5)

	//sampledata.LineNumToIndex()

}
