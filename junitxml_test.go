package junitxml

import (
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"
	"time"
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
	var (
		testDuration = time.Nanosecond * 10
	)
	ju := &JUnitXML{}
	good := ju.Suite("all good")
	good.Case("alpha")
	good.Case("beta")
	good.Case("gamma")
	mixed := ju.Suite("mixed")
	success := mixed.Case("success")
	success.Success(testDuration)
	bad := mixed.Case("bad")
	bad.Fail("once", testDuration)
	bad.Fail("twice", testDuration)
	mixed.Case("ugly").Abort(errors.New("buggy"), testDuration)
	ju.Suite("fast").Fail("fail early", time.Duration(0))
	skipped := mixed.Case("skipped")
	skipped.Skip("skipped")

	return ju
}
