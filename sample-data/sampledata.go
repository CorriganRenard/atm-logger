package sampledata

//go:generate atm-logger
func CompareInt(a, b, c, d int) *logger {
	l := &logger{}
	if a < b {
		// RULE: a (%v)lessn than  b(%v), do something
		// some details here %v
		l.SetTitle(a, b).SetDetail(8)
	} else if a > b {
		// RULE: a (%v) greater than  b(%v), do something else
		// and here are some more details... %s
		l.SetTitle(a, b).SetDetail("some more detail")
	} else {
		// RULE: a (%v) == b (%v), do something
		l.SetTitle(a, b)

		if c < d {
			// RULE: c (%v) < d (%v), do something else
			l.SetTitle(c, d)
		} else if c > d {
			// RULE: c (%v) > d (%v), do something else once more
			l.SetTitle(c, d)
		} else {
			// RULE: c (%v) == d (%v), do something else again
			l.SetTitle(c, d)
		}
	}
	return l

}
