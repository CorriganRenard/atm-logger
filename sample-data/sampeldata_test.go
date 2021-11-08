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
	level := 0
	writeChildren(&level, ls)

	//log.Printf("log summary: %#v", ls)

}

func writeChildren(level *int, rd RuleData) {
	newLevel := *level + 1
	level = &newLevel

	log.Printf("hints: %v", rd.Summary)
	for _, ld := range rd.Children {
		//log.Printf("rutnime line %d key: %d", v, k)
		var tabs strings.Builder
		for i := 0; i < *level; i++ {
			tabs.WriteString("\t")
		}
		log.Printf("ran: %v %s%v children: %d", ld.Triggered, tabs.String(), ld.Title, len(ld.Children))
		if len(ld.Children) > 0 {
			writeChildren(level, ld)
		}
		//log.Println("finished get rule")
	}

}
