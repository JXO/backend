// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import "github.com/jxo/lime"

type (
	// NopApplication performs NOP.
	NopApplication struct {
		lime.BypassUndoCommand
	}
	// NopWindow performs NOP.
	NopWindow struct {
		lime.BypassUndoCommand
	}
	// NopText performs NOP.
	NopText struct {
		lime.BypassUndoCommand
	}
)

// Run executes the NopApplication command.
func (c *NopApplication) Run() error {
	return nil
}

// IsChecked represents if the command
// contains a checkbox in the frontend
func (c *NopApplication) IsChecked() bool {
	return false
}

// Run executes the NopWindow command.
func (c *NopWindow) Run(w *lime.Window) error {
	return nil
}

// Run executes the NopText command.
func (c *NopText) Run(v *lime.View, e *lime.Edit) error {
	return nil
}

func init() {
	registerByName([]namedCmd{
		{"nop", &NopApplication{}},
		{"nop", &NopWindow{}},
		{"nop", &NopText{}},
	})
}
