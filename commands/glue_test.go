// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"testing"

	"github.com/jxo/lime"
	"github.com/jxo/lime/text"
)

func TestGlueCmds(t *testing.T) {
	ed := lime.GetEditor()
	ch := ed.CommandHandler()
	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	v.SetScratch(true)
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)
	v.SetScratch(false)
	ch.RunTextCommand(v, "mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", lime.Args{"characters": "a"})
	ch.RunTextCommand(v, "insert", lime.Args{"characters": "b"})
	ch.RunTextCommand(v, "insert", lime.Args{"characters": "c"})
	ch.RunTextCommand(v, "glue_marked_undo_groups", nil)
	if v.UndoStack().Position() != 1 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Substr(text.Region{A: 0, B: v.Size()}); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "undo", nil)
	if d := v.Substr(text.Region{A: 0, B: v.Size()}); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "redo", nil)
	if d := v.Substr(text.Region{A: 0, B: v.Size()}); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	if v.UndoStack().Position() != 1 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Substr(text.Region{A: 0, B: v.Size()}); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "undo", nil)
	if d := v.Substr(text.Region{A: 0, B: v.Size()}); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}

	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", lime.Args{"characters": "a"})
	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", lime.Args{"characters": "b"})
	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", lime.Args{"characters": "c"})
	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "glue_marked_undo_groups", nil)
	if v.UndoStack().Position() != 1 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Substr(text.Region{A: 0, B: v.Size()}); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "undo", nil)
	if d := v.Substr(text.Region{A: 0, B: v.Size()}); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "redo", nil)
	if d := v.Substr(text.Region{A: 0, B: v.Size()}); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	if v.UndoStack().Position() != 1 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Substr(text.Region{A: 0, B: v.Size()}); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}

	ch.RunTextCommand(v, "mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "move", lime.Args{"forward": false, "extend": true, "by": "lines"})
	ch.RunTextCommand(v, "move", lime.Args{"forward": false, "extend": true, "by": "lines"})
	ch.RunTextCommand(v, "move", lime.Args{"forward": false, "extend": true, "by": "lines"})
	ch.RunTextCommand(v, "left_delete", nil)
	ch.RunTextCommand(v, "insert", lime.Args{"characters": "a"})
	ch.RunTextCommand(v, "insert", lime.Args{"characters": "b"})
	ch.RunTextCommand(v, "insert", lime.Args{"characters": "c"})
	ch.RunTextCommand(v, "glue_marked_undo_groups", nil)
	if v.UndoStack().Position() != 2 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Substr(text.Region{A: 0, B: v.Size()}); d != "Helabc" {
		t.Error(d)
	}

	ch.RunTextCommand(v, "undo", nil)
	if v.UndoStack().Position() != 1 {
		t.Error(v.UndoStack().Position())
	} else if d := v.Substr(text.Region{A: 0, B: v.Size()}); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
}
