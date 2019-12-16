// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.


package ctypes

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	type thing struct {
		num uint64
	}
	thing1 := thing{42}

	if IsEmpty(thing1) {
		t.Error("empty not empty")
	}

	thing2 := thing{}
	if !IsEmpty(thing2) {
		t.Error("not empty empty")
	}

	if !IsEmpty(nil) {
		t.Error("nil not empty")
	}
}
