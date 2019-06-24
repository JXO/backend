// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"testing"

	"github.com/jxo/lime"
	_ "github.com/jxo/lime/sublime/api"
)

func TestPlugin(t *testing.T) {
	newPlugin("testdata/plugin.py").Load()
	pyTest(t, "plugin_test")
}

func pyTest(t *testing.T, imp string) {
	if _, err := pyImport(imp); err != nil {
		t.Errorf("Error importing %s: %s", imp, err)
	}
}

func init() {
	pyAddPath("testdata")

	ed := lime.GetEditor()
	ed.Init()
	ed.NewWindow()
}
