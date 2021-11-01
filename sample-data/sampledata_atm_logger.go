package sampledata

import "fmt"
import "strconv"
import "log"
import "runtime"

const _atm_logger_name = "a (%v)lessn than  b(%v), do somethinga (%v) greater than  b(%v), do something elsea (%v) == b (%v), do somethingnested func  a (%v)lessn than  b(%v), do somethingnested func a (%v) greater than  b(%v), do something elsenested a (%v) == b (%v), do somethingc (%v) < d (%v), do something elsec (%v) > d (%v), do something else once morec (%v) == d (%v), do something else again"

var _atm_logger_index = [...]uint16{0, 37, 82, 112, 162, 219, 256, 290, 334, 375}
var _atm_logger_line_nums = [...]int{7, 11, 15, 56, 60, 64, 21, 24, 27}
var _atm_logger_runtime_line_nums = [...]int{9, 13, 16, 58, 62, 65, 22, 25, 28}

const _atm_logger_detail = "some details here %vand here are some more details... %ssome details here %vand here are some more details... %s"

var _atm_logger_detail_index = [...]uint8{0, 20, 56, 56, 76, 112, 112, 112, 112, 112}
var _atm_logger_tab_counts = [...]int{2, 2, 2, 3, 3, 3, 3, 3, 3}

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
	k := searchInts(_atm_logger_runtime_line_nums, i)
	if k > 0 {
		return k
	}
	return -1

}

func searchInts(a [9]int, x int) int {
	for k, v := range a {

		if v == x {
			return k
		}
	}
	return -1
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

func newLogger() *logger {

	return &logger{
		InitFunc: getCallerFunc(),
	}

}

func getCallerFunc() string {

	pcs := make([]uintptr, 10)
	n := runtime.Callers(3, pcs)
	pcs = pcs[:n]

	//frameLen := 0
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		//log.Printf("frame.Function: %!v(MISSING)", frame.Function)
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
	n := runtime.Callers(2, pcs)
	pcs = pcs[:n]

	//frameLen := 0
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		// ffSlice := strings.Split(frame.Function, "/")
		// ff := strings.TrimLeft(ffSlice[len(ffSlice)-1], "sample-data.")

		log.Printf("frame func: %!v(MISSING) line: %!v(MISSING), initfunc: %!v(MISSING) ", frame.Function, frame.Line, l.InitFunc)
		previousLines = append(previousLines, frame.Line)
		if frame.Function == l.InitFunc {

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
	log.Printf("l.RuntimeLines: %!v(MISSING)", l.RuntimeLines)
	nextTriggeredIdx := lineNumToIndex(l.RuntimeLines[runtimeIdx])
	lastTab := 0
	firstTab := 0
	var tabNumCount = make(map[int]int, 0)
	for k, _ := range _atm_logger_line_nums {

		tab := _atm_logger_tab_counts[k]

		// log.Printf("\n\nidx: %!d(MISSING) lineNum: %!d(MISSING)", k, ln)
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
		//log.Printf("tabDiff: %!!(MISSING)d(MISSING)", tabDiff)

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
		//log.Printf("tabNumCount: %!!(MISSING)d(MISSING)", tabNumCount)
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
