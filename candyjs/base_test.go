package candyjs

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type CandySuite struct {
	ctx    *Context
	stored interface{}
}

var _ = Suite(&CandySuite{})

func (s *CandySuite) SetUpTest(c *C) {
	s.ctx = NewContext()
	s.stored = nil
	s.ctx.PushGlobalGoFunction("store", func(value interface{}) {
		s.stored = value
	})
}

func (s *CandySuite) TestPushGlobalCandyJSObject(c *C) {
	c.Assert(s.ctx.PevalString(`store(CandyJS._functions.toString())`), IsNil)
	c.Assert(s.stored, Equals, "[object Object]")

	c.Assert(s.ctx.PevalString(`store(CandyJS._call.toString())`), IsNil)
	c.Assert(s.stored, Equals, "function anon() {/* ecmascript */}")

	c.Assert(s.ctx.PevalString(`store(CandyJS.proxy.toString())`), IsNil)
	c.Assert(s.stored, Equals, "function anon() {/* ecmascript */}")

	c.Assert(s.ctx.PevalString(`store(CandyJS.require.toString())`), IsNil)
	c.Assert(s.stored, Equals, "function anon() {/* native */}")
}

func (s *CandySuite) TestPushGlobalCandyJSObject_Require(c *C) {
	fn := func(ctx *Context) {
		ctx.PushString("qux")
	}

	RegisterPackagePusher("foo", fn)
	c.Assert(s.ctx.PevalString(`store(CandyJS.require("foo"))`), IsNil)
	c.Assert(s.stored, Equals, "qux")
}

func (s *CandySuite) TestSetRequireFunction(c *C) {
	s.ctx.SetRequireFunction(func(id string, a ...interface{}) string {
		return fmt.Sprintf(`exports.store = function () { store("%s"); };`, id)
	})

	c.Assert(s.ctx.PevalString("require('foo').store()"), IsNil)
	c.Assert(s.stored, Equals, "foo")
}

func (s *CandySuite) TestPushType(c *C) {
	s.ctx.PushGlobalObject()
	s.ctx.PushObject()
	s.ctx.PushType(MyStruct{})
	s.ctx.PutPropString(-2, "MyStruct")
	s.ctx.PutPropString(-2, "foo")
	s.ctx.Pop()

	c.Assert(s.ctx.PevalString(`
		obj = new foo.MyStruct()
		obj.int = 42
		store(obj)
	`), IsNil)

	c.Assert(s.stored.(*MyStruct).Int, Equals, 42)
}

func (s *CandySuite) TestGlobalPushType(c *C) {
	s.ctx.PushGlobalType("MyStruct", MyStruct{})

	c.Assert(s.ctx.PevalString(`
		obj = new MyStruct()
		obj.int = 42
		store(obj)
	`), IsNil)

	c.Assert(s.stored.(*MyStruct).Int, Equals, 42)
}

func (s *CandySuite) TestPushProxy(c *C) {
	s.ctx.PushGlobalObject()
	s.ctx.PushObject()
	s.ctx.PushProxy(&MyStruct{Int: 142})
	s.ctx.PutPropString(-2, "obj")
	s.ctx.PutPropString(-2, "foo")
	s.ctx.Pop()

	err := s.ctx.PevalString(`store(foo.obj.int)`)
	c.Assert(err, IsNil)
	c.Assert(s.stored, Equals, 142.0)
}

func (s *CandySuite) TestPushGlobalProxy_GetMap(c *C) {
	s.ctx.PushGlobalProxy("test", &map[string]int{"foo": 42})

	s.ctx.PevalString(`store(test.foo)`)
	c.Assert(s.stored, Equals, 42.0)
}

