package flags

import (
	"flag"
	"os"
	"strings"
	"testing"

	. "github.com/pingcap/check"
)

func Test(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testFlagSuite{})

type testFlagSuite struct{}

func (s *testFlagSuite) TestSetFlagsFromEnv(c *C) {
	fs := flag.NewFlagSet("test1", flag.ExitOnError)
	fs.String("f1", "", "")
	fs.String("f2", "", "")
	fs.String("f3", "", "")
	fs.Parse([]string{})

	os.Clearenv()
	// 1. flag is set with env vars
	os.Setenv("TEST_F1", "abc")
	// 2. flag is set by command-line args
	mustSuccess(c, fs.Set("f2", "xyz"))
	// 3. command-line flags take precedence over env vars
	os.Setenv("TEST_F3", "123")
	mustSuccess(c, fs.Set("f3", "789"))

	// before
	for fl, expected := range map[string]string{
		"f1": "",
		"f2": "xyz",
		"f3": "789",
	} {
		c.Assert(fs.Lookup(fl).Value.String(), Equals, expected)
	}

	mustSuccess(c, SetFlagsFromEnv("TEST", fs))

	// after
	for fl, expected := range map[string]string{
		"f1": "abc",
		"f2": "xyz",
		"f3": "789",
	} {
		c.Assert(fs.Lookup(fl).Value.String(), Equals, expected)
	}
}

func (s *testFlagSuite) TestSetFlagsFromEnvMore(c *C) {
	fs := flag.NewFlagSet("test2", flag.ExitOnError)
	fs.String("str", "", "")
	fs.Int("int", 0, "")
	fs.Bool("bool", false, "")
	fs.String("a-hyphen", "", "")
	fs.String("lowercase", "", "")
	fs.Parse([]string{})

	os.Clearenv()
	os.Setenv("TEST_STR", "ijk")
	os.Setenv("TEST_INT", "654")
	os.Setenv("TEST_BOOL", "1")
	os.Setenv("TEST_A_HYPHEN", "foo")
	os.Setenv("TEST_lowertest", "bar")

	mustSuccess(c, SetFlagsFromEnv("TEST", fs))

	for fl, expected := range map[string]string{
		"str":       "ijk",
		"int":       "654",
		"bool":      "true",
		"a-hyphen":  "foo",
		"lowercase": "",
	} {
		c.Assert(fs.Lookup(fl).Value.String(), Equals, expected)
	}
}

func (s *testFlagSuite) TestSetFlagsFromEnvBad(c *C) {
	fs := flag.NewFlagSet("test3", flag.ExitOnError)
	fs.Int("num", 0, "")
	fs.Parse([]string{})

	os.Clearenv()
	os.Setenv("TEST_NUM", "abc123")

	mustFail(c, SetFlagsFromEnv("TEST", fs))
}

func (s *testFlagSuite) TestURLStrsFromFlag(c *C) {
	urlv, err := NewURLsValue("http://127.0.0.1:1234")
	c.Assert(err, IsNil)

	fs := flag.NewFlagSet("testUrlFlag", flag.ExitOnError)
	fs.Var(urlv, "urls", "")
	fs.Parse([]string{})

	urls := "http://192.168.1.1:1234,http://192.168.1.2:1234,http://192.168.1.3:1234"
	mustSuccess(c, fs.Set("urls", urls))
	c.Assert(strings.Join(URLStrsFromFlag(fs, "urls"), ","), Equals, urls)
}

func mustSuccess(c *C, err error) {
	c.Assert(err, IsNil)
}

func mustFail(c *C, err error) {
	c.Assert(err, NotNil)
}
