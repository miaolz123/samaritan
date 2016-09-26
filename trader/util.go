package trader

import (
	"github.com/robertkrimen/otto"
)

func parseObj(o *otto.Object) map[string]interface{} {
	obj := make(map[string]interface{})
	if o != nil {
		for _, k := range o.Keys() {
			v, _ := o.Get(k)
			if v.IsObject() {
				obj[k] = parseObj(v.Object())
			} else {
				obj[k] = v
			}
		}
	}
	return obj
}
