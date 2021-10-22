package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {

	log.Printf("hello world!")

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

		log.Printf("package name: %s", pkg.Name)
		pathSlice := strings.Split(pkg.GoFiles[0], "/")
		pathSlice = pathSlice[:len(pathSlice)-1]
		pathSlice = append(pathSlice, fmt.Sprintf("%s_atm_logger.go", pkg.Name))

		fmt.Fprintf(buf, "package %s\n\n", pkg.Name)
		//fmt.Fprintf(buf, "import \"fmt\"\n\n")
		fmt.Fprintf(buf, "import \"strconv\"\n")
		fmt.Fprintf(buf, "import \"sort\"\n")
		fmt.Fprintf(buf, "import \"runtime\"\n\n")

		for _, file := range pkg.GoFiles {
			//fmt.Println(file)
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "open file: %v\n", err)
				os.Exit(1)

			}
			defer f.Close()

			log.Println("got here 1")
			tabs := make([]int, 0)
			titles := make([]string, 0)
			details := make([]string, 0)
			lineNums := make([]int, 0)

			// this tells us we just finished a rule
			var checkDetails bool

			var detailBuilder strings.Builder

			scanner := bufio.NewScanner(f)
			lineNum := 1
			// optionally, resize scanner's capacity for lines over 64K, see next example
			for scanner.Scan() {
				lineText := scanner.Text()
				if strings.HasPrefix(strings.TrimSpace(lineText), "// RULE:") {
					log.Printf("line: %d: %s", lineNum, strings.TrimPrefix(strings.TrimSpace(lineText), "// RULE: "))
					titles = append(titles, strings.TrimPrefix(strings.TrimSpace(lineText), "// RULE: "))
					tabs = append(tabs, countTabs(lineText))
					lineNums = append(lineNums, lineNum)
					checkDetails = true

					// make map of hash to index

				} else if checkDetails && strings.HasPrefix(strings.TrimSpace(lineText), "//") {
					detailBuilder.WriteString(strings.TrimPrefix(strings.TrimSpace(lineText), "// "))

				} else if checkDetails {
					details = append(details, detailBuilder.String())
					checkDetails = false
					detailBuilder.Reset()
				}
				lineNum++
			}

			if err := scanner.Err(); err != nil && err != io.EOF {
				log.Fatal(err)
			}

			if len(titles) > 0 {
				declareIndexAndNameVar(buf, titles, lineNums)
				declareTabAndDetailVar(buf, details, tabs)
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, indexToRule)
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, numToIdx)
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, searchInts, len(lineNums))
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, getRule)
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, logger)
				fmt.Fprintf(buf, "\n\n")
				_, err = fmt.Fprintf(buf, summary)

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

func countTabs(s string) int {

	for k, v := range s {
		if v != '\t' {
			return k
		}
	}
	return 0
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
	RuntimeLines  []int
	TitleArgs     [][]interface{}
        DetailArgs    [][]interface{}
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
        
	runtimeIdx := 0
	nextTriggeredIdx := lineNumToIndex(l.RuntimeLines[runtimeIdx])
	// _ = nextTriggered
	// var _log_summary RuleSummary
	// var currentLevel int
	for k, _ := range _atm_logger_line_nums {
		rd := RuleData{
			Title:     idxToRule(k),
			Detail:    idxToDetail(k),
			HasDetail: len(idxToDetail(k)) > 0,
			TabNum:    _atm_logger_tab_counts[k],
		}

		if nextTriggeredIdx == k && runtimeIdx < len(l.RuntimeLines) {
			rd.Triggered = true

			nextTriggeredIdx = lineNumToIndex(l.RuntimeLines[runtimeIdx])
			runtimeIdx++
		}

		// // set title, detail

		// if k == nextTriggeredIdx {
		// 	// set to triggered
		// }
		// //_atm_logger_tab_counts[k]

	}

	return RuleSummary{}
}

func (l *logger) GetSummaryTriggered() RuleSummary {


	return RuleSummary{}
}
`

const summary = `type RuleData struct{
        Title     string
        HasDetail bool
        Detail    string
        TabNum    int
        Triggered bool
        Children []RuleData
}

type RuleSummary []RuleData
`

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
