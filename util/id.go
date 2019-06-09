// Copyright 2013 Fredrik Ehnbom
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package util

import (
	"sync/atomic"
)

// ID is a unique uint32 value
type ID uint32

//IDInterface has an ID method
type IDInterface interface {
	ID() ID
}

// HasID is an utility struct typically embedded to give
// the type a unique id
type HasID struct {
	id ID
}

var (
	idCount = ID(0)
)

// ID return the id from type HasID
func (i *HasID) ID() ID {
	if i.id == 0 {
		i.id = ID(atomic.AddUint32((*uint32)(&idCount), 1))
	}
	return i.id
}
