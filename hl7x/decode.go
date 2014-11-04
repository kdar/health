package hl7x

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kdar/health/hl7"
)

func Unmarshal(src hl7.Data, dst interface{}) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("interface must be a pointer to struct")
	}
	v = v.Elem()

	d := newDecoder()
	d.decode(src, reflect.ValueOf(dst).Elem())
	if d.err.Errors != nil {
		return d.err
	}
	return nil
}

type decoder struct {
	err *Error
}

func newDecoder() *decoder {
	return &decoder{
		err: &Error{},
	}
}

func (d *decoder) decode(src hl7.Data, dst reflect.Value) {
	dstKind := dst.Kind()
	switch dstKind {
	case reflect.String:
		d.decodeString(src, dst)
	case reflect.Struct:
		d.decodeStruct(src, dst)
	case reflect.Slice:
		d.decodeSlice(src, dst)
	default:
		d.err.append(fmt.Errorf("unsupported type: %s", dstKind))
		return
	}
}

func (d *decoder) decodeString(src hl7.Data, dst reflect.Value) {
	v, ok := src.(hl7.Field)
	if !ok {
		d.err.append(fmt.Errorf("decodeString: src is not a hl7.Field, it's: %s", reflect.TypeOf(src)))
		return
	}

	dst.SetString(string(v))
}

func (d *decoder) decodeStruct(src hl7.Data, dst reflect.Value) {
	typ := dst.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		fieldValue := dst.FieldByName(fieldName)

		// So just skip it.
		if !fieldValue.CanSet() {
			d.err.append(fmt.Errorf("%s.%s is not settable", typ.Name(), fieldName))
			continue
		}

		index := i
		if _, ok := src.(hl7.Segment); ok {
			index += 1
		}

		newSrc, ok := src.Index(index)
		if !ok {
			//d.err.append(fmt.Errorf("could not find src index for %s.%s", typ.Name(), fieldName))
			continue
		}

		d.decode(newSrc, fieldValue)
	}
}

func (d *decoder) decodeSlice(src hl7.Data, dst reflect.Value) {
	dstType := dst.Type()
	dstElemType := dstType.Elem()
	sliceType := reflect.SliceOf(dstElemType)

	srcLen := src.Len()
	dstSlice := reflect.MakeSlice(sliceType, srcLen, srcLen)
	for i := 0; i < srcLen; i++ {
		currentField := dstSlice.Index(i)
		if newSrc, ok := src.Index(i); ok {
			d.decode(newSrc, currentField)
		}
	}

	dst.Set(dstSlice)
}
