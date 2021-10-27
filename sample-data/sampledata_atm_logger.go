package sampledata

import (
	"fmt"
	"log"
	"runtime"
	"sort"
	"strconv"
)

const _atm_logger_name = "a (%v)lessn than  b(%v), do somethinga (%v) greater than  b(%v), do something elsea (%v) == b (%v), do somethingnested func  a (%v)lessn than  b(%v), do somethingnested func a (%v) greater than  b(%v), do something elsenested a (%v) == b (%v), do somethingc (%v) < d (%v), do something elsec (%v) > d (%v), do something else once morec (%v) == d (%v), do something else againnested func  a (%v)lessn than  b(%v), do somethingnested func a (%v) greater than  b(%v), do something elsenested a (%v) == b (%v), do somethingnested func  a (%v)lessn than  b(%v), do somethingnested func a (%v) greater than  b(%v), do something elsenested a (%v) == b (%v), do something"

var _atm_logger_index = [...]uint16{0, 37, 82, 112, 162, 219, 256, 290, 334, 375, 425, 482, 519, 569, 626, 663}
var _atm_logger_line_nums = [...]int{7, 11, 15, 56, 60, 64, 21, 24, 27, 38, 42, 46, 53, 57, 61}

const _atm_logger_detail = "some details here %vand here are some more details... %ssome details here %vand here are some more details... %ssome details here %vand here are some more details... %ssome details here %vand here are some more details... %s"

var _atm_logger_detail_index = [...]uint8{0, 20, 56, 56, 76, 112, 112, 112, 112, 112, 132, 168, 168, 188, 224, 224}
var _atm_logger_tab_counts = [...]int{2, 2, 2, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2}

func idxToRule(i int) string {
	if i >= len(_atm_logger_index)-1 {
		return strconv.FormatInt(int64(i), 10)
	}
	return _atm_logger_name[_atm_logger_index[i]:_atm_logger_index[i+1]]
}

func idxToDetail(i int) string {
	if i >= len(_atm_logger_detail_index)-1 {
		return strconv.FormatInt(int64(i), 10)
	}
	return _atm_logger_detail[_atm_logger_detail_index[i]:_atm_logger_detail_index[i+1]]
}

func lineNumToIndex(i int) int {
	k := searchInts(_atm_logger_line_nums, i)
	if _atm_logger_line_nums[k-1] == i-1 {
		return k - 1
	}
	return -1
}

func searchInts(a [15]int, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

// GetRule takes the line number at runtime and converts it to the nearest rule comment above it
func GetRule(runtimeLine int) string {
	return idxToRule(lineNumToIndex(runtimeLine))
}

type logger struct {
	RuntimeLines []int
	TitleArgs    [][]interface{}
	DetailArgs   [][]interface{}
	InitFunc     string
	InitFuncInt  int
}

func (l *logger) newLogger() {
	l.InitFunc = getCallerFunc()

}

func getCallerFunc() string {

	pcs := make([]uintptr, 10)
	n := runtime.Callers(2, pcs)
	pcs = pcs[:n]

	//frameLen := 0
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		log.Printf("frame.Function: %s", frame.Function)
		return frame.Function
	}

	return ""
}

func (l *logger) SetTitle(args ...interface{}) *logger {

	var previousLines []int
	// _, _, line, _ := runtime.Caller(1)
	// l.RuntimeLines = append(l.RuntimeLines, line)
	// l.TitleArgs = append(l.TitleArgs, args)

	pcs := make([]uintptr, 10)
	n := runtime.Callers(1, pcs)
	pcs = pcs[:n]

	//frameLen := 0
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		ff := frame.Function
		previousLines = append(previousLines, frame.Line)
		if ff == l.InitFunc {

			break
		}

	}
	//l.FrameLen = append(l.FrameLen, frameLen)
	l.RuntimeLines = append(l.RuntimeLines, sum(previousLines))
	l.TitleArgs = append(l.TitleArgs, args)
	return l
}

func (l *logger) SetDetail(args ...interface{}) {
	l.DetailArgs = append(l.DetailArgs, args)
}

func (l *logger) GetSummaryAll() RuleData {

	var rs RuleData
	runtimeIdx := 0
	nextTriggeredIdx := lineNumToIndex(l.RuntimeLines[runtimeIdx])
	lastTab := 0
	firstTab := 0
	var tabNumCount = make(map[int]int, 0)
	for k, _ := range _atm_logger_line_nums {

		tab := _atm_logger_tab_counts[k]

		// log.Printf("\n\nidx: %!d(MISSING)", k)
		// log.Printf("tab: %!d(MISSING)", tab)
		// log.Printf("nextTriggeredIdx: %!d(MISSING)", nextTriggeredIdx)
		if firstTab == 0 {
			firstTab = tab
		}
		rd := RuleData{
			Title:     idxToRule(k),
			Detail:    idxToDetail(k),
			HasDetail: len(idxToDetail(k)) > 0,
			TabNum:    tab,
		}

		if nextTriggeredIdx == k {

			rd.Triggered = true
			rd.Open = true
			rd.Title = fmt.Sprintf(idxToRule(k), l.TitleArgs[runtimeIdx]...)
			rd.Detail = fmt.Sprintf(idxToDetail(k), l.DetailArgs[runtimeIdx]...)
			runtimeIdx++
			if runtimeIdx < len(l.RuntimeLines) {
				nextTriggeredIdx = lineNumToIndex(l.RuntimeLines[runtimeIdx])

			}

		}
		tabDiff := tab - firstTab
		//log.Printf("tabDiff: %!d(MISSING)", tabDiff)

		if _, ok := tabNumCount[tabDiff]; ok {
			if lastTab > tab {
				// reset tabNumCount for anything > tabDiff
				for k, _ := range tabNumCount {
					if k > tabDiff {
						delete(tabNumCount, k)
					}
				}
			}

			tabNumCount[tabDiff] = tabNumCount[tabDiff] + 1

		} else {
			tabNumCount[tabDiff] = 0
		}
		//log.Printf("tabNumCount: %!d(MISSING)", tabNumCount)
		rs.AppendChild(rd, tabNumCount, tabDiff)
		lastTab = tab
	}

	return rs

}

func (l *logger) GetSummaryTriggered() RuleData {

	return RuleData{}
}

func sum(s []int) int {
	var tot int
	for _, v := range s {
		tot += v
	}
	return tot
}

type RuleData struct {
	Title     string
	HasDetail bool
	Detail    string
	TabNum    int
	Triggered bool
	Open      bool
	Children  []RuleData
}

func (rs *RuleData) AppendChild(rd RuleData, tabNameCount map[int]int, lastTab int) {
	switch lastTab {
	case 0:
		rs.Children = append(rs.Children, rd)
	case 1:
		rs.Children[tabNameCount[0]].Children = append(rs.Children[tabNameCount[0]].Children, rd)
	case 2:
		rs.Children[tabNameCount[0]].Children[tabNameCount[1]].Children = append(rs.Children[tabNameCount[0]].Children[tabNameCount[1]].Children, rd)
	}
}
