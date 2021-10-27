package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {

	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedName | packages.NeedSyntax, Tests: false}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		//os.Exit(1)
	}

	// Print the names of the source files
	// for each package listed on the command line.
	for _, pkg := range pkgs {

		if len(pkg.GoFiles) == 0 {
			log.Printf("no go filesin pkg: %v", pkg.Name)
			continue
		}

		buf := new(bytes.Buffer)

		cp := newCodeParser()

		log.Printf("package name: %s", pkg.Name)
		pathSlice := strings.Split(pkg.GoFiles[0], "/")
		pathSlice = pathSlice[:len(pathSlice)-1]
		pathSlice = append(pathSlice, fmt.Sprintf("%s_atm_logger.go", pkg.Name))

		fmt.Fprintf(buf, "package %s\n\n", pkg.Name)
		fmt.Fprintf(buf, "import \"fmt\"\n")
		fmt.Fprintf(buf, "import \"strconv\"\n")
		fmt.Fprintf(buf, "import \"sort\"\n")
		fmt.Fprintf(buf, "import \"log\"\n")
		fmt.Fprintf(buf, "import \"runtime\"\n\n")

		//expFuncMap := make(map[string]posRange, 0)
		//localFuncMap := make(map[string]posRange, 0)

		for _, file := range pkg.GoFiles {

			if strings.HasSuffix(file, "_test.go") {
				continue
			}
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
			if err != nil {
				log.Fatal(err)
			}

			ast.Inspect(node, func(n ast.Node) bool {
				// Find  Functions
				fn, ok := n.(*ast.FuncDecl)
				if ok {

					pr := posRange{
						OffsetStart: int64(fset.Position(fn.Body.Lbrace).Offset),
						StartLine:   fset.Position(fn.Body.Lbrace).Line,
						EndLine:     fset.Position(fn.Body.Rbrace).Line,
					}
					if fn.Name.IsExported() {

						//expFuncMap[fn.Name.Name] = pr
						//log.Printf("exp fname: %v", fn.Name)
						//log.Printf("file line %d", fset.Position(fn.Pos()).Line)
					}
					cp.localFuncMap[fn.Name.Name] = pr
					//log.Printf("non-exp fname: %v", fn.Name)
					//log.Printf("file line %d", fset.Position(fn.Pos()).Line)

					// only go one level down
					return false
				}
				return true

			})

		}

		log.Printf("completed ast inspect funcs")

		for _, file := range pkg.GoFiles {

			cp.reset()
			if strings.HasSuffix(file, "_test.go") {
				continue
			}
			fmt.Println(file)
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "open file: %v\n", err)
				os.Exit(1)

			}
			defer f.Close()

			//cp.setScanner(f)

			// tabs := make([]int, 0)
			// titles := make([]string, 0)
			// details := make([]string, 0)
			// lineNums := make([]int, 0)

			// call new func here

			if err := cp.parseCode(f, 0, 0, 0, 0); err != nil {
				log.Fatalf("error parsing code: %v", err)
			}

			if len(cp.titles) > 0 {
				declareIndexAndNameVar(buf, cp.titles, cp.lineNums)
				declareTabAndDetailVar(buf, cp.details, cp.tabs)
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, indexToRule)
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, numToIdx)
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, searchInts, len(cp.lineNums))
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, getRule)
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, logger)
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, summary)

				fmt.Fprintf(buf, "\n\n")
				writeAppendChildFunc(buf, maxTabs(cp.tabs))

			}

		}

		//log.Printf("bytes string: %v", buf.String())
		src, err := format.Source(buf.Bytes())
		if err != nil {
			log.Fatalf("error formatting code: %v", err)

		}

		err = ioutil.WriteFile(strings.Join(pathSlice, "/"), src, 0644)
		if err != nil {
			log.Fatalf("writing output: %s", err)
		}

	}
}

type codeParser struct {
	titles       []string
	details      []string
	tabs         []int
	lineNums     []int
	tabOffset    int
	localFuncMap map[string]posRange
	//scanner      *bufio.Scanner
	scanOffset     int64
	callerLineHist []int
	//	buf          *bytes.Buffer
}

