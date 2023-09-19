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
	"slices"
	"sort"
)

// Slices utilities

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func SameList(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	slices.Sort(s1)
	slices.Sort(s2)
	for i, v := range s1 {
		if s2[i] != v {
			return false
		}
	}
	return true
}

// Mapping utilities

func GetValues(vm map[string]string) []string {
	var keys []string
	for k := range vm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func HasSameKeys(vm1, vm2 map[string]string) bool {
	k1 := GetValues(vm1)
	k2 := GetValues(vm2)
	return SameList(k1, k2)
}

// String utilities

func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
