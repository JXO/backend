// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/jxo/lime"
)

const limeCmdMark = "lime.cmd.mark"

type (
	// MarkUndoGroupsForGluing Command marks the current position
	// in the undo stack as the start of commands to glue, potentially
	// overwriting any existing marks.
	MarkUndoGroupsForGluing struct {
		lime.BypassUndoCommand
	}

	// GlueMarkedUndoGroups Command merges commands from the previously
	// marked undo stack location to the current location into a single
	// entry in the undo stack.
	GlueMarkedUndoGroups struct {
		lime.BypassUndoCommand
	}

	// MaybeMarkUndoGroupsForGluing Command is similar to
	// MarkUndoGroupsForGluingCommand with the exception that if there
	// is already a mark set, it is not overwritten.
	MaybeMarkUndoGroupsForGluing struct {
		lime.BypassUndoCommand
	}

	// UnmarkUndoGroupsForGluing Command removes the glue mark set by
	// either MarkUndoGroupsForGluingCommand or MaybeMarkUndoGroupsForGluingCommand
	// if it was set.
	UnmarkUndoGroupsForGluing struct {
		lime.BypassUndoCommand
	}
)

// Run executes the MarkUndoGroupsForGluing command.
func (c *MarkUndoGroupsForGluing) Run(v *lime.View, e *lime.Edit) error {
	v.Settings().Set(limeCmdMark, v.UndoStack().Position())
	return nil
}

// Run executes the UnmarkUndoGroupsForGluing command.
func (c *UnmarkUndoGroupsForGluing) Run(v *lime.View, e *lime.Edit) error {
	v.Settings().Erase(limeCmdMark)
	return nil
}

// Run executes the GlueMarkedUndoGroups command.
func (c *GlueMarkedUndoGroups) Run(v *lime.View, e *lime.Edit) error {
	pos := v.UndoStack().Position()
	mark, ok := v.Settings().Get(limeCmdMark).(int)
	if !ok {
		return fmt.Errorf("No mark in the current view")
	}
	if l, p := pos-mark, mark; p != -1 && (l-p) > 1 {
		v.UndoStack().GlueFrom(mark)
	}
	return nil
}

// Run executes the MaybeMarkUndoGroupsForGluing command.
func (c *MaybeMarkUndoGroupsForGluing) Run(v *lime.View, e *lime.Edit) error {
	if !v.Settings().Has(limeCmdMark) {
		v.Settings().Set(limeCmdMark, v.UndoStack().Position())
	}
	return nil
}

func init() {
	register([]lime.Command{
		&MarkUndoGroupsForGluing{},
		&GlueMarkedUndoGroups{},
		&MaybeMarkUndoGroupsForGluing{},
		&UnmarkUndoGroupsForGluing{},
	})
}
