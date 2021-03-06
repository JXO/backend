// Copyright 2019 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// This file was generated by gen_python_api.go and shouldn't be manually modified

package api

import (
	"fmt"

	"github.com/jxo/lime"
	"github.com/jxo/lime/text"
	"github.com/jxo/lime/util"
	"github.com/limetext/gopy"
)

var (
	_ = lime.View{}
	_ = text.Region{}
	_ = fmt.Errorf
	_ = util.Settings{}
)

func sublime_Register(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
		arg2 interface{}
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for lime.commandHandler.Register() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	if v, err := tu.GetItem(1); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(interface{}); !ok {
				return nil, fmt.Errorf("Expected type interface {} for lime.commandHandler.Register() arg2, not %s", v.Type())
			} else {
				arg2 = v2
			}
		}
	}
	if err := lime.GetEditor().CommandHandler().Register(arg1, arg2); err != nil {
		return nil, err
	} else {
		return toPython(nil)
	}
}

func sublime_Unregister(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for lime.commandHandler.Unregister() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	if err := lime.GetEditor().CommandHandler().Unregister(arg1); err != nil {
		return nil, err
	} else {
		return toPython(nil)
	}
}
