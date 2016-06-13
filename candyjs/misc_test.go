package candyjs

import (
	. "gopkg.in/check.v1"
)

func (s *CandySuite) TestIsExported(c *C) {
	c.Assert(isExported("Foo"), Equals, true)
	c.Assert(isExported("foo"), Equals, false)
}

func (s *CandySuite) TestNameToJavaScript(c *C) {
	c.Assert(nameToJavaScript("FooQux"), Equals, "fooQux")
	c.Assert(nameToJavaScript("FOOQux"), Equals, "fooQux")
	c.Assert(nameToJavaScript("Foo"), Equals, "foo")
	c.Assert(nameToJavaScript("FOO"), Equals, "foo")
}

func (s *CandySuite) TestNameToGo(c *C) {
	c.Assert(nameToGo("fooQux"), DeepEquals, []string{"FooQux", "FOOQux"})
	c.Assert(nameToGo("FooQux"), DeepEquals, []string(nil))
}
