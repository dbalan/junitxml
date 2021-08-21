package junitxml

import (
	"bytes"
	"encoding/xml"
	"os"
	"regexp"
	"testing"
)

func TestJunitOutput(t *testing.T) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Encode(FakeTestSuites())
	output := buf.Bytes()

	os.WriteFile("fakeresults.xml", output, 0666) // to validate later

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
	if len(cases) != 7 {
		t.Errorf("JUnit output had %d testcase; expected 7\n", len(cases))
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
	mixed.Case("bad").Fail("failed")
	mixed.Case("ugly").Abort("buggy")
	ju.Suite("fast").Fail("too soon")
	return ju
}
