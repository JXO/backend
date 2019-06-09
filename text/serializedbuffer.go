// Copyright 2013 Fredrik Ehnbom
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package text

import (
	log "github.com/jxo/log4go"
	"runtime/debug"
)

// SerializedBuffer is a type that serializes all read/write operations
// from/to the inner buffer implementation.
type SerializedBuffer struct {
	inner   InnerBufferInterface
	ops     chan SerializedOperation
	lockret chan interface{}
}

// SerializedOperation is a function return interface {}.
type SerializedOperation func() interface{}

func (s *SerializedBuffer) init(bi InnerBufferInterface) {
	s.inner = bi
	s.ops = make(chan SerializedOperation)
	s.lockret = make(chan interface{})
	go s.worker()
}

// Close closes the buffer operations.
func (s *SerializedBuffer) Close() {
	if s.inner == nil {
		return
	}
	// Close s.ops indicating that we won't send new data on in this channel.
	close(s.ops)
	// The rest of the cleanup will happen once the worker has finished
	// receiving everything on the s.ops channel
}

func (s *SerializedBuffer) worker() {
	for o := range s.ops {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error("Recovered from panic: %v, %s", r, debug.Stack())
					s.lockret <- r
				}
			}()
			s.lockret <- o()
		}()
	}
	// Done processing all ops, so freeing the other resources here
	s.inner = nil
	close(s.lockret)
}

// Size returns the size of the buffer.
func (s *SerializedBuffer) Size() int {
	s.ops <- func() interface{} { return s.inner.Size() }
	r := <-s.lockret
	if r2, ok := r.(int); ok {
		return r2
	}

	return 0
}

// SubstrR returns the runes in the specified Region.
func (s *SerializedBuffer) SubstrR(re Region) []rune {
	s.ops <- func() interface{} { return s.inner.SubstrR(re) }
	r := <-s.lockret
	if r2, ok := r.([]rune); ok {
		return r2
	}

	log.Error("Error: %v", r)
	return nil
}

// InsertR inserts the given rune data at the specified point in the buffer.
func (s *SerializedBuffer) InsertR(point int, data []rune) error {
	s.ops <- func() interface{} { return s.inner.InsertR(point, data) }
	r := <-s.lockret
	if r2, ok := r.(error); ok {
		return r2
	}
	return nil
}

// Erase erases "length" units from "point".
func (s *SerializedBuffer) Erase(point, length int) error {
	s.ops <- func() interface{} { return s.inner.Erase(point, length) }
	r := <-s.lockret

	if r2, ok := r.(error); ok {
		return r2
	}

	return nil
}

// Index returns the rune at the given index
func (s *SerializedBuffer) Index(i int) rune {
	s.ops <- func() interface{} { return s.inner.Index(i) }
	r := <-s.lockret
	if r2, ok := r.(rune); ok {
		return r2
	}

	log.Error("Error: %v", r)
	return 0
}

// RowCol convert a text position into a row and column.
func (s *SerializedBuffer) RowCol(point int) (row, col int) {
	s.ops <- func() interface{} { r, c := s.inner.RowCol(point); return [2]int{r, c} }
	r := <-s.lockret
	if r2, ok := r.([2]int); ok {
		return r2[0], r2[1]
	}

	log.Error("Error: %v", r)
	return 0, 0
}

// TextPoint inverse of #RowCol, converting a row and column into a text position.
func (s *SerializedBuffer) TextPoint(row, col int) (i int) {
	s.ops <- func() interface{} { return s.inner.TextPoint(row, col) }
	r := <-s.lockret
	if r2, ok := r.(int); ok {
		return r2
	}

	log.Error("Error: %v", r)
	return 0
}
