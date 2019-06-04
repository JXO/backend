// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	"github.com/jxo/lime"
	"github.com/jxo/lime/log"
)

type namedCmd struct {
	name string
	cmd  lime.Command
}

func registerByName(cmds []namedCmd) {
	ch := lime.GetEditor().CommandHandler()
	for _, cmd := range cmds {
		if err := ch.Register(cmd.name, cmd.cmd); err != nil {
			log.Error("Failed to register command %s: %s", cmd.name, err)
		}
	}
}

func register(cmds []lime.Command) {
	ch := lime.GetEditor().CommandHandler()
	for _, cmd := range cmds {
		if err := ch.RegisterWithDefault(cmd); err != nil {
			log.Error("Failed to register command: %s", err)
		}
	}
}
