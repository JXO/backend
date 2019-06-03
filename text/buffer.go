// Copyright 2013 Fredrik Ehnbom
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package text

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

// The InnerBufferInterface defines a minimal
// interface that different buffer implementations
// need to implement.
type InnerBufferInterface interface {
	// Returns the size of a buffer
	Size() int

	// Returns the runes in the specified Region.
	// Implementations should clamp the region
	// as appropriate.
	SubstrR(r Region) []rune

	// Inserts the given rune data at the
	// specified point in the buffer.
	InsertR(point int, data []rune) error

	// Erases "length" units from "point".
	Erase(point, length int) error

	// Returns the rune at the given index
	Index(int) rune

	// Convert a text position into a row and column.
	RowCol(point int) (row, col int)

	// Inverse of #RowCol, converting a row and column
	// into a text position.
	TextPoint(row, col int) (i int)

	// Close the buffer, freeing any associated resources.
	Close()
}

// BufferObserver tracking changes made to a buffer.
// (http://en.wikipedia.org/wiki/Observer_pattern)
//
// Modifying the buffer from within the callback is not allowed and will result in a NOP.
// This is because if there are multiple observers attached to the buffer, they should all first
// be up to date with the current change before another one is introduced.
type BufferObserver interface {
	// Called after Buffer.Erase has executed.
	//
	// Modifying the buffer from within the callback will result in an error/nop, but if
	// it would work the following would restore the buffer's contents as it was before
	// Erase was executed:
	//     bufferChanged.InsertR(regionRemoved.Begin(), dataRemoved)
	Erased(bufferChanged Buffer, regionRemoved Region, dataRemoved []rune)

	// Called after Buffer.Insert/InsertR has executed.
	//
	// Modifying the buffer from within the callback will result in an error/nop, but if
	// it would work the following would restore the buffer's contents as it was before
	// InsertR was executed:
	//     bufferChanged.Erase(regionInserted.Begin(), regionInserted.Size())
	Inserted(bufferChanged Buffer, regionInserted Region, dataInserted []rune)
}

// Buffer defines the full buffer interface
type Buffer interface {
	fmt.Stringer
	InnerBufferInterface
	IDInterface

	// Adds the given observer to this buffer's list of observers
	AddObserver(BufferObserver) error

	// Removes the given observer from this buffer's list of observers
	RemoveObserver(BufferObserver) error

	SetName(string) error
	Name() string
	SetFileName(string) error
	FileName() string

	// Inserts the given string at the given location.
	// Typically just a wrapper around #InsertR
	Insert(point int, svalue string) error

	// Returns the string of the specified Region.
	// Typically just a wrapper around #SubstrR
	Substr(r Region) string

	ChangeCount() int
	// Returns the line region at the given offset
	Line(offset int) Region
	// Returns a Region starting at the start of a line and ending at the end of a (possibly different) line
	LineR(r Region) Region
	// Returns the lines intersecting the region
	Lines(r Region) []Region
	// Like #Line, but includes the line endings
	FullLine(offset int) Region
	// Like #LineR, but includes the line endings
	FullLineR(r Region) Region
	// Returns the word region at the given text offset
	Word(offset int) Region
	// Returns the Region covering the start of the word in r.Begin()
	// to the end of the word in r.End()
	WordR(r Region) Region
}

// The BufferChangedCallback is called everytime a buffer is
// changed.
type BufferChangedCallback func(buf Buffer, position, delta int)

type buffer struct {
	HasID
	SerializedBuffer
	changecount int
	name        string
	filename    string
	callbacks   []BufferChangedCallback
	observers   map[BufferObserver]bool

	inCallbacks int32

	lock sync.Mutex // All lock
}

const (
	wordSeps    = "./\\()\"'-:,.;<>~!@#$%^&*|+=[]{}`~?"
	wordSpaces  = " \n\t\r"
	wordAllSeps = wordSpaces + wordSeps
)

var (
	errObserverAlreadyAdded = fmt.Errorf("Observer has already been added")
	errObserverNotInList    = fmt.Errorf("Observer is not in the list of observers")
	errNothingToInsert      = fmt.Errorf("Nothing to insert")
	errNothingToErase       = fmt.Errorf("Nothing to erase")
	errBufferInCallbacks    = fmt.Errorf("Buffer can not be modified when in a callback")
)

// NewBuffer returns a new empty Buffer
func NewBuffer() Buffer {
	b := buffer{
		observers: make(map[BufferObserver]bool),
	}
	b.SerializedBuffer.init(&rebalancingNode{})
	r := &b
	runtime.SetFinalizer(r, func(b *buffer) { b.Close() })

	return r
}

// implement the fmt.Stringer interface
func (b *buffer) String() string {
	return fmt.Sprintf("Buffer{id: %d, name: \"%s\", filename: \"%s\"}", b.ID(), b.Name(), b.FileName())
}

func (b *buffer) modLock() error {
	if !atomic.CompareAndSwapInt32(&b.inCallbacks, 0, 1) {
		return errBufferInCallbacks
	}
	return nil
}

func (b *buffer) modUnlock() {
	if !atomic.CompareAndSwapInt32(&b.inCallbacks, 1, 0) {
		panic("this shouldn't happen...")
	}
}

func (b *buffer) AddObserver(obs BufferObserver) error {
	if err := b.modLock(); err != nil {
		return err
	}
	defer b.modUnlock()
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.observers[obs] {
		return errObserverAlreadyAdded
	}
	b.observers[obs] = true
	return nil
}

func (b *buffer) RemoveObserver(obs BufferObserver) error {
	if err := b.modLock(); err != nil {
		return err
	}
	defer b.modUnlock()
	b.lock.Lock()
	defer b.lock.Unlock()
	if !b.observers[obs] {
		return errObserverNotInList
	}
	delete(b.observers, obs)
	return nil
}

