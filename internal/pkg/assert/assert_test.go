// Copyright © 2018 Bitgrip <berlin@bitgrip.de>
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

package assert

import (
	"testing"
)

func TestTextByteSlice(t *testing.T) {
	expected := "test text"
	slice1 := []byte(expected)
	Equals(t, expected, buildCompareString(slice1))
}

func TestBinaryByteSlice(t *testing.T) {
	slice1 := []byte{0xF3, 0xF4, 0xF5}
	expected := string(slice1)
	Equals(t, expected, buildCompareString(slice1))
}
