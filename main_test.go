package main

import (
	"testing"

	sampledata "github.com/corriganrenard/atm-logger/sample-data"
)

func TestFoo(t *testing.T) {
	t.Log("test")
	sampledata.CompareInt(1, 3)

	//sampledata.LineNumToIndex()

}
