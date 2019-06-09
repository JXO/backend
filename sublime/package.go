// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jxo/lime"
	"github.com/jxo/lime/keys"
	"github.com/jxo/lime/log"
	"github.com/jxo/lime/packages"
	_ "github.com/jxo/lime/sublime/api"
	"github.com/jxo/lime/util"
)

// Represents a sublime package
// TODO: iss#71
type pkg struct {
	dir  string
	name string
	util.HasSettings
	keys.HasKeyBindings
	platformSettings *util.HasSettings
	defaultSettings  *util.HasSettings
	defaultKB        *keys.HasKeyBindings
	plugins          map[string]*plugin
	syntaxes         map[string]*syntax
	colorSchemes     map[string]*colorScheme
}

func newPKG(dir string) packages.Package {
	p := &pkg{
		dir:              dir,
		name:             pkgName(dir),
		platformSettings: new(util.HasSettings),
		defaultSettings:  new(util.HasSettings),
		defaultKB:        new(keys.HasKeyBindings),
		plugins:          make(map[string]*plugin),
		syntaxes:         make(map[string]*syntax),
		colorSchemes:     make(map[string]*colorScheme),
	}

	ed := lime.GetEditor()

	// Initializing settings hierarchy
	// editor <- default <- platform <- user(package)
	p.Settings().SetParent(p.platformSettings)
	p.platformSettings.Settings().SetParent(p.defaultSettings)
	p.defaultSettings.Settings().SetParent(ed)

	// Initializing keybidings hierarchy
	// default <- platform(package) <- editor.default
	edDefault := ed.KeyBindings().Parent().KeyBindings().Parent().KeyBindings().Parent()
	tmp := edDefault.KeyBindings().Parent()
	edDefault.KeyBindings().SetParent(p)
	p.KeyBindings().SetParent(p.defaultKB)
	if tmp != nil {
		p.defaultKB.KeyBindings().SetParent(tmp)
	}

	lime.OnUserPathAdd.Add(p.loadUserSettings)

	return p
}

func (p *pkg) Load() {
	log.Debug("Loading package %s", p.Name())
	p.loadKeyBindings()
	p.loadSettings()
	p.loadUserSettings(lime.GetEditor().UserPath())
	// When we failed on importing sublime_plugin module we continue
	// loading packages but not package plugins
	if module != nil {
		p.loadPlugins()
	}
	// load files that could be anywhere in the package dir like syntax,
	// colour scheme and preferences
	filepath.Walk(p.Path(), p.scan)
}

func (p *pkg) UnLoad() {}

func (p *pkg) Path() string {
	return p.dir
}

func (p *pkg) Name() string {
	return p.name
}

// TODO: how we should watch the package and the files containing?
func (p *pkg) FileCreated(name string) {}

func (p *pkg) loadPlugins() {
	log.Fine("Loading %s plugins", p.Name())
	fis, err := ioutil.ReadDir(p.Path())
	if err != nil {
		log.Warn("Error on reading directory %s, %s", p.Path(), err)
		return
	}
	for _, fi := range fis {
		if isPlugin(fi.Name()) {
			p.loadPlugin(filepath.Join(p.Path(), fi.Name()))
		}
	}
}

func (p *pkg) loadPlugin(path string) {
	pl := newPlugin(path)
	pl.Load()

	p.plugins[path] = pl.(*plugin)
}

func (p *pkg) loadColorScheme(path string) {
	log.Fine("Loading %s package color scheme %s", p.Name(), path)
	cs, err := newColorScheme(path)
	if err != nil {
		log.Warn("Error loading %s color scheme %s: %s", p.Name(), path, err)
		return
	}

	p.colorSchemes[path] = cs
	lime.GetEditor().AddColorScheme(path, cs)
}

func (p *pkg) loadSyntax(path string) {
	log.Fine("Loading %s package syntax %s", p.Name(), path)
	syn, err := newSyntax(path)
	if err != nil {
		log.Warn("Error loading %s syntax: %s", p.Name(), err)
		return
	}

	p.syntaxes[path] = syn
	lime.GetEditor().AddSyntax(path, syn)
}

func (p *pkg) loadKeyBindings() {
	log.Fine("Loading %s keybindings", p.Name())
	ed := lime.GetEditor()

	pt := filepath.Join(p.Path(), "Default.sublime-keymap")
	log.Finest("Loading %s", pt)
	packages.LoadJSON(pt, p.defaultKB.KeyBindings())

	pt = filepath.Join(p.Path(), "Default ("+ed.Plat()+").sublime-keymap")
	log.Finest("Loading %s", pt)
	packages.LoadJSON(pt, p.KeyBindings())
}

func (p *pkg) loadSettings() {
	log.Fine("Loading %s settings", p.Name())
	ed := lime.GetEditor()

	pt := filepath.Join(p.Path(), "Preferences.sublime-settings")
	log.Finest("Loading %s", pt)
	packages.LoadJSON(pt, p.defaultSettings.Settings())

	pt = filepath.Join(p.Path(), "Preferences ("+ed.Plat()+").sublime-settings")
	log.Finest("Loading %s", pt)
	packages.LoadJSON(pt, p.platformSettings.Settings())
}

func (p *pkg) loadUserSettings(dir string) {
	log.Fine("Loading %s user settings", p.Name())
	pt := filepath.Join(dir, p.Name()+".sublime-settings")
	log.Finest("Loading %s", pt)
	packages.LoadJSON(pt, p.Settings())
}

func (p *pkg) scan(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	if isColorScheme(path) {
		p.loadColorScheme(path)
	}
	if isSyntax(path) {
		p.loadSyntax(path)
	}
	return nil
}

func pkgName(dir string) string {
	return filepath.Base(dir)
}

// Any directory in sublime is a package
func isPKG(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		return false
	}

	name := pkgName(dir)
	ed := lime.GetEditor()
	if ignoreds, ok := ed.Settings().Get("ignored_packages").([]interface{}); ok {
		for _, ignored := range ignoreds {
			if ignored == name {
				return false
			}
		}
	}

	return true
}

var packageRecord = &packages.Record{isPKG, newPKG}

func onInit() {
	// Assuming there is a sublime_plugin.py file in the current directory
	// for that we should add current directory to python paths
	// Every package that imports sublime package should have a copy of
	// sublime_plugin.py file in the "." directory
	pyAddPath(".")
	packages.Register(packageRecord)
	var err error
	if module, err = pyImport("sublime_plugin"); err != nil {
		log.Error("Error importing sublime_plugin: %s", err)
		return
	}
	lime.OnPackagesPathAdd.Add(pyAddPath)
	packages.Register(pluginRecord)
}

func init() {
	lime.OnInit.Add(onInit)
}
