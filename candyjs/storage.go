package candyjs

import "C"
import (
	"sync"
	"unsafe"
)

type storage struct {
	vars map[unsafe.Pointer]interface{}
	sync.Mutex
}

func newStorage() *storage {
	return &storage{
		vars: make(map[unsafe.Pointer]interface{}, 0),
	}
}

func (s *storage) add(v interface{}) unsafe.Pointer {
	s.Lock()
	defer s.Unlock()

	ptr := C.malloc(1)
	s.vars[ptr] = v

	return ptr
}

func (s *storage) get(ptr unsafe.Pointer) interface{} {
	s.Lock()
	defer s.Unlock()

	return s.vars[ptr]
}
