/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"reflect"
	"testing"
)

type TestSubObject struct {
	SubObject string
}

type TestObject struct {
	StringObject string
	IntObject    int
	ListObject   []string
	MapObject    map[string]string
	StructObject TestSubObject
}

func TestGetValues(t *testing.T) {
	result := GetValues(map[string]string{"a": "A"})
	if !reflect.DeepEqual(result, []string{"a"}) {
		t.Fail()
	}
	result = GetValues(map[string]string{})
	if len(result) != 0 {
		t.Fail()
	}
	result = GetValues(map[string]string{"a": "B", "b": "C"})
	if !reflect.DeepEqual(result, []string{"a", "b"}) {
		t.Fail()
	}
}

func TestHasSameKeys(t *testing.T) {
	result := HasSameKeys(map[string]string{"a": "A"}, map[string]string{"a": "A"})
	if !result {
		t.Fail()
	}
	result = HasSameKeys(map[string]string{}, map[string]string{"a": "A"})
	if result {
		t.Fail()
	}
	result = HasSameKeys(map[string]string{}, map[string]string{})
	if !result {
		t.Fail()
	}
	result = HasSameKeys(map[string]string{"a": "A"}, map[string]string{"b": "B"})
	if result {
		t.Fail()
	}
}

func TestContains(t *testing.T) {
	result := Contains([]string{"a", "b", "c"}, "a")
	if !result {
		t.Fail()
	}
	result = Contains([]string{"a", "b", "c"}, "c")
	if !result {
		t.Fail()
	}
	result = Contains([]string{"a", "b", "c"}, "d")
	if result {
		t.Fail()
	}
	result = Contains([]string{"a", "a"}, "a")
	if !result {
		t.Fail()
	}
	result = Contains([]string{}, "a")
	if result {
		t.Fail()
	}
}

func TestSameList(t *testing.T) {
	result := SameList([]string{"a", "b", "c"}, []string{"a", "b", "c"})
	if !result {
		t.Fail()
	}
	result = SameList([]string{"a", "b", "c"}, []string{"c", "b", "a"})
	if !result {
		t.Fail()
	}
	result = SameList([]string{"a", "b", "c"}, []string{"a", "a", "b", "c"})
	if result {
		t.Fail()
	}
	result = SameList([]string{}, []string{})
	if !result {
		t.Fail()
	}
	result = SameList([]string{"a", "a", "b"}, []string{"a", "b", "b"})
	if result {
		t.Fail()
	}
}

func TestRemoveString(t *testing.T) {
	result := RemoveString([]string{"a", "b", "c"}, "a")
	if !SameList(result, []string{"b", "c"}) {
		t.Fail()
	}
	result = RemoveString([]string{"a", "b", "a"}, "a")
	if !SameList(result, []string{"b"}) {
		t.Fail()
	}
	result = RemoveString([]string{"b", "b", "a"}, "a")
	if !SameList(result, []string{"b", "b"}) {
		t.Fail()
	}
	result = RemoveString([]string{"a", "b", "c"}, "d")
	if !SameList(result, []string{"a", "b", "c"}) {
		t.Fail()
	}
	result = RemoveString([]string{}, "a")
	if !SameList(result, []string{}) {
		t.Fail()
	}
}

func TestGetValueOf(t *testing.T) {
	obj := TestObject{
		"string", 1, []string{"list", "object"}, map[string]string{"map": "object"}, TestSubObject{"struct"},
	}
	result := GetValueOf(obj, "StringObject")
	if result.Kind() != reflect.String {
		t.Fail()
	}
	if result.String() != "string" {
		t.Fail()
	}
	result = GetValueOf(obj, "IntObject")
	if result.Kind() != reflect.Int {
		t.Fail()
	}
	if result.Int() != 1 {
		t.Fail()
	}
	result = GetValueOf(obj, "ListObject")
	if result.Kind() != reflect.Slice {
		t.Fail()
	}
	if !SameList(result.Interface().([]string), []string{"list", "object"}) {
		t.Fail()
	}
	result = GetValueOf(obj, "MapObject")
	if result.Kind() != reflect.Map {
		t.Fail()
	}
	m := result.Interface().(map[string]string)
	if m["map"] != "object" {
		t.Fail()
	}
	result = GetValueOf(obj, "StructObject")
	if result.Kind() != reflect.Struct {
		t.Fail()
	}
	s := result.Interface().(TestSubObject)
	if s.SubObject != "struct" {
		t.Fail()
	}
}
