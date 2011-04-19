package main

import "g/xml"
import "testing"
import __regexp__ "regexp"

var tests = []testing.InternalTest{
	{"xml.TestCursor", xml.TestCursor},
	{"xml.TestChildrenString", xml.TestChildrenString},
	{"xml.TestChildrenSlice", xml.TestChildrenSlice},
}
var benchmarks = []testing.InternalBenchmark{ //
}

func main() {
	testing.Main(__regexp__.MatchString, tests)
	testing.RunBenchmarks(__regexp__.MatchString, benchmarks)
}
