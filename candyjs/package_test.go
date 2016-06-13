package candyjs

import (
	. "gopkg.in/check.v1"
)

func (s *CandySuite) TestRegisterPackagePusher(c *C) {
	fn := func(ctx *Context) {}
	RegisterPackagePusher("foo", fn)

	c.Assert(pushers, HasLen, 1)
}

func (s *CandySuite) TestPushGlobalPackage(c *C) {
	fn := func(ctx *Context) {
		ctx.PushString("qux")
	}

	RegisterPackagePusher("foo", fn)
	c.Assert(s.ctx.PushGlobalPackage("foo", "bar"), IsNil)

	s.ctx.PevalString(`store(bar)`)
	c.Assert(s.stored, Equals, "qux")
}

func (s *CandySuite) TestPushGlobalPackage_NotFound(c *C) {
	c.Assert(s.ctx.PushGlobalPackage("qux", "qux"), Equals, ErrPackageNotFound)
}
