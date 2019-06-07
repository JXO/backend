// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// The packages package handles lime package management.
//
// The key idea of lime packages is modularity/optionality. The core lime
// shouldn't know about tmbundle nor sublime-package, but rather it should
// make it possible to use these and other variants. @quarnster
// Ideally packages implemented in such a way that we can just do:
// import (
// _ "github.com/jxo/lime/sublime/textmate"
// _ "github.com/jxo/lime/sublime"
// _ "github.com/jxo/lime/emacs"
// )
//
// Package type
//
// Each plugin or package that wants to communicate with lime should
// implement this interface.
//
// Record type
//
// For enabling lime to detect and load a package it should register itself as
// a Record
//
package packages
