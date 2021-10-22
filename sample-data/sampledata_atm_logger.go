package sampledata

import (
	"runtime"
	"sort"
	"strconv"
)

const _atm_logger_name = "a (%v)lessn than  b(%v), do somethinga (%v) greater than  b(%v), do something elsea (%v) == b (%v), do somethingc (%v) < d (%v), do something elsec (%v) > d (%v), do something else once morec (%v) == d (%v), do something else again"

var _atm_logger_index = [...]uint8{0, 37, 82, 112, 146, 190, 231}
var _atm_logger_line_nums = [...]int{7, 11, 15, 19, 22, 25}

const _atm_logger_detail = "some details here %vand here are some more details... %s"

var _atm_logger_detail_index = [...]uint8{0, 20, 56, 56, 56, 56, 56}
var _atm_logger_tab_counts = [...]int{2, 2, 2, 3, 3, 3}

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

func searchInts(a [6]int, x int) int {
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
}

func (l *logger) SetTitle(args ...interface{}) *logger {
	_, _, line, _ := runtime.Caller(1)
	l.RuntimeLines = append(l.RuntimeLines, line)
	l.TitleArgs = append(l.TitleArgs, args)
	return l
}

func (l *logger) SetDetail(args ...interface{}) {
	l.DetailArgs = append(l.DetailArgs, args)
}

func (l *logger) GetSummaryAll() RuleSummary {

	var rs RuleSummary
	runtimeIdx := 0
	nextTriggeredIdx := lineNumToIndex(l.RuntimeLines[runtimeIdx])
	// _ = nextTriggered
	// var _log_summary RuleSummary
	// var currentLevel int
	lastTab := 0
	for k, _ := range _atm_logger_line_nums {
		tab := _atm_logger_tab_counts[k]
		rd := RuleData{
			Title:     idxToRule(k),
			Detail:    idxToDetail(k),
			HasDetail: len(idxToDetail(k)) > 0,
			TabNum:    tab,
		}

		if nextTriggeredIdx == k && runtimeIdx < len(l.RuntimeLines) {
			rd.Triggered = true
			rd.Open = true
			nextTriggeredIdx = lineNumToIndex(l.RuntimeLines[runtimeIdx])
			runtimeIdx++
		}

		if tab > lastTab {
			// child
		} else if tab < lastTab {
			// parent

		} else {
			// sibling
		}

		rs = append(rs, rd)

		// // set title, detail

		// if k == nextTriggeredIdx {
		// 	// set to triggered
		// }
		// //_atm_logger_tab_counts[k]

	}

	return rs
}

func (l *logger) GetSummaryTriggered() RuleSummary {

	return RuleSummary{}
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

type RuleSummary []RuleData
