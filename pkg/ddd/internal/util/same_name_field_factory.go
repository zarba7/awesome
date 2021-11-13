package util

import (
	"log"
	"reflect"
	"unsafe"
)

type SameNameStructFactory struct {
	slice []*sameNameField
}

func (the SameNameStructFactory) RegisterStruct(tid uint32, elem interface{}) {
	if len(the.slice) <= int(tid) { //SameNameStructFactory todo 优化成数组
		old := the.slice
		the.slice = make([]*sameNameField, tid+1)
		copy(the.slice, old)
	}
	if the.slice[tid] != nil {
		log.Panicf("{%d} same elem %v", tid, elem)
	}
	the.slice[tid] = newSameNameField(elem)
}

func (the *SameNameStructFactory) FindStructByFieldVal(tid uint32, fVal interface{}) interface{} {
	return the.slice[tid].newStructByFieldVal(fVal)
}

type Unmarshal func(args interface{}) error

func (the *SameNameStructFactory) FindStructByUnmarshal(tid uint32, unmarshal Unmarshal) (interface{}, error) {
	fty := the.slice[tid]
	fVal := fty.newFieldVal()
	return fty.newStructByFieldVal(fVal), unmarshal(fVal)
}

type sameNameField struct {
	structType  reflect.Type
	fieldType   reflect.Type
	fieldOffset uintptr
}

func (the *sameNameField) newFieldVal() interface{} { return reflect.New(the.fieldType).Interface() }

func (the *sameNameField) newStructByFieldVal(val interface{}) interface{} {
	V := reflect.New(the.structType)
	fPtr := V.Elem().UnsafeAddr() + the.fieldOffset
	*(*uintptr)(unsafe.Pointer(fPtr)) = reflect.ValueOf(val).Elem().UnsafeAddr()
	return V.Interface()
}

func newSameNameField(elem interface{}) *sameNameField {
	the := &sameNameField{}
	the.structType = reflect.TypeOf(elem)
	if the.structType.Kind() == reflect.Ptr {
		the.structType = the.structType.Elem()
	}
	if the.structType.Kind() != reflect.Struct {
		log.Panicf("{%s.%s} is no struct", the.structType.PkgPath(), the.structType.Name())
	}
	for i := 0; i < the.structType.NumField(); i++ {
		f := the.structType.Field(i)
		if f.Type.Kind() == reflect.Ptr && the.structType.Name() == f.Name {
			the.fieldType = f.Type.Elem()
			the.fieldOffset = f.Offset
			return the
		}
	}
	log.Panicf("{%s.%s} have no same name field", the.structType.PkgPath(), the.structType.Name())
	return nil
}
