// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package syntax

import "github.com/jxo/lime/sublime/textmate"

type (
	Context []Pattern

	Pattern struct {
		Include              string
		MetaScope            string `yaml:"meta_scope"`
		MetaContentScope     string `yaml:"meta_content_scope"`
		MetaIncludePrototype string `yaml:"meta_include_prototype"`
		Match                textmate.Regex
		Scope                string
		Captures             textmate.Captures
		Push                 Context
		Pop                  bool
		Set                  string
		Syntax               string
	}
)
