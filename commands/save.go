// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/jxo/lime"
)

type (

	// Save command writes the currently
	// opened file to the disk.
	Save struct {
		lime.DefaultCommand
	}

	// PromptSaveAs command lets us save
	// the currently active
	// file with a different name.
	PromptSaveAs struct {
		lime.DefaultCommand
	}

	// SaveAll command saves all the open files to the disk.
	SaveAll struct {
		lime.DefaultCommand
	}
)

// Run executes the Save command.
func (c *Save) Run(v *lime.View, e *lime.Edit) error {
	err := v.Save()
	if err != nil {
		lime.GetEditor().Frontend().ErrorMessage(fmt.Sprintf("Failed to save %s:n%s", v.FileName(), err))
		return err
	}
	return nil
}

// Run executes the PromptSaveAs command.
func (c *PromptSaveAs) Run(v *lime.View, e *lime.Edit) error {
	dir := viewDirectory(v)
	fe := lime.GetEditor().Frontend()
	files := fe.Prompt("Save file", dir, lime.PROMPT_SAVE_AS)
	if len(files) == 0 {
		return nil
	}

	name := files[0]
	if err := v.SaveAs(name); err != nil {
		fe.ErrorMessage(fmt.Sprintf("Failed to save as %s:%s", name, err))
		return err
	}
	return nil
}

// Run executes the SaveAll command.
func (c *SaveAll) Run(w *lime.Window) error {
	fe := lime.GetEditor().Frontend()
	for _, v := range w.Views() {
		if err := v.Save(); err != nil {
			fe.ErrorMessage(fmt.Sprintf("Failed to save %s:n%s", v.FileName(), err))
			return err
		}
	}
	return nil
}

func init() {
	register([]lime.Command{
		&Save{},
		&PromptSaveAs{},
		&SaveAll{},
	})
}
