package sampledata

import "strconv"
import "sort"

const _atm_logger_name = "a lessn than  b, do somethinga greater than  b, do something elsea (%s) == b (%s), do something"

var _atm_logger_index = [...]uint8{0, 29, 65, 95}
var _atm_logger_line_nums = [...]int{11, 14, 17}

func IdxToRule(i int) string {
	if i >= len(_atm_logger_index)-1 {
		return strconv.FormatInt(int64(i), 10)
	}
	return _atm_logger_name[_atm_logger_index[i]:_atm_logger_index[i+1]]
}

func LineNumToIndex(i int) int {
	k := searchInts(_atm_logger_line_nums, i)
	if _atm_logger_line_nums[k] == i {
		return k + 1
	}
	return -1
}

func searchInts(a [3]int, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}