func (p *codeParser) reset() {
	p.titles = []string{}
	p.details = []string{}
	p.tabs = []int{}
	p.lineNums = []int{}
	p.tabOffset = 0

}

func sum(s []int) int {
	var tot int
	for _, v := range s {
		tot += v
	}
	return tot
}

func newCodeParser() codeParser {

	return codeParser{
		localFuncMap: make(map[string]posRange, 0),

		//buf:          new(bytes.Buffer),
	}

}

func (p *codeParser) setScanner(f io.ReadSeeker) {

}

type posRange struct {
	OffsetStart int64
	StartLine   int
	EndLine     int
}

func (p *codeParser) parseCode(f io.ReadSeeker, offsetStart int64, parseLineOffsetStart, parseLineOffsetEnd, callerLine int) error {

	scanner := bufio.NewScanner(f)

	scanOffset := offsetStart
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		scanOffset += int64(advance)
		return
	}
	scanner.Split(scanLines)

	if callerLine > 0 {
		p.callerLineHist = append(p.callerLineHist, callerLine)
	}

	defer func() {
		p.tabOffset--
		if len(p.callerLineHist) > 0 {
			p.callerLineHist = p.callerLineHist[:len(p.callerLineHist)-1]
		}
		log.Printf("seek back to where we left off: %v", offsetStart)
		if _, err := f.Seek(offsetStart, 0); err != nil {
			log.Fatalf("couldn't seek back to offset start: %v", err)
		}
	}()

	// this tells us we just finished a rule
	var checkDetails bool
	var detailBuilder strings.Builder

	lineNum := 1
	if parseLineOffsetStart > lineNum {
		lineNum = parseLineOffsetStart
	}

	if _, err := f.Seek(offsetStart, 0); err != nil {
		return err
	}

	for scanner.Scan() {

		//log.Printf("linenum: %d", lineNum)
		lineText := scanner.Text()
		// skip the first line
		if lineNum < parseLineOffsetStart {
			continue
		}
		if parseLineOffsetEnd > 0 && lineNum >= parseLineOffsetEnd {
			break
		}
		if strings.HasPrefix(strings.TrimSpace(lineText), "// RULE:") {
			log.Printf("line: %d: %s", lineNum, strings.TrimPrefix(strings.TrimSpace(lineText), "// RULE: "))
			p.titles = append(p.titles, strings.TrimPrefix(strings.TrimSpace(lineText), "// RULE: "))
			p.tabs = append(p.tabs, countTabs(lineText, p.tabOffset))
			p.lineNums = append(p.lineNums, lineNum+sum(p.callerLineHist))
			checkDetails = true

			// make map of hash to index

		} else if checkDetails && strings.HasPrefix(strings.TrimSpace(lineText), "//") {
			detailBuilder.WriteString(strings.TrimPrefix(strings.TrimSpace(lineText), "// "))

		} else if checkDetails {
			p.details = append(p.details, detailBuilder.String())
			checkDetails = false
			detailBuilder.Reset()
		} else if strings.HasPrefix(strings.TrimSpace(lineText), "//") {

		} else {

			// loop over package func decl names
			// check if this line calls any of them
			// TODO: should be a better way of doing this in the ast package
			for k, v := range p.localFuncMap {

				if strings.Contains(lineText, fmt.Sprintf("%s(", k)) && v.StartLine != lineNum {

					log.Printf("%d found func: %v parsing from line %d to %d", lineNum, k, v.StartLine, v.EndLine)
					// go to this func and scan it next
					//currentOffset := pos
					p.tabOffset++

					if err := p.parseCode(f, v.OffsetStart, v.StartLine, v.EndLine, lineNum); err != nil {
						log.Fatalf("error parsing code %v tab offset: %d", err, p.tabOffset)
					}

				}
			}

		}
		lineNum++
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return err
	}

	return nil

}

func maxTabs(tabs []int) int {
	maxTabs := 0
	for k, v := range tabs {
		if k == 0 {
			maxTabs = v
		} else if v > tabs[k-1] {
			maxTabs = v
		}
	}
	return maxTabs
}

