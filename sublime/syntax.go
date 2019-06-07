// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"path/filepath"

	"github.com/jxo/lime/parser"
	"github.com/jxo/lime/sublime/textmate/language"
)

// wrapper around Language implementing lime.Syntax interface
type syntax struct {
	l *language.Language
}

func newSyntax(path string) (*syntax, error) {
	if l, err := language.Load(path); err != nil {
		return nil, err
	} else {
		return &syntax{l: l}, nil
	}
}

func (s *syntax) Parser(data string) (parser.Parser, error) {
	// we can't use syntax language(s.l) because it causes race conditions
	// on concurrent parsing. We could load the language from the file again
	// but I think copying would be much faster
	l := s.l.Copy()
	return language.NewParser(l, []rune(data)), nil
}

func (s *syntax) Name() string {
	return s.l.Name
}

func (s *syntax) FileTypes() []string {
	return s.l.FileTypes
}

func isSyntax(path string) bool {
	if filepath.Ext(path) == ".tmLanguage" {
		return true
	}
	return false
}
