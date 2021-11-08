package sampledata

import "log"

//go:generate atm-logger -func=CompareInt
func CompareInt(a, b, c, d int) *Logger {
	l := newLogger()

	// RULE: Start compare int
	l.SetTitle().SetDetail()

	if a < b {
		// RULE: a (%v)lessn than  b(%v), do something HINT: a<b
		// some details here %v
		l.SetTitle(a, b).SetDetail(8)
	} else if a > b {
		// RULE: a (%v) greater than  b(%v), do something else HINT: a>b
		// and here are some more details... %s
		l.SetTitle(a, b).SetDetail("some more detail")
	} else {
		// RULE: a (%v) == b (%v), do something HINT: else
		l.SetTitle(a, b).SetDetail()

		ret := CompareInt2(a, b, l)

		log.Printf("ret: %v", ret)

		if c < d {
			// RULE: c (%v) < d (%v), do something else
			l.SetTitle(c, d).SetDetail()
		} else if c > d {
			// RULE: c (%v) > d (%v), do something else once more
			l.SetTitle(c, d).SetDetail()
		} else {
			// RULE: c (%v) == d (%v), do something else again
			l.SetTitle(c, d).SetDetail()
		}
	}

	return l

}

func CompareInt2(a, b int, l *Logger) int {

	// RULE: Start compare int 2 HINT: CompareInt2
	l.SetTitle().SetDetail()

	if a < b {
		// RULE: nested func  a (%v)lessn than  b(%v), do something
		// some details here %v
		l.SetTitle(a, b).SetDetail(8)
	} else if a > b {
		// RULE: nested func a (%v) greater than  b(%v), do something else
		// and here are some more details... %s
		l.SetTitle(a, b).SetDetail("some nested detail")
	} else {
		// RULE: nested a (%v) == b (%v), do something
		l.SetTitle(a, b).SetDetail()
	}

	return 0

}
