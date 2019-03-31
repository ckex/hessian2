// Copyright (c) 2016 ~ 2019, dubbogo.
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
	"testing"
)
import (
	"github.com/stretchr/testify/assert"
)

func doTestReflectResponse(t *testing.T, in interface{}, out interface{}) {
	err := ReflectResponse(in, out)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	result := UnpackPtrValue(reflect.ValueOf(out)).Interface()

	equal := reflect.DeepEqual(in, result)
	if !equal {
		t.Errorf("expect [%v]: %v, but got [%v]: %v", reflect.TypeOf(in), in, reflect.TypeOf(result), result)
	}
}

func TestReflectResponse(t *testing.T) {
	var b bool
	doTestReflectResponse(t, true, &b)
	doTestReflectResponse(t, false, &b)

	var i int
	doTestReflectResponse(t, 123, &i)
	doTestReflectResponse(t, 234, &i)

	var i16 int16
	doTestReflectResponse(t, int16(456), &i16)

	var i64 int64
	doTestReflectResponse(t, int64(789), &i64)

	var s string
	doTestReflectResponse(t, "hello world", &s)

	type rr struct {
		Name string
		Num  int
	}

	var r1 rr
	doTestReflectResponse(t, rr{"dubbogo", 32}, &r1)

	// ------ map test -------
	// NOTE: map type currently MUST be map[interface{}]interface{}
	m1 := make(map[interface{}]interface{})
	var m1r map[interface{}]interface{}
	m1["hello"] = "world"
	m1[1] = "go"
	m1["dubbo"] = 666
	doTestReflectResponse(t, m1, &m1r)

	m2 := make(map[string]string)
	var m2r map[string]string
	m2["hello"] = "world"
	m2["dubbo"] = "666"
	doTestReflectResponse(t, m2, &m2r)

	m3 := make(map[string]rr)
	var m3r map[string]rr
	m3["dubbo"] = rr{"hello", 123}
	m3["go"] = rr{"world", 456}
	doTestReflectResponse(t, m3, &m3r)

	// ------ slice test -------
	s1 := []string{"abc", "def", "hello", "world"}
	var s1r []string
	doTestReflectResponse(t, s1, &s1r)

	s2 := []rr{rr{"dubbo", 666}, rr{"go", 999}}
	var s2r []rr
	doTestReflectResponse(t, s2, &s2r)

	s3 := []interface{}{rr{"dubbo", 666}, 123, "hello"}
	var s3r []interface{}
	doTestReflectResponse(t, s3, &s3r)
}

// separately test copy map to map[interface{}]interface{}
func TestCopyMap(t *testing.T) {
	type rr struct {
		Name string
		Num  int
	}

	m3 := make(map[string]rr)
	var m3r map[interface{}]interface{}
	r1 := rr{"hello", 123}
	r2 := rr{"world", 456}
	m3["dubbo"] = r1
	m3["go"] = r2

	err := ReflectResponse(m3, &m3r)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, 2, len(m3r))

	rr1, ok := m3r["dubbo"]
	assert.True(t, ok)
	assert.True(t, reflect.DeepEqual(r1, rr1))

	rr2, ok := m3r["go"]
	assert.True(t, ok)
	assert.True(t, reflect.DeepEqual(r2, rr2))
}

func TestCopySlice(t *testing.T) {
	type rr struct {
		Name string
		Num  int
	}

	r1 := rr{"hello", 123}
	r2 := rr{"world", 456}

	s1 := []rr{r1, r2}
	var s1r []interface{}

	err := ReflectResponse(s1, &s1r)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, 2, len(s1r))
	assert.True(t, reflect.DeepEqual(r1, s1r[0]))
	assert.True(t, reflect.DeepEqual(r2, s1r[1]))
}
