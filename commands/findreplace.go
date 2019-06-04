// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"errors"
	"github.com/jxo/lime"
	"github.com/jxo/lime/text"
)

type (
	// FindUnderExpand Command extends the selection to the current word
	// if the current selection region is empty.
	// If one character or more is selected, the text buffer is scanned for
	// the next occurrence of the selection and that region too is added to
	// the selection set.
	FindUnderExpand struct {
		lime.DefaultCommand
	}
	// FindNext command searches for the last search term, starting at
	// the end of the last selection in the buffer, and wrapping around. If
	// it finds the term, it clears the current selections and selects the
	// newly-found regions.
	FindNext struct {
		lime.DefaultCommand
	}

	// ReplaceNext Command searches for the "old" argument text,
	// and at the first occurance of the text, replaces it with the
	// "new" argument text. If there are multiple regions, the find
	// starts from the max region.
	ReplaceNext struct {
		lime.DefaultCommand
	}

	FindAll struct {
		lime.DefaultCommand
		SearchText []rune
	}

	ReplaceAll struct {
		lime.DefaultCommand
		SearchText  []rune
		ReplaceText []rune
	}
)

var (
	// Remembers the last sequence of runes searched for.
	lastSearch  []rune
	replaceText string
)

// Run executes the FindUnderExpand command.
func (c *FindUnderExpand) Run(v *lime.View, e *lime.Edit) error {
	sel := v.Sel()
	rs := sel.Regions()

	if sel.HasEmpty() {
		for i, r := range rs {
			if r2 := v.Word(r.A); r2.Size() > r.Size() {
				rs[i] = r2
			}
		}
		sel.Clear()
		sel.AddAll(rs)
		lastSearch = v.SubstrR(rs[len(rs)-1])
		return nil
	}
	last := rs[len(rs)-1]
	lastSearch = v.SubstrR(last)
	r := v.Find(string(lastSearch), last.End(), lime.IGNORECASE|lime.LITERAL)
	if r.A != -1 {
		sel.Add(r)
	}
	return nil
}

func nextSelection(v *lime.View, search string) (text.Region, error) {
	sel := v.Sel()
	rs := sel.Regions()
	last := 0
	wrap := v.Settings().Bool("find_wrap")

	// Regions are not sorted, so finding the last one requires a search.
	for _, r := range rs {
		last = text.Max(last, r.End())
	}

	// Start the search right after the last selection.
	start := last
	r := v.Find(search, start, lime.IGNORECASE|lime.LITERAL)
	// If not found yet and find_wrap setting is true, search
	// from the start of the buffer to our original starting point.
	if r.A == -1 && wrap {
		r = v.Find(search, 0, lime.IGNORECASE|lime.LITERAL)
	}
	// If we found our string, select it.
	if r.A != -1 {
		return r, nil
	}
	return text.Region{-1, -1}, errors.New("Selection not Found")
}

func (c *FindAll) Run(v *lime.View, e *lime.Edit) error {
	if len(c.SearchText) == 0 {
		return nil
	}
	/* Original state of find_wrap is stored and then is made false so that nextSection doesn't go into an infinate loop.
	   Later in the end, it is returned to it's original state.
	*/
	wrap := v.Settings().Bool("find_wrap")
	v.Settings().Set("find_wrap", false)
	search := string(c.SearchText)
	sel := v.Sel()
	sel.Clear()
	for {
		selection, err := nextSelection(v, search)
		if err != nil {
			break
		}
		sel.Add(selection)
	}
	v.Settings().Set("find_wrap", wrap)
	return nil

}

// Run executes the FindNext command.
func (c *FindNext) Run(v *lime.View, e *lime.Edit) error {
	/*
		Correct behavior of FindNext:
			- If there is no previous search, do nothing
			- Find the last region in the buffer, start the
			  search immediately after that.
			- If the search term is found, clear any existing
			  selections, and select the newly-found region.
			- Right now this is doing a case-sensitive search. In ST3
			  that's a setting.
	*/

	// If there is no last search term, nothing to do here.
	if len(lastSearch) == 0 {
		return nil
	}
	newr, err := nextSelection(v, string(lastSearch))
	if err != nil {
		return err
	}
	sel := v.Sel()
	sel.Clear()
	sel.Add(newr)
	return nil
}

func (c *ReplaceAll) Run(v *lime.View, e *lime.Edit) error {
	if len(c.SearchText) == 0 {
		return nil
	}
	search := string(c.SearchText)
	replace := string(c.ReplaceText)
	sel := v.Sel()
	sel.Clear()
	/* Original state of find_wrap is stored and then is made false so that nextSection doesn't go into an infinate loop.
	   Later in the end, it is returned to it's original state.
	*/
	wrap := v.Settings().Bool("find_wrap")
	v.Settings().Set("find_wrap", false)
	for {
		selection, err := nextSelection(v, string(search))
		if err != nil {
			break
		}
		v.Erase(e, selection)
		v.Insert(e, selection.Begin(), replace)
		sel := v.Sel()
		sel.Clear()
		sel.Add(text.Region{selection.Begin(), selection.Begin() + len(replace)})
	}
	v.Settings().Set("find_wrap", wrap)
	return nil
}

// Run executes the ReplaceNext command.
func (c *ReplaceNext) Run(v *lime.View, e *lime.Edit) error {
	// use selection function from find.go to get the next region
	selection, err := nextSelection(v, string(lastSearch))
	if err != nil {
		return err
	}
	v.Erase(e, selection)
	v.Insert(e, selection.Begin(), replaceText)
	return nil
}

func init() {
	register([]lime.Command{
		&FindUnderExpand{},
		&FindNext{},
		&ReplaceNext{},
		&ReplaceAll{},
		&FindAll{},
	})
}