func countTabs(s string, tabOffset int) int {

	for k, v := range s {
		if v != '\t' {
			return k + tabOffset
		}
	}
	return 0 + tabOffset
}

// // custom stringer
// // Code generated by stringer -type Pill pill.go; DO NOT EDIT.

// package painkiller

// import "fmt"

// const _Log_name = "PlaceboAspirinIbuprofenParacetamol"

// var _Log_index = [...]uint8{0, 7, 14, 23, 34}

// func (i Log) String() string {
//     if i < 0 || i+1 >= Log(len(_Log_index)) {
//         return fmt.Sprintf("Log(%d)", i)
//     }
//     return _Log_name[_Log_index[i]:_Log_index[i+1]]
// }v

// Helpers

// usize returns the number of bits of the smallest unsigned integer
// type that will hold n. Used to create the smallest possible slice of
// integers to use as indexes into the concatenated strings.
func usize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	default:
		// 2^32 is enough constants for anyone.
		return 32
	}
}

// writeAppendChildFunc writes the method to append rules to the RuleData.Children field
func writeAppendChildFunc(b *bytes.Buffer, maxTab int) {
	funcStr := createAppendChildFunc(maxTab)
	fmt.Fprintf(b, "func %s\n", funcStr)
}

func createAppendChildFunc(maxTab int) string {
	b := new(bytes.Buffer)

	b.WriteString(`(rs *RuleData) AppendChild(rd RuleData, tabNameCount map[int]int, lastTab int) {
	switch lastTab {
`)
	for i := 0; i < maxTab; i++ {
		fmt.Fprintf(b, "case %d:\n", i)
		fmt.Fprintf(b, "rs.Children")
		for j := 0; j < i; j++ {
			fmt.Fprintf(b, "[tabNameCount[%d]].Children", j)
		}
		fmt.Fprintf(b, " = append(rs.Children")
		for j := 0; j < i; j++ {
			fmt.Fprintf(b, "[tabNameCount[%d]].Children", j)
		}
		fmt.Fprintf(b, ",  rd)\n")
	}

	fmt.Fprintf(b, "}\n")
	fmt.Fprintf(b, "}")
	return b.String()
}

// declareIndexAndNameVar is the single-run version of declareIndexAndNameVars
func declareTabAndDetailVar(b *bytes.Buffer, details []string, tabs []int) {
	index, name, tabStr := createTabAndDetailDecl(details, tabs)
	fmt.Fprintf(b, "const %s\n", name)
	fmt.Fprintf(b, "var %s\n", index)
	fmt.Fprintf(b, "var %s\n", tabStr)

	//fmt.Fprintf(b, stringOneRun, "_atm_logger_name")

}

// createTabAndDetailDecl returns the pair of declarations for the run. The caller will add "const" and "var".
func createTabAndDetailDecl(comments []string, lineNums []int) (string, string, string) {
	b := new(bytes.Buffer)
	indexes := make([]int, len(comments))
	for i := range comments {
		b.WriteString(comments[i])
		indexes[i] = b.Len()
	}
	nameConst := fmt.Sprintf("_atm_logger_detail = %q", b.String())
	nameLen := b.Len()
	b.Reset()
	fmt.Fprintf(b, "_atm_logger_detail_index = [...]uint%d{0, ", usize(nameLen))
	for i, v := range indexes {
		if i > 0 {
			fmt.Fprintf(b, ", ")
		}
		fmt.Fprintf(b, "%d", v)
	}
	fmt.Fprintf(b, "}")
	index := b.String()
	b.Reset()
	fmt.Fprintf(b, "_atm_logger_tab_counts = [...]int{")
	for i, v := range lineNums {
		if i > 0 {
			fmt.Fprintf(b, ", ")
		}
		fmt.Fprintf(b, "%d", v)
	}
	fmt.Fprintf(b, "}")
	return index, nameConst, b.String()
}

// declareIndexAndNameVar is the single-run version of declareIndexAndNameVars
func declareIndexAndNameVar(b *bytes.Buffer, comments []string, lineNums []int) {
	index, name, lineNumStr := createIndexAndNameDecl(comments, lineNums)
	fmt.Fprintf(b, "const %s\n", name)
	fmt.Fprintf(b, "var %s\n", index)
	fmt.Fprintf(b, "var %s\n", lineNumStr)

	//fmt.Fprintf(b, stringOneRun, "_atm_logger_name")

}