func (s *CandySuite) TestPushGlobalProxy_GetPtr(c *C) {
	s.ctx.PushGlobalProxy("test", &MyStruct{Int: 42})

	s.ctx.PevalString(`store(test.int)`)
	c.Assert(s.stored, Equals, 42.0)

	s.ctx.PevalString(`try { x = test.baz; } catch(err) { store(true); }`)
	c.Assert(s.stored, Equals, true)
}

func (s *CandySuite) TestPushGlobalProxy_Set(c *C) {
	s.ctx.PushGlobalProxy("test", &MyStruct{Int: 42})

	s.ctx.PevalString(`test.int = 21; store(test.int)`)
	c.Assert(s.stored, Equals, 21.0)

	s.ctx.PevalString(`try { test.baz = 21; } catch(err) { store(true); }`)
	c.Assert(s.stored, Equals, true)
}

func (s *CandySuite) TestPushGlobalProxy_Has(c *C) {
	s.ctx.PushGlobalProxy("test", &MyStruct{})
	s.ctx.PevalString(`store("int" in test)`)
	c.Assert(s.stored, Equals, true)

	s.ctx.PevalString(`store("qux" in test)`)
	c.Assert(s.stored, Equals, false)
}

func (s *CandySuite) TestPushGlobalProxy_Nested(c *C) {
	s.ctx.PushGlobalProxy("test", &MyStruct{
		Int:     42,
		Float64: 21.0,
		Nested:  &MyStruct{Int: 21},
	})

	c.Assert(s.ctx.PevalString(`store([
		test.int,
	    test.multiply(2),
	    test.nested.int,
	    test.nested.multiply(3)
	])`), IsNil)

	c.Assert(s.stored, DeepEquals, []interface{}{42.0, 84.0, 21.0, 63.0})
}

func (s *CandySuite) TestPushGlobalProxy_Integration(c *C) {
	now := time.Now()
	after := now.Add(time.Millisecond)

	s.ctx.PushGlobalProxy("a", now)
	s.ctx.PushGlobalProxy("b", after)

	s.ctx.PevalString(`store(b.sub(a))`)
	c.Assert(s.stored, Equals, 1000000.0)
}

func (s *CandySuite) TestPushGlobalInterface(c *C) {
	s.ctx.PushGlobalInterface("int", 42)

	c.Assert(s.ctx.PevalString(`store(int)`), IsNil)
	c.Assert(s.stored, Equals, 42.0)
}

func (s *CandySuite) TestPushGlobalStruct(c *C) {
	s.ctx.PushGlobalStruct("test", &MyStruct{
		Int:     42,
		Float64: 21.0,
		Nested:  &MyStruct{Int: 21},
	})

	c.Assert(s.ctx.PevalString(`store([
		test.int,
		test.multiply(2),
		test.nested.int,
		test.nested.multiply(3)
	])`), IsNil)

	c.Assert(s.stored, DeepEquals, []interface{}{42.0, 84.0, 21.0, 63.0})
}

func (s *CandySuite) TestPushGlobalValueInt(c *C) {
	s.ctx.pushGlobalValue("test", reflect.ValueOf(42))
	c.Assert(s.ctx.PevalString(`store(test)`), IsNil)
	c.Assert(s.stored, Equals, 42.0)
}

func (s *CandySuite) TestPushGlobalValueUint(c *C) {
	s.ctx.pushGlobalValue("test", reflect.ValueOf(uint(42)))
	c.Assert(s.ctx.PevalString(`store(test)`), IsNil)
	c.Assert(s.stored, Equals, 42.0)
}

func (s *CandySuite) TestPushGlobalValueFloat(c *C) {
	s.ctx.pushGlobalValue("test", reflect.ValueOf(42.2))
	c.Assert(s.ctx.PevalString(`store(test)`), IsNil)
	c.Assert(s.stored, Equals, 42.2)
}

func (s *CandySuite) TestPushGlobalValueString(c *C) {
	s.ctx.pushGlobalValue("test", reflect.ValueOf("foo"))
	c.Assert(s.ctx.PevalString(`store(test)`), IsNil)
	c.Assert(s.stored, Equals, "foo")
}

