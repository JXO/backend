// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package lime

import "github.com/jxo/lime/render"

// Any color scheme view should implement this interface
// also it should register it self from editor.AddColorScheme
type ColorScheme interface {
	render.ColourScheme
	Name() string
}

type scheme struct {
	settings render.Settings
}

func (s *scheme) Spice(*render.ViewRegions) render.Flavour {
	return render.Flavour{
		Background: s.GlobalSettings().Background,
		Foreground: s.GlobalSettings().Foreground,
	}
}

func (s *scheme) GlobalSettings() render.Settings {
	return s.settings
}

func (s *scheme) Name() string {
	return "Plain theme"
}

// default colorscheme used when there is a problem
var colorscheme *scheme

func defaultScheme() ColorScheme {
	if colorscheme == nil {
		colorscheme = &scheme{
			render.Settings{
				Background: render.Colour{255, 255, 255, 1},
			},
		}
	}
	return colorscheme
}
