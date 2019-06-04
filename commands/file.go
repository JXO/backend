// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"os"
	"os/user"
	"path"

	"github.com/jxo/lime"
)

type (
	// NewFile command creates a new file.
	NewFile struct {
		lime.DefaultCommand
	}
	// PromptOpenFile command prompts opening
	// an existing file from the filesystem.
	PromptOpenFile struct {
		lime.DefaultCommand
	}
)

// Run executes the NewFile command.
func (c *NewFile) Run(w *lime.Window) error {
	ed := lime.GetEditor()
	ed.ActiveWindow().NewFile()
	return nil
}

// Run executes the PromptOpenFile command.
func (o *PromptOpenFile) Run(w *lime.Window) error {
	dir := viewDirectory(w.ActiveView())
	fe := lime.GetEditor().Frontend()
	files := fe.Prompt("Open file", dir, lime.PROMPT_SELECT_MULTIPLE)
	for _, file := range files {
		w.OpenFile(file, 0)
	}
	return nil
}

func viewDirectory(v *lime.View) string {
	if v != nil && v.FileName() != "" {
		p := path.Dir(v.FileName())
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return "/"
}

func init() {
	register([]lime.Command{
		&NewFile{},
		&PromptOpenFile{},
	})
}