func (s *CandySuite) TestPushGlobalValueStruct(c *C) {
	s.ctx.pushGlobalValue("test", reflect.ValueOf(MyStruct{Int: 42}))
	c.Assert(s.ctx.PevalString(`store(test.int)`), IsNil)
	c.Assert(s.stored, Equals, 42.0)
}

func (s *CandySuite) TestPushGlobalValueStructPtr(c *C) {
	s.ctx.pushGlobalValue("test", reflect.ValueOf(&MyStruct{Int: 42}))
	c.Assert(s.ctx.PevalString(`store(test.int)`), IsNil)
	c.Assert(s.stored, Equals, 42.0)
}

func (s *CandySuite) TestPushGlobalValueNil(c *C) {
	s.ctx.pushGlobalValue("test", reflect.ValueOf(nil))
	c.Assert(s.ctx.PevalString(`store(test)`), IsNil)
	c.Assert(s.stored, Equals, nil)
}

func (s *CandySuite) TestPushGlobalValueDefault(c *C) {
	s.ctx.pushGlobalValue("test", reflect.ValueOf([]string{"foo", "bar"}))
	c.Assert(s.ctx.PevalString(`store(test)`), IsNil)
	c.Assert(s.stored, DeepEquals, []interface{}{"foo", "bar"})
}

func (s *CandySuite) TestPushGlobalValueStringPtr(c *C) {
	foo := "foo"
	s.ctx.pushGlobalValue("test", reflect.ValueOf(&foo))
	c.Assert(s.ctx.PevalString(`store(test)`), IsNil)
	c.Assert(s.stored, Equals, "foo")
}

func (s *CandySuite) PendingTestPushGlobalValueWithMethods(c *C) {
	s.ctx.pushGlobalValue("test", reflect.ValueOf(time.Duration(1e5)))
	c.Assert(s.ctx.PevalString(`store(test.string())`), IsNil)
	c.Assert(s.stored, Equals, 42.0)
}

func (s *CandySuite) TestPushGlobalValues(c *C) {
	s.ctx.pushGlobalValues("test", []reflect.Value{
		reflect.ValueOf("foo"), reflect.ValueOf("qux"),
	})

	c.Assert(s.ctx.PevalString(`store(test)`), IsNil)
	c.Assert(s.stored, DeepEquals, []interface{}{"foo", "qux"})
}

func (s *CandySuite) TestPushGlobalGoFunction_String(c *C) {
	var called interface{}
	s.ctx.PushGlobalGoFunction("test_in_string", func(s string) {
		called = s
	})

	s.ctx.EvalString("test_in_string('foo')")
	c.Assert(called, Equals, "foo")
}

func (s *CandySuite) TestPushGlobalGoFunction_Int(c *C) {
	var ri, ri8, ri16, ri32, ri64 interface{}
	s.ctx.PushGlobalGoFunction("test_in_int", func(i int, i8 int8, i16 int16, i32 int32, i64 int64) {
		ri = i
		ri8 = i8
		ri16 = i16
		ri32 = i32
		ri64 = i64
	})

	s.ctx.EvalString("test_in_int(42, 8, 16, 32, 64)")
	c.Assert(ri, Equals, 42)
	c.Assert(ri8, Equals, int8(8))
	c.Assert(ri16, Equals, int16(16))
	c.Assert(ri32, Equals, int32(32))
	c.Assert(ri64, Equals, int64(64))
}

