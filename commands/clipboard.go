// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"sort"
	"strings"

	"github.com/jxo/lime"
	"github.com/jxo/lime/text"
)

type (
	// Copy copies the current selection to the clipboard. If there
	// are multiple selections, they are concatenated in order from
	// top to bottom of the file, separated by newlines
	Copy struct {
		lime.DefaultCommand
	}

	// Cut copies the current selection to the clipboard, removing it from the
	// buffer. If there are multiple selections, they are concatenated in order
	// from top to bottom of the file, separated by newlines.
	Cut struct {
		lime.DefaultCommand
	}

	// Paste pastes the contents of the clipboard, overwriting the current
	// selection, if any. If there are multiple selections, the clipboard is
	// split into lines. If the number of lines equals the number of selections,
	// the lines are pasted separately into each selection in order from top to
	// bottom of the file. Otherwise the entire clipboard is pasted over every
	// selection.
	Paste struct {
		lime.DefaultCommand
	}
)

func getRegions(v *lime.View, cut bool) *text.RegionSet {
	rs := &text.RegionSet{}
	regions := v.Sel().Regions()
	sort.Sort(regionSorter(regions))
	rs.AddAll(regions)

	he, ae := rs.HasEmpty(), !rs.HasNonEmpty() || cut
	for _, r := range rs.Regions() {
		if ae && r.Empty() {
			rs.Add(v.FullLineR(r))
		} else if he && r.Empty() {
			rs.Subtract(r)
		}
	}

	return rs
}

func getSelForCopy(v *lime.View, rs *text.RegionSet) (s string, ex bool) {
	ss := make([]string, rs.Len())

	for i, r := range rs.Regions() {
		sub := v.Substr(r)

		if !v.Sel().HasNonEmpty() && !strings.HasSuffix(sub, "\n") {
			sub += "\n"
			ex = true
		}

		ss[i] = sub
	}

	s = strings.Join(ss, "\n")

	return
}

// Run executes the Copy command.
func (c *Copy) Run(v *lime.View, e *lime.Edit) error {
	rs := getRegions(v, false)
	s, ex := getSelForCopy(v, rs)

	cb := lime.GetEditor().Clipboard()
	cb.Set(s, ex)

	return nil
}

// Run executes the Cut command.
func (c *Cut) Run(v *lime.View, e *lime.Edit) error {
	s, ex := getSelForCopy(v, getRegions(v, false))

	rs := getRegions(v, true)
	regions := rs.Regions()
	sort.Sort(sort.Reverse(regionSorter(regions)))

	for _, r := range regions {
		v.Erase(e, r)
	}

	cb := lime.GetEditor().Clipboard()
	cb.Set(s, ex)

	return nil
}

// Run executes the Paste command.
func (c *Paste) Run(v *lime.View, e *lime.Edit) error {
	cb := lime.GetEditor().Clipboard()

	rs := &text.RegionSet{}
	regions := v.Sel().Regions()
	sort.Sort(regionSorter(regions))
	rs.AddAll(regions)

	s, ex := cb.Get()

	ss := strings.Split(s, "\n")
	split := !ex && len(ss) == rs.Len()

	for i := rs.Len() - 1; i >= 0; i-- {
		r := rs.Get(i)

		if split {
			v.Replace(e, r, ss[i])
		} else if !ex {
			v.Replace(e, r, s)
		} else {
			l := v.FullLineR(r)
			v.Insert(e, l.Begin(), s)
		}
	}

	return nil
}

func init() {
	register([]lime.Command{
		&Copy{},
		&Cut{},
		&Paste{},
	})
}
