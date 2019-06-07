// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package api

import (
	"fmt"

	"github.com/jxo/lime"
	"github.com/jxo/lime/log"
	"github.com/limetext/gopy"
	"github.com/jxo/lime/text"
	"github.com/jxo/lime/util"
)

var (
	_ = lime.View{}
	_ = text.Region{}
)

var (
	_onQueryContextGlueClass = py.Class{
		Name:    "sublime.OnQueryContextGlue",
		Pointer: (*OnQueryContextGlue)(nil),
	}
	_viewEventGlueClass = py.Class{
		Name:    "sublime.ViewEventGlue",
		Pointer: (*ViewEventGlue)(nil),
	}
)

type (
	OnQueryContextGlue struct {
		py.BaseObject
		inner py.Object
	}
	ViewEventGlue struct {
		py.BaseObject
		inner py.Object
	}
)

var evmap = map[string]*lime.ViewEvent{
	"on_new":                &lime.OnNew,
	"on_load":               &lime.OnLoad,
	"on_activated":          &lime.OnActivated,
	"on_deactivated":        &lime.OnDeactivated,
	"on_pre_close":          &lime.OnPreClose,
	"on_close":              &lime.OnClose,
	"on_pre_save":           &lime.OnPreSave,
	"on_post_save":          &lime.OnPostSave,
	"on_modified":           &lime.OnModified,
	"on_selection_modified": &lime.OnSelectionModified,
}

func (c *ViewEventGlue) PyInit(args *py.Tuple, kwds *py.Dict) error {
	if args.Size() != 2 {
		return fmt.Errorf("Expected 2 arguments not %d", args.Size())
	}
	if v, err := args.GetItem(0); err != nil {
		return err
	} else {
		c.inner = v
	}
	if v, err := args.GetItem(1); err != nil {
		return err
	} else if v2, ok := v.(*py.Unicode); !ok {
		return fmt.Errorf("Second argument not a string: %v", v)
	} else {
		ev := evmap[v2.String()]
		if ev == nil {
			return fmt.Errorf("Unknown event: %s", v2)
		}
		ev.Add(c.onEvent)
		c.inner.Incref()
		c.Incref()
	}
	return nil
}

func (c *ViewEventGlue) onEvent(v *lime.View) {
	l := py.NewLock()
	defer l.Unlock()
	pv, err := toPython(v)
	if err != nil {
		log.Error(err)
	}
	defer pv.Decref()
	log.Fine("onEvent: %v, %v, %v", c, c.inner, pv)

	if ret, err := c.inner.Base().CallFunctionObjArgs(pv); err != nil {
		log.Error(err)
	} else if ret != nil {
		ret.Decref()
	}
}

func (c *OnQueryContextGlue) PyInit(args *py.Tuple, kwds *py.Dict) error {
	if args.Size() != 1 {
		return fmt.Errorf("Expected only 1 argument not %d", args.Size())
	}
	if v, err := args.GetItem(0); err != nil {
		return err
	} else {
		c.inner = v
	}
	c.inner.Incref()
	c.Incref()

	lime.OnQueryContext.Add(c.onQueryContext)
	return nil
}

func (c *OnQueryContextGlue) onQueryContext(v *lime.View, key string, operator util.Op, operand interface{}, match_all bool) lime.QueryContextReturn {
	l := py.NewLock()
	defer l.Unlock()

	var (
		pv, pk, po, poa, pm, ret py.Object
		err                      error
	)
	if pv, err = toPython(v); err != nil {
		log.Error(err)
		return lime.Unknown
	}
	defer pv.Decref()

	if pk, err = toPython(key); err != nil {
		log.Error(err)
		return lime.Unknown
	}
	defer pk.Decref()

	if po, err = toPython(operator); err != nil {
		log.Error(err)
		return lime.Unknown
	}
	defer po.Decref()

	if poa, err = toPython(operand); err != nil {
		log.Error(err)
		return lime.Unknown
	}
	defer poa.Decref()

	if pm, err = toPython(match_all); err != nil {
		log.Error(err)
		return lime.Unknown
	}
	defer pm.Decref()

	if ret, err = c.inner.Base().CallFunctionObjArgs(pv, pk, po, poa, pm); err != nil {
		log.Error(err)
		return lime.Unknown
	}
	defer ret.Decref()

	if r2, ok := ret.(*py.Bool); ok {
		if r2.Bool() {
			return lime.True
		} else {
			return lime.False
		}
	} else {
		log.Fine("other: %v", ret)
	}
	return lime.Unknown
}