func (s *CandySuite) TestPushGlobalGoFunction_Uint(c *C) {
	var ri, ri8, ri16, ri32, ri64 interface{}
	s.ctx.PushGlobalGoFunction("test_in_uint", func(i uint, i8 uint8, i16 uint16, i32 uint32, i64 uint64) {
		ri = i
		ri8 = i8
		ri16 = i16
		ri32 = i32
		ri64 = i64
	})

	s.ctx.EvalString("test_in_uint(42, 8, 16, 32, 64)")
	c.Assert(ri, Equals, uint(42))
	c.Assert(ri8, Equals, uint8(8))
	c.Assert(ri16, Equals, uint16(16))
	c.Assert(ri32, Equals, uint32(32))
	c.Assert(ri64, Equals, uint64(64))
}

func (s *CandySuite) TestPushGlobalGoFunction_Float(c *C) {
	var called64 interface{}
	var called32 interface{}
	s.ctx.PushGlobalGoFunction("test_in_float", func(f64 float64, f32 float32) {
		called64 = f64
		called32 = f32
	})

	s.ctx.EvalString("test_in_float(42, 42)")
	c.Assert(called64, Equals, 42.0)
	c.Assert(called32, Equals, float32(42.0))
}

func (s *CandySuite) TestPushGlobalGoFunction_Bool(c *C) {
	var called interface{}
	s.ctx.PushGlobalGoFunction("test_in_bool", func(b bool) {
		called = b
	})

	s.ctx.EvalString("test_in_bool(true)")
	c.Assert(called, Equals, true)
}

func (s *CandySuite) TestPushGlobalGoFunction_Interface(c *C) {
	var called interface{}
	s.ctx.PushGlobalGoFunction("test", func(i interface{}) {
		called = i
	})

	s.ctx.EvalString("test('qux')")
	c.Assert(called, Equals, "qux")
}

func (s *CandySuite) TestPushGlobalGoFunction_Struct(c *C) {
	var called *MyStruct
	s.ctx.PushGlobalGoFunction("test", func(m *MyStruct) {
		called = m
	})

	s.ctx.EvalString("test({'int':42})")
	c.Assert(called.Int, Equals, 42)
}

func (s *CandySuite) TestPushGlobalGoFunction_Slice(c *C) {
	var called interface{}
	s.ctx.PushGlobalGoFunction("test_in_slice", func(s []interface{}) {
		called = s
	})

	s.ctx.EvalString("test_in_slice(['foo', 42])")
	c.Assert(called, DeepEquals, []interface{}{"foo", 42.0})
}

func (s *CandySuite) TestPushGlobalGoFunction_Map(c *C) {
	var called interface{}
	s.ctx.PushGlobalGoFunction("test_in_map", func(s map[string]interface{}) {
		called = s
	})

	s.ctx.EvalString("test_in_map({foo: 42, qux: {bar: 'bar'}})")

	c.Assert(called, DeepEquals, map[string]interface{}{
		"foo": 42.0,
		"qux": map[string]interface{}{"bar": "bar"},
	})
}

func (s *CandySuite) TestPushGlobalGoFunction_Nil(c *C) {
	var cm, cs, ci, cst interface{}
	s.ctx.PushGlobalGoFunction("test_nil", func(m map[string]interface{}, s []interface{}, i int, st string) {
		cm = m
		cs = s
		ci = i
		cst = st
	})

	s.ctx.EvalString("test_nil(null, null, null, null)")
	c.Assert(cm, DeepEquals, map[string]interface{}(nil))
	c.Assert(cs, DeepEquals, []interface{}(nil))
	c.Assert(ci, DeepEquals, 0)
	c.Assert(cst, DeepEquals, "")
}

func (s *CandySuite) TestPushGlobalGoFunction_Optional(c *C) {
	var cm, cs, ci, cst interface{}
	s.ctx.PushGlobalGoFunction("test_optional", func(m map[string]interface{}, s []interface{}, i int, st string) {
		cm = m
		cs = s
		ci = i
		cst = st
	})

	s.ctx.EvalString("test_optional()")
	c.Assert(cm, DeepEquals, map[string]interface{}(nil))
	c.Assert(cs, DeepEquals, []interface{}(nil))
	c.Assert(ci, DeepEquals, 0)
	c.Assert(cst, DeepEquals, "")
}

