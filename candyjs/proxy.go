package candyjs

import "C"
import (
	"errors"
	"reflect"
)

// ErrUndefinedProperty is throw when a property for a given proxied object on
// javascript cannot be found, basically a valid method or field cannot found.
var ErrUndefinedProperty = errors.New("undefined property")

var (
	p = &proxy{}

	//internalKeys map contains the keys that are called by duktape and cannot
	//throw an error, the value of the map is the value returned when this keys
	//are requested.
	internalKeys = map[string]interface{}{
		"toJSON": nil, "valueOf": nil,
		"toString": func() string { return "[candyjs Proxy]" },
	}
)

type proxy struct{}

func (p *proxy) has(t interface{}, k string) bool {
	_, err := p.getProperty(t, k)
	return err != ErrUndefinedProperty
}

func (p *proxy) get(t interface{}, k string, recv interface{}) (interface{}, error) {
	f, err := p.getProperty(t, k)
	if err != nil {
		if v, isInternal := internalKeys[k]; isInternal {
			return v, nil
		}

		return nil, err
	}

	return f.Interface(), nil
}

func (p *proxy) set(t interface{}, k string, v, recv interface{}) (bool, error) {
	f, err := p.getProperty(t, k)
	if err != nil {
		return false, err
	}

	if !f.CanSet() {
		return false, nil
	}

	value := reflect.Zero(f.Type())
	if v != nil {
		value = reflect.ValueOf(castNumberToGoType(f.Kind(), v))
	}

	f.Set(value)
	return true, nil
}

func (p *proxy) getProperty(t interface{}, key string) (reflect.Value, error) {
	v := reflect.ValueOf(t)
	r, found := p.getValueFromKind(key, v)
	if !found {
		return r, ErrUndefinedProperty
	}

	return r, nil
}

func (p *proxy) getValueFromKind(key string, v reflect.Value) (reflect.Value, bool) {
	var value reflect.Value
	var found bool
	switch v.Kind() {
	case reflect.Ptr:
		value, found = p.getValueFromKindPtr(key, v)
	case reflect.Struct:
		value, found = p.getValueFromKindStruct(key, v)
	case reflect.Map:
		value, found = p.getValueFromKindMap(key, v)
	}

	if !found {
		return p.getMethod(key, v)
	}

	return value, found
}

func (p *proxy) getValueFromKindPtr(key string, v reflect.Value) (reflect.Value, bool) {
	r, found := p.getMethod(key, v)
	if !found {
		return p.getValueFromKind(key, v.Elem())
	}

	return r, found
}

func (p *proxy) getValueFromKindStruct(key string, v reflect.Value) (reflect.Value, bool) {
	var r reflect.Value
	for _, name := range nameToGo(key) {
		r = v.FieldByName(name)
		if r.IsValid() {
			break
		}
	}

	return r, r.IsValid()
}

func (p *proxy) getValueFromKindMap(key string, v reflect.Value) (reflect.Value, bool) {
	keyValue := reflect.ValueOf(key)
	r := v.MapIndex(keyValue)

	return r, r.IsValid()
}

func (p *proxy) getMethod(key string, v reflect.Value) (reflect.Value, bool) {
	var r reflect.Value
	for _, name := range nameToGo(key) {
		r = v.MethodByName(name)
		if r.IsValid() {
			break
		}
	}

	return r, r.IsValid()
}

func (p *proxy) enumerate(t interface{}) (interface{}, error) {
	return p.getPropertyNames(t)
}

func (p *proxy) getPropertyNames(t interface{}) ([]string, error) {
	v := reflect.ValueOf(t)

	var names []string
	var err error
	switch v.Kind() {
	case reflect.Ptr:
		names, err = p.getPropertyNames(v.Elem().Interface())
		if err != nil {
			return nil, err
		}
	case reflect.Struct:
		cFields := v.NumField()
		for i := 0; i < cFields; i++ {
			fieldName := v.Type().Field(i).Name
			if !isExported(fieldName) {
				continue
			}

			names = append(names, nameToJavaScript(fieldName))
		}
	}

	mCount := v.NumMethod()
	for i := 0; i < mCount; i++ {
		methodName := v.Type().Method(i).Name
		if !isExported(methodName) {
			continue
		}

		names = append(names, nameToJavaScript(methodName))
	}

	return names, nil
}

func castNumberToGoType(k reflect.Kind, v interface{}) interface{} {
	switch k {
	case reflect.Int:
		v = int(v.(float64))
	case reflect.Int8:
		v = int8(v.(float64))
	case reflect.Int16:
		v = int16(v.(float64))
	case reflect.Int32:
		v = int32(v.(float64))
	case reflect.Int64:
		v = int64(v.(float64))
	case reflect.Uint:
		v = uint(v.(float64))
	case reflect.Uint8:
		v = uint8(v.(float64))
	case reflect.Uint16:
		v = uint16(v.(float64))
	case reflect.Uint32:
		v = uint32(v.(float64))
	case reflect.Uint64:
		v = uint64(v.(float64))
	case reflect.Float32:
		v = float32(v.(float64))
	}

	return v
}
