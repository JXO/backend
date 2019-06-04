// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import "github.com/jxo/lime"

type (
	// IncreaseFontSize command increases the font size by 1.
	IncreaseFontSize struct {
		lime.DefaultCommand
	}
	// DecreaseFontSize command decreases the font size by 1.
	DecreaseFontSize struct {
		lime.DefaultCommand
	}
)

// Run executes the IncreaseFontSize command.
func (i *IncreaseFontSize) Run(w *lime.Window) error {
	fontSize := w.Settings().Int("font_size")
	fontSize++
	w.Settings().Set("font_size", fontSize)
	return nil
}

// Run executes the DecreaseFontSize command.
func (d *DecreaseFontSize) Run(w *lime.Window) error {
	fontSize := w.Settings().Int("font_size")
	fontSize--
	w.Settings().Set("font_size", fontSize)
	return nil
}

func init() {
	register([]lime.Command{
		&IncreaseFontSize{},
		&DecreaseFontSize{},
	})
}
