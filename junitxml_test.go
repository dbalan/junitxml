package junitxml

import (
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"
)

func TestJunitOutput(t *testing.T) {
	var buf bytes.Buffer
	FakeTestSuites().WriteXML(&buf)
	output := buf.Bytes()

	os.WriteFile("fakeresults.xml", output, 0644) // to validate later

	reTop := regexp.MustCompile(`(?s)^<testsuites\W.*</testsuites>$`)
	reSuites := regexp.MustCompile(`(?s)<testsuite .*?</testsuite>`)
	reCases := regexp.MustCompile(`(?s)<testcase .*?</testcase>`)

	if !reTop.Match(output) {
		t.Errorf("JUnit output has no outer <testsuites>\n")
	}
	suites := reSuites.FindAll(output, -1)
	if len(suites) != 3 {
		t.Errorf("JUnit output had %d testsuite elements; expected 3\n", len(suites))
	}
	cases := reCases.FindAll(output, -1)
	if len(cases) != 8 {
		t.Errorf("JUnit output had %d testcase; expected 8\n", len(cases))
	}
}

func FakeTestSuites() *JUnitXML {
	ju := &JUnitXML{}
	good := ju.Suite("all good")
	good.Case("alpha")
	good.Case("beta")
	good.Case("gamma")
	mixed := ju.Suite("mixed")
	mixed.Case("good")
	bad := mixed.Case("bad")
	bad.Fail("once", "0.01000")
	bad.Fail("twice", "0.01000")
	mixed.Case("ugly").Abort(errors.New("buggy"), "0.1000")
	ju.Suite("fast").Fail("fail early", "0.0")
	skipped := mixed.Case("skipped")
	skipped.Skip("skipped")

	return ju
}
