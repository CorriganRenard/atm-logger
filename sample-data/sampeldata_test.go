package sampledata

import (
	"log"
	"strings"
	"testing"
)

func TestCompareInt(t *testing.T) {
	t.Log("test")
	l := CompareInt(3, 3, 4, 5)

	ls := l.GetSummaryAll()

	//log.Printf("log summary: %#v", ls)
	for _, ld := range ls {
		//log.Printf("rutnime line %d key: %d", v, k)
		var tabs strings.Builder
		for i := 0; i < ld.TabNum; i++ {
			tabs.WriteString("\t")
		}
		log.Printf("%s%v", tabs.String(), ld)
		//log.Println("finished get rule")
	}

}
