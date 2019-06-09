// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

//go:generate go run gen_python_api.go

package api

import (
	"github.com/jxo/lime"
	"github.com/jxo/lime/log"
	"github.com/jxo/lime/render"
	"github.com/jxo/lime/util"
	"github.com/limetext/gopy"
)

var classes = []struct {
	name string
	c    *py.Class
}{
	{"Region", &_regionClass},
	{"RegionSet", &_region_setClass},
	{"View", &_viewClass},
	{"Window", &_windowClass},
	{"Edit", &_editClass},
	{"Settings", &_settingsClass},
	{"WindowCommandGlue", &_windowCommandGlueClass},
	{"TextCommandGlue", &_textCommandGlueClass},
	{"ApplicationCommandGlue", &_applicationCommandGlueClass},
	{"OnQueryContextGlue", &_onQueryContextGlueClass},
	{"ViewEventGlue", &_viewEventGlueClass},
}

var constants = []struct {
	name     string
	constant int
}{
	{"OP_EQUAL", int(util.OpEqual)},
	{"OP_NOT_EQUAL", int(util.OpNotEqual)},
	{"OP_REGEX_MATCH", int(util.OpRegexMatch)},
	{"OP_NOT_REGEX_MATCH", int(util.OpNotRegexMatch)},
	{"OP_REGEX_CONTAINS", int(util.OpRegexContains)},
	{"OP_NOT_REGEX_CONTAINS", int(util.OpNotRegexContains)},
	{"INHIBIT_WORD_COMPLETIONS", 0},
	{"INHIBIT_EXPLICIT_COMPLETIONS", 0},
	{"LITERAL", int(lime.IGNORECASE)},
	{"IGNORECASE", int(lime.LITERAL)},
	{"CLASS_WORD_START", int(lime.CLASS_WORD_START)},
	{"CLASS_WORD_END", int(lime.CLASS_WORD_END)},
	{"CLASS_PUNCTUATION_START", int(lime.CLASS_PUNCTUATION_START)},
	{"CLASS_PUNCTUATION_END", int(lime.CLASS_PUNCTUATION_END)},
	{"CLASS_SUB_WORD_START", int(lime.CLASS_SUB_WORD_START)},
	{"CLASS_SUB_WORD_END", int(lime.CLASS_SUB_WORD_END)},
	{"CLASS_LINE_START", int(lime.CLASS_LINE_START)},
	{"CLASS_LINE_END", int(lime.CLASS_LINE_END)},
	{"CLASS_EMPTY_LINE", int(lime.CLASS_EMPTY_LINE)},
	{"CLASS_MIDDLE_WORD", int(lime.CLASS_MIDDLE_WORD)},
	{"CLASS_WORD_START_WITH_PUNCTUATION", int(lime.CLASS_WORD_START_WITH_PUNCTUATION)},
	{"CLASS_WORD_END_WITH_PUNCTUATION", int(lime.CLASS_WORD_END_WITH_PUNCTUATION)},
	{"CLASS_OPENING_PARENTHESIS", int(lime.CLASS_OPENING_PARENTHESIS)},
	{"CLASS_CLOSING_PARENTHESIS", int(lime.CLASS_CLOSING_PARENTHESIS)},
	{"DRAW_EMPTY", int(render.DRAW_EMPTY)},
	{"HIDE_ON_MINIMAP", int(render.HIDE_ON_MINIMAP)},
	{"DRAW_EMPTY_AS_OVERWRITE", int(render.DRAW_EMPTY_AS_OVERWRITE)},
	{"DRAW_NO_FILL", int(render.DRAW_NO_FILL)},
	{"DRAW_NO_OUTLINE", int(render.DRAW_NO_OUTLINE)},
	{"DRAW_SOLID_UNDERLINE", int(render.DRAW_SOLID_UNDERLINE)},
	{"DRAW_STIPPLED_UNDERLINE", int(render.DRAW_STIPPLED_UNDERLINE)},
	{"DRAW_SQUIGGLY_UNDERLINE", int(render.DRAW_SQUIGGLY_UNDERLINE)},
	{"PERSISTENT", int(render.PERSISTENT)},
	{"HIDDEN", int(render.HIDDEN)},
}

func init() {
	l := py.InitAndLock()
	defer l.Unlock()

	if sys, err := py.Import("sys"); err != nil {
		log.Warn("Couldn't import sys: %s", err)
	} else {
		if pyc, err := py.NewUnicode("dont_write_bytecode"); err != nil {
			log.Warn(err)
		} else {
			// avoid generating pyc files
			sys.Base().SetAttr(pyc, py.True)
		}
		sys.Decref()
	}

	methods := append(generated_methods, manual_methods...)
	m, err := py.InitModule("sublime", methods)
	if err != nil {
		// TODO: we should handle this as error
		panic(err)
	}

	for _, cl := range classes {
		c, err := cl.c.Create()
		if err != nil {
			panic(err)
		}
		if err := m.AddObject(cl.name, c); err != nil {
			panic(err)
		}
	}

	for _, c := range constants {
		if err := m.AddIntConstant(c.name, c.constant); err != nil {
			panic(err)
		}
	}
}
