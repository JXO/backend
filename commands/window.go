// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"github.com/jxo/lime"
)

type (
	// NewWindow command opens a new window.
	NewWindow struct {
		lime.DefaultCommand
	}

	// CloseAll command closes all the
	// open views inside the current window.
	CloseAll struct {
		lime.DefaultCommand
	}

	// CloseWindow command lets us close the current window.
	CloseWindow struct {
		lime.DefaultCommand
	}

	// NewWindowApp creates a new window, setting it as active.
	NewWindowApp struct {
		lime.DefaultCommand
	}

	// CloseWindowApp command closes all the active windows.
	CloseWindowApp struct {
		lime.DefaultCommand
	}
)

// Run executes the NewWindow command.
func (c *NewWindow) Run(w *lime.Window) error {
	ed := lime.GetEditor()
	ed.SetActiveWindow(ed.NewWindow())
	return nil
}

// Run executes the CloseAll command.
func (c *CloseAll) Run(w *lime.Window) error {
	if !w.CloseAllViews() {
		return fmt.Errorf("Window{id:%d} failed to close all windows", w.ID())
	}
	return nil
}

// Run executes the CloseWindow command.
func (c *CloseWindow) Run(w *lime.Window) error {
	if !w.Close() {
		return fmt.Errorf("Window{id:%d} failed to close window", w.ID())
	}
	return nil
}

// Run executes the NewWindowApp command.
func (c *NewWindowApp) Run() error {
	ed := lime.GetEditor()
	ed.SetActiveWindow(ed.NewWindow())
	return nil
}

// Run executes the CloseWindowApp command.
func (c *CloseWindowApp) Run() error {
	ed := lime.GetEditor()
	if !ed.ActiveWindow().Close() {
		return fmt.Errorf("Failed to close window app")
	}
	return nil
}

// IsChecked shows if NewWindowApp has a
// checkbox in the frontend.
func (c *NewWindowApp) IsChecked() bool {
	return false
}

// IsChecked shows if CloseWindowApp has a
// checkbox in the frontend.
func (c *CloseWindowApp) IsChecked() bool {
	return false
}

func init() {
	register([]lime.Command{
		&NewWindow{},
		&CloseAll{},
		&CloseWindow{},
	})

	registerByName([]namedCmd{
		{"new_window", &NewWindowApp{}},
		{"close_window", &CloseWindowApp{}},
	})
}
