// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"testing"

	"github.com/jxo/lime/text"
)

func TestSingleSelection(t *testing.T) {
	tests := []findTest{
		{
			"Hello World!\nTest123123\nAbrakadabra\n",
			[]text.Region{{1, 1}, {2, 2}, {3, 3}, {6, 6}},
			[]text.Region{{1, 1}},
			false,
		},
		{
			"Hello World!\nTest123123\nAbrakadabra\n",
			[]text.Region{{2, 2}, {3, 3}, {6, 6}},
			[]text.Region{{2, 2}},
			false,
		},
		{
			"Hello World!\nTest123123\nAbrakadabra\n",
			[]text.Region{{5, 5}},
			[]text.Region{{5, 5}},
			false,
		},
	}

	runFindTest(tests, t, "single_selection")
}

func TestSelectAll(t *testing.T) {
	tests := []findTest{
		{
			"Hello World!\nTest123123\nAbrakadabra\n",
			[]text.Region{{1, 1}, {2, 2}, {3, 3}, {6, 6}},
			[]text.Region{{0, 36}},
			false,
		},
		{
			"Hello World!\nTest123123\nAbrakadabra\n",
			[]text.Region{{2, 2}, {3, 3}, {6, 6}},
			[]text.Region{{0, 36}},
			false,
		},
		{
			"Hello World!\nTest123123\nAbrakadabra\n",
			[]text.Region{{5, 5}},
			[]text.Region{{0, 36}},
			false,
		},
	}

	runFindTest(tests, t, "select_all")
}
