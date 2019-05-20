// Copyright (c) 2016 ~ 2019, Alex Stocks.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hessian

import (
	"reflect"
	"unsafe"
)

import (
	jerrors "github.com/juju/errors"
)

// used to ref object,list,map
type _refElem struct {
	// record the kind of target, objects are the same only if the address and kind are the same
	kind reflect.Kind

	// ref index
	index int
}

// _refHolder is used to record decode list, the address of which may change when appending more element.
type _refHolder struct {
	// destinations
	destinations []reflect.Value

	value reflect.Value
}

var _refHolderType = reflect.TypeOf(_refHolder{})

// change ref value
func (h *_refHolder) change(v reflect.Value) {
	if h.value.CanAddr() && v.CanAddr() && h.value.Pointer() == v.Pointer() {
		return
	}
	h.value = v
}

// notice all destinations ref to the value
func (h *_refHolder) notify() {
	for _, dest := range h.destinations {
		SetValue(dest, h.value)
	}
}

// add destination
func (h *_refHolder) add(dest reflect.Value) {
	h.destinations = append(h.destinations, dest)
}

// Add reference
func (d *Decoder) appendRefs(v interface{}) *_refHolder {
	var holder *_refHolder
	vv := EnsurePackValue(v)
	// only slice and array need ref holder , for its address changes when decoding
	if vv.Kind() == reflect.Slice || vv.Kind() == reflect.Array {
		holder = &_refHolder{
			value: vv,
		}
		// pack holder value
		v = reflect.ValueOf(holder)
	}

	d.refs = append(d.refs, v)
	return holder
}

//encRef encode ref index
func encRef(b []byte, index int) []byte {
	return encInt32(append(b, BC_REF), int32(index))
}

// return the order number of ref object if found ,
// otherwise, add the object into the encode ref map
func (e *Encoder) checkRefMap(v reflect.Value) (int, bool) {
	var (
		kind reflect.Kind
		addr unsafe.Pointer
	)

	if v.Kind() == reflect.Ptr {
		for v.Elem().Kind() == reflect.Ptr {
			v = v.Elem()
		}
		kind = v.Elem().Kind()
		if kind == reflect.Slice || kind == reflect.Map {
			addr = unsafe.Pointer(v.Elem().Pointer())
		} else {
			addr = unsafe.Pointer(v.Pointer())
		}
	} else {
		kind = v.Kind()
		switch kind {
		case reflect.Slice, reflect.Map:
			addr = unsafe.Pointer(v.Pointer())
		default:
			addr = unsafe.Pointer(PackPtr(v).Pointer())
		}
	}

	if elem, ok := e.refMap[addr]; ok {
		// the array addr is equal to the first elem, which must ignore
		if elem.kind == kind {
			return elem.index, ok
		}
		return 0, false
	}

	n := len(e.refMap)
	e.refMap[addr] = _refElem{kind, n}
	return 0, false
}

/////////////////////////////////////////
// Ref
/////////////////////////////////////////

// # value reference (e.g. circular trees and graphs)
// ref        ::= x51 int            # reference to nth map/list/object
func (d *Decoder) decRef(flag int32) (interface{}, error) {
	var (
		err error
		tag byte
		i   int32
	)

	if flag != TAG_READ {
		tag = byte(flag)
	} else {
		tag, _ = d.readByte()
	}

	switch {
	case tag == BC_REF:
		i, err = d.decInt32(TAG_READ)
		if err != nil {
			return nil, err
		}

		if len(d.refs) <= int(i) {
			return nil, ErrIllegalRefIndex
		}
		// return the exact ref object, which maybe a _refHolder
		return d.refs[i], nil

	default:
		return nil, jerrors.Errorf("decRef illegal ref type tag:%+v", tag)
	}
}
