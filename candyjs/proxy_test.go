package candyjs

import (
	"encoding/json"

	. "gopkg.in/check.v1"
)

func (s *CandySuite) TestProxy_Has(c *C) {
	c.Assert(p.has(&MyStruct{Int: 42}, "int"), Equals, true)
	c.Assert(p.has(&MyStruct{Int: 42}, "Int"), Equals, false)
}

func (s *CandySuite) TestProxy_Get(c *C) {
	v, err := p.get(&MyStruct{Int: 42}, "int", nil)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, 42)
}

func (s *CandySuite) TestProxy_GetUndefinedProperty(c *C) {
	v, err := p.get(&MyStruct{Int: 42}, "foo", nil)
	c.Assert(err, Equals, ErrUndefinedProperty)
	c.Assert(v, Equals, nil)
}

func (s *CandySuite) TestProxy_GetInternal(c *C) {
	v, err := p.get(&MyStruct{Int: 42}, "toJSON", nil)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, nil)
}

func (s *CandySuite) TestProxy_Set(c *C) {
	providers := [][]interface{}{
		{"int", nil, 0},
		{"int", 42.0, 42},
		{"int8", 42.0, int8(42)},
		{"int16", 42.0, int16(42)},
		{"int32", 42.0, int32(42)},
		{"int64", 42.0, int64(42)},
		{"uInt", 42.0, uint(42)},
		{"uInt8", 42.0, uint8(42)},
		{"uInt16", 42.0, uint16(42)},
		{"uInt32", 42.0, uint32(42)},
		{"uInt64", 42.0, uint64(42)},
		{"float32", 42.0, float32(42)},
	}

	for _, p := range providers {
		s.testProxy_Set(c, p[0], p[1], p[2])
	}
}

func (s *CandySuite) testProxy_Set(c *C, key, set, get interface{}) {
	t := &MyStruct{}

	setted, err := p.set(t, key.(string), set, nil)
	c.Assert(err, IsNil)
	c.Assert(setted, Equals, true)

	v, err := p.get(t, key.(string), nil)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, get)
}

func (s *CandySuite) TestProxy_Enumerate(c *C) {
	keys, err := p.enumerate(&MyStruct{Int: 42})
	c.Assert(err, IsNil)
	c.Assert(keys, DeepEquals, []string{
		"bool", "int", "int8", "int16", "int32", "int64", "uInt", "uInt8",
		"uInt16", "uInt32", "uInt64", "string", "bytes", "float32", "float64",
		"empty", "nested", "slice", "multiply",
	})
}

func (s *CandySuite) TestProxy_SetOnFunction(c *C) {
	setted, err := p.set(&MyStruct{Int: 21}, "multiply", 42.0, nil)
	c.Assert(err, IsNil)
	c.Assert(setted, Equals, false)
}

func (s *CandySuite) TestProxy_Properties(c *C) {
	provider := [][]interface{}{
		{&MyStruct{Int: 32}, "int", 32},
		{MyStruct{Int: 42}, "int", 42},
		{map[string]int{"foo": 21}, "foo", 21},
		{&map[string]int{"foo": 42}, "foo", 42},
	}

	for _, v := range provider {
		s.testProxyProperties(c, v[0], v[1], v[2])
	}
}

func (s *CandySuite) testProxyProperties(c *C, value, key, expected interface{}) {
	val, err := p.get(value, key.(string), nil)
	c.Assert(err, IsNil)
	c.Assert(val, Equals, expected)
}

func (s *CandySuite) TestProxy_Functions(c *C) {
	provider := [][]interface{}{
		{&MyStruct{}, "string"},
		{&customMap{}, "functionWithPtr"},
		{customMap{}, "functionWithoutPtr"},
		{customInt(1), "functionWithoutPtr"},
	}

	for _, v := range provider {
		s.testProxyFunction(c, v[0], v[1])
	}
}

func (s *CandySuite) testProxyFunction(c *C, value, key interface{}) {
	val, err := p.get(value, key.(string), nil)
	c.Assert(err, IsNil)
	c.Assert(val, NotNil)
}

func (s *CandySuite) TestProxyInternalKeys(c *C) {
	s.ctx.PushGlobalObject()
	s.ctx.PushObject()
	s.ctx.PushProxy(&MyStruct{Int: 142})
	s.ctx.PutPropString(-2, "obj")
	s.ctx.PutPropString(-2, "foo")
	s.ctx.Pop()

	//calls valueOf
	err := s.ctx.PevalString(`store(1 == foo.obj)`)
	c.Assert(err, IsNil)
	c.Assert(s.stored, Equals, false)

	//calls valueOf also toString
	err = s.ctx.PevalString(`store("[candyjs Proxy]" == foo.obj)`)
	c.Assert(err, IsNil)
	c.Assert(s.stored, Equals, true)

	err = s.ctx.PevalString(`foo.obj`)
	c.Assert(err, IsNil)

	//calls toJson
	js := s.ctx.JsonEncode(-1)
	r := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(js), &r)
	c.Assert(r["int"], Equals, 142.0)
}

type customInt int

func (c customInt) FunctionWithoutPtr() {}

type customMap map[string]int

func (c customMap) FunctionWithoutPtr() {}
func (c *customMap) FunctionWithPtr()   {}
