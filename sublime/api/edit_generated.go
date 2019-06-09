// Copyright 2017 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// This file was generated by gen_python_api.go and shouldn't be manually modified

package api

import (
	"fmt"

	"github.com/jxo/lime"
	"github.com/limetext/gopy"
	"github.com/jxo/lime/text"
)

var (
	_ = lime.View{}
	_ = text.Region{}
	_ = fmt.Errorf
)

var _editClass = py.Class{
	Name:    "sublime.Edit",
	Pointer: (*Edit)(nil),
}

type Edit struct {
	py.BaseObject
	data *lime.Edit
}

func (o *Edit) PyInit(args *py.Tuple, kwds *py.Dict) error {
	return fmt.Errorf("Can't initialize type Edit")
}
func (o *Edit) PyStr() string {
	return o.data.String()
}