// createIndexAndNameDecl returns the pair of declarations for the run. The caller will add "const" and "var".
func createIndexAndNameDecl(comments []string, lineNums []int) (string, string, string) {
	b := new(bytes.Buffer)
	indexes := make([]int, len(comments))
	for i := range comments {
		b.WriteString(comments[i])
		indexes[i] = b.Len()
	}
	nameConst := fmt.Sprintf("_atm_logger_name = %q", b.String())
	nameLen := b.Len()
	b.Reset()
	fmt.Fprintf(b, "_atm_logger_index = [...]uint%d{0, ", usize(nameLen))
	for i, v := range indexes {
		if i > 0 {
			fmt.Fprintf(b, ", ")
		}
		fmt.Fprintf(b, "%d", v)
	}
	fmt.Fprintf(b, "}")
	index := b.String()
	b.Reset()
	fmt.Fprintf(b, "_atm_logger_line_nums = [...]int{")
	for i, v := range lineNums {
		if i > 0 {
			fmt.Fprintf(b, ", ")
		}
		fmt.Fprintf(b, "%d", v)
	}
	fmt.Fprintf(b, "}")
	return index, nameConst, b.String()
}

const indexToRule = `func idxToRule(i int) string {
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

`

const numToIdx = `func lineNumToIndex(i int) int {
    k := searchInts(_atm_logger_line_nums, i)
    if _atm_logger_line_nums[k-1] == i-1 {	
	return k - 1
    }
    return -1
}
`

const searchInts = `func searchInts(a [%d]int, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}
`

const getRule = `// GetRule takes the line number at runtime and converts it to the nearest rule comment above it
func GetRule(runtimeLine int) string {
	return idxToRule(lineNumToIndex(runtimeLine))
}
`

const logger = `type logger struct {
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

	frameLen := 0
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		log.Printf("frame.Function: %v", frame.Function)
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

	frameLen := 0
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

		// log.Printf("\n\nidx: %d", k)
		// log.Printf("tab: %d", tab)
		// log.Printf("nextTriggeredIdx: %d", nextTriggeredIdx)
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
		//log.Printf("tabDiff: %d", tabDiff)

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
		//log.Printf("tabNumCount: %d", tabNumCount)
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
`

const summary = `type RuleData struct{
        Title     string
        HasDetail bool
        Detail    string
        TabNum    int
        Triggered bool
        Open      bool
        Children []RuleData
}
`

// // AppendChild appends child elements to RuleSummary at level of tab
// func (rs *RuleData) AppendChild(rd RuleData, tabs []int) {
// 	switch len(tabs) {
// 	case 0:
// 		rs.Children = append(rs.Children, rd)
// 	case 1:
// 		rs.Children[tabs[0]].Children = append(rs.Children[tabs[0]].Children, rd)
// 	case 2:
// 		rs.Children[tabs[0]].Children[tabs[1]].Children = append(rs.Children[tabs[0]].Children[tabs[1]].Children, rd)

// 	}

// }

// func Search(n int, f func(int) bool) int {
// 	// Define f(-1) == false and f(n) == true.
// 	// Invariant: f(i-1) == false, f(j) == true.
// 	i, j := 0, n
// 	for i < j {
// 		h := int(uint(i+j) >> 1) // avoid overflow when computing h
// 		// i ≤ h < j
// 		if !f(h) {
// 			i = h + 1 // preserves f(i-1) == false
// 		} else {
// 			j = h // preserves f(j) == true
// 		}
// 	}
// 	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
// 	return i
// }

// Convenience wrappers for common cases.

// SearchInts searches for x in a sorted slice of ints and returns the index
// as specified by Search. The return value is the index to insert x if x is
// not present (it could be len(a)).
// The slice must be sorted in ascending order.
//
// func SearchInts(a []int, x int) int {
// 	return Search(len(a), func(i int) bool { return a[i] >= x })
// }
