// Copyright 2015 The etcd Authors
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

package types

// ID represents a generic identifier which is canonically
// stored as a uint64 but is typically represented as a
// base-16 string for input/output
type ID string

func (i ID) String() string {
	return string(i)
}

// IDFromString attempts to create an ID from a base-16 string.
func IDFromString(s string) (ID, error) {
	return ID(s), nil
}

// IDSlice implements the sort interface
type IDSlice []ID

func (p IDSlice) Len() int           { return len(p) }
func (p IDSlice) Less(i, j int) bool { return string(p[i]) < string(p[j]) }
func (p IDSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
