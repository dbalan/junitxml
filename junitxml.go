package junitxml

import "encoding/xml"

type JUnitXML struct {
	XMLName xml.Name     `xml:"testsuites"`
	Suites  []*TestSuite `xml:"testsuite"`
}
type TestSuite struct {
	Name         string      `xml:"name,attr"`
	TestCount    int         `xml:"tests,attr"`
	FailureCount int         `xml:"failures,attr"`
	ErrorCount   int         `xml:"errors,attr"`
	Cases        []*TestCase `xml:"testcase"`
}
type TestCase struct {
	Name     string   `xml:"name,attr"`
	Failures []string `xml:"failure,omitempty"`
	Error    string   `xml:"error,omitempty"`
}

func (j *JUnitXML) Suite(name string) *TestSuite {
	ts := &TestSuite{Name: name}
	j.Suites = append(j.Suites, ts)
	return ts
}

func (ts *TestSuite) Case(name string) *TestSuite {
	ts.TestCount++
	tc := &TestCase{Name: name}
	ts.Cases = append(ts.Cases, tc)
	return ts
}

func (ts *TestSuite) Fail(msg string) {
	ts.FailureCount++
	if len(ts.Cases) == 0 {
		ts.Case("unknown")
	}
	curt := ts.Cases[len(ts.Cases)-1]
	curt.Failures = append(curt.Failures, msg)
}

func (ts *TestSuite) Abort(msg string) {
	ts.ErrorCount++
	if len(ts.Cases) == 0 {
		ts.Case("unknown")
	}
	curt := ts.Cases[len(ts.Cases)-1]
	curt.Error = msg
}