func (b *buffer) SetName(n string) error {
	if err := b.modLock(); err != nil {
		return err
	}
	defer b.modUnlock()
	b.lock.Lock()
	defer b.lock.Unlock()
	b.name = n
	return nil
}

func (b *buffer) Name() string {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.name
}

func (b *buffer) FileName() string {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.filename
}

func (b *buffer) SetFileName(n string) error {
	if err := b.modLock(); err != nil {
		return err
	}
	defer b.modUnlock()
	b.lock.Lock()
	defer b.lock.Unlock()
	b.filename = n
	return nil
}

func (b *buffer) notify(position, delta int) {
	for i := range b.callbacks {
		b.callbacks[i](b, position, delta)
	}
}

func (b *buffer) InsertR(point int, value []rune) error {
	if err := b.modLock(); err != nil {
		return err
	}
	defer b.modUnlock()
	if err := b.SerializedBuffer.InsertR(point, value); err != nil {
		return err
	}
	b.lock.Lock()
	b.changecount++
	b.lock.Unlock()
	b.notify(point, len(value))
	for obs := range b.observers {
		obs.Inserted(b, Region{point, point + len(value)}, value)
	}
	return nil
}
func (b *buffer) Insert(point int, svalue string) error {
	if len(svalue) == 0 {
		return errNothingToInsert
	}
	value := []rune(svalue)
	return b.InsertR(point, value)
}

func (b *buffer) Erase(point, length int) error {
	if err := b.modLock(); err != nil {
		return err
	}
	defer b.modUnlock()
	if length <= 0 {
		return errNothingToErase
	}
	b.lock.Lock()
	b.changecount++
	b.lock.Unlock()
	re := Region{point, point + length}
	data := b.SubstrR(re)
	// Adjust the region in case original region was longer than the actual buffer
	re.B = re.A + len(data)
	if err := b.SerializedBuffer.Erase(point, length); err != nil {
		return err
	}

	b.notify(point+length, -length)
	for obs := range b.observers {
		obs.Erased(b, re, data)
	}
	return nil
}

func (b *buffer) Substr(r Region) string {
	return string(b.SubstrR(r))
}

func (b *buffer) ChangeCount() int {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.changecount
}

func (b *buffer) Line(offset int) Region {
	if offset < 0 {
		return Region{0, 0}
	} else if s := b.Size(); offset >= s {
		return Region{s, s}
	}
	soffset := offset
sloop:
	o := Clamp(0, soffset, soffset-32)
	sub := b.SubstrR(Region{o, soffset})
	s := soffset
	for s > o && sub[s-o-1] != '\n' {
		s--
	}
	if s == o && o > 0 && sub[0] != '\n' {
		soffset = o
		goto sloop
	}

	l := b.Size()
	eoffset := offset
eloop:
	o = Clamp(eoffset, l, eoffset+32)
	sub = b.SubstrR(Region{eoffset, o})
	e := eoffset
	for e < o && sub[e-eoffset] != '\n' {
		e++
	}
	if e == o && o < l && sub[o-eoffset-1] != '\n' {
		eoffset = o
		goto eloop
	}
	return Region{s, e}
}

func (b *buffer) Lines(r Region) (lines []Region) {
	r = b.FullLineR(r)
	buf := b.SubstrR(r)
	last := r.Begin()
	for i, ru := range buf {
		if ru == '\n' {
			lines = append(lines, Region{last, r.Begin() + i})
			last = r.Begin() + i + 1
		}
	}
	if last != r.End() {
		lines = append(lines, Region{last, r.End()})
	}
	return
}

func (b *buffer) LineR(r Region) Region {
	s := b.Line(r.Begin())
	e := b.Line(r.End())
	return Region{s.Begin(), e.End()}
}

func (b *buffer) FullLine(offset int) Region {
	r := b.Line(offset)
	s := b.Size()
	for r.B < s {
		if i := b.Index(r.B); i == '\r' || i == '\n' {
			break
		}
		r.B++
	}
	if r.B != b.Size() {
		r.B++
	}
	return r
}

func (b *buffer) FullLineR(r Region) Region {
	s := b.FullLine(r.Begin())
	e := b.FullLine(r.End())
	return Region{s.Begin(), e.End()}
}

func (b *buffer) Word(offset int) Region {
	if offset < 0 {
		offset = 0
	}
	lr := b.FullLine(offset)
	col := offset - lr.Begin()

	line := b.SubstrR(lr)
	if len(line) == 0 {
		return Region{offset, offset}
	}

	if col >= len(line) {
		col = len(line) - 1
	}

	last := true
	li := 0
	ls := false
	lc := 0
	for i, r := range line {
		cur := strings.ContainsRune(wordAllSeps, r)
		cs := r == ' '
		if !cs {
			lc = i
		}
		if last == cur && ls == cs {
			continue
		}
		ls = cs
		r := Region{li, i}
		if r.Contains(col) && i != 0 {
			r.A, r.B = r.A+lr.Begin(), r.B+lr.Begin()
			if !(r.B == offset && last) {
				return r
			}
		}
		li = i
		last = cur
	}

	r := Region{lr.Begin() + li, lr.End()}
	lc += lr.Begin()
	if lc != offset && !strings.ContainsRune(wordSpaces, b.Index(r.A)) {
		r.B = lc
	}
	if r.A == offset && r.B == r.A+1 {
		r.B--
	}

	return r
}

// Returns a region that starts at the first character in a word
// and ends with the last character in a (possibly different) word
func (b *buffer) WordR(r Region) Region {
	s := b.Word(r.Begin())
	e := b.Word(r.End())
	return Region{s.Begin(), e.End()}
}
