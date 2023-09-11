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