func (s *CandySuite) TestPushGlobalGoFunction_Variadic(c *C) {
	var calledA interface{}
	var calledB interface{}
	s.ctx.PushGlobalGoFunction("test_in_variadic", func(s string, is ...int) {
		calledA = s
		calledB = is
	})

	s.ctx.EvalString("test_in_variadic('foo', 21, 42)")
	c.Assert(calledA, DeepEquals, "foo")
	c.Assert(calledB, DeepEquals, []int{21, 42})
}

func (s *CandySuite) TestPushGlobalGoFunction_EmptyVariadic(c *C) {
	var calledA interface{}
	var calledB interface{}
	s.ctx.PushGlobalGoFunction("test_in_variadic", func(s string, is ...int) {
		calledA = s
		calledB = is
	})

	s.ctx.EvalString("test_in_variadic('foo')")
	c.Assert(calledA, DeepEquals, "foo")
	c.Assert(calledB, DeepEquals, []int{})
}

func (s *CandySuite) TestPushGlobalGoFunction_ReturnMultiple(c *C) {
	s.ctx.PushGlobalGoFunction("test", func() (int, int, error) {
		return 2, 4, nil
	})

	c.Assert(s.ctx.PevalString("store(test())"), IsNil)
	c.Assert(s.stored, HasLen, 2)
	c.Assert(s.stored.([]interface{})[0], Equals, 2.0)
	c.Assert(s.stored.([]interface{})[1], Equals, 4.0)
}

func (s *CandySuite) TestPushGlobalGoFunction_ReturnStruct(c *C) {
	s.ctx.PushGlobalGoFunction("test", func() *MyStruct {
		return &MyStruct{Int: 42}
	})

	c.Assert(s.ctx.PevalString("store(test().multiply(3))"), IsNil)
	c.Assert(s.stored, Equals, 126.0)
}

func (s *CandySuite) TestPushGlobalGoFunction_Function(c *C) {
	s.ctx.PushGlobalGoFunction("test", func(fn func(int, int) int) {
		s.stored = fn
	})

	c.Assert(s.ctx.PevalString(`
		test(CandyJS.proxy(function(a, b) { return a * b; }));
	`), IsNil)

	c.Assert(s.stored.(func(int, int) int)(10, 5), Equals, 50)
}

func (s *CandySuite) TestPushGlobalGoFunction_FunctionMultiple(c *C) {
	s.ctx.PushGlobalGoFunction("test", func(fn func(int, int) (int, int)) {
		s.stored = fn
	})

	c.Assert(s.ctx.PevalString(`
		test(CandyJS.proxy(function(a, b) { return [b, a]; }));
	`), IsNil)

	a, b := s.stored.(func(int, int) (int, int))(10, 5)
	c.Assert(a, Equals, 5)
	c.Assert(b, Equals, 10)
}

func (s *CandySuite) TestPushGlobalGoFunction_Error(c *C) {
	s.ctx.PushGlobalGoFunction("test", func() (string, error) {
		return "foo", fmt.Errorf("foo")
	})

	c.Assert(s.ctx.PevalString(`
		try {
			test();
		} catch(err) {
			store(true);
		}
	`), IsNil)
	c.Assert(s.stored, Equals, true)
}

func (s *CandySuite) TearDownTest(c *C) {
	s.ctx.DestroyHeap()
}

type MyStruct struct {
	Bool    bool
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	UInt    uint
	UInt8   uint8
	UInt16  uint16
	UInt32  uint32
	UInt64  uint64
	String  string
	Bytes   []byte
	Float32 float32
	Float64 float64
	Empty   *MyStruct
	Nested  *MyStruct
	Slice   []int
	private int
}

func (m *MyStruct) Multiply(x int) int {
	return m.Int * x
}

func (m *MyStruct) privateMethod() int {
	return 1
}
