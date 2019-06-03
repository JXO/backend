// Copyright 2013 Fredrik Ehnbom
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package text

const chunkSize = 256 * 1024

type (
	naiveBuffer struct {
		data []rune
	}
)

func (b *naiveBuffer) Close() {
}

func (b *naiveBuffer) Size() int {
	return len(b.data)
}

func (b *naiveBuffer) Index(pos int) rune {
	return b.data[pos]
}

func (b *naiveBuffer) SubstrR(r Region) []rune {
	l := len(b.data)
	s, e := Clamp(0, l, r.Begin()), Clamp(0, l, r.End())
	return b.data[s:e]
}

func (b *naiveBuffer) InsertR(point int, value []rune) error {
	point = Clamp(0, len(b.data), point)
	req := len(b.data) + len(value)
	if cap(b.data) < req {
		alloc := (req + chunkSize - 1) &^ (chunkSize - 1)
		n := make([]rune, len(b.data), alloc)
		copy(n, b.data)
		b.data = n
	}
	if point == len(b.data) {
		copy(b.data[point:req], value)
	} else {
		copy(b.data[point+len(value):cap(b.data)], b.data[point:len(b.data)])
		copy(b.data[point:req], value)
	}
	b.data = b.data[:req]
	return nil
}

func (b *naiveBuffer) Erase(point, length int) error {
	if length == 0 {
		return nil
	}
	b.data = append(b.data[0:point], b.data[point+length:len(b.data)]...)
	return nil
}

func (b *naiveBuffer) RowCol(point int) (row, col int) {
	if point < 0 {
		point = 0
	} else if l := b.Size(); point > l {
		point = l
	}

	sub := b.SubstrR(Region{0, point})
	for _, r := range sub {
		if r == '\n' {
			row++
			col = 0
		} else {
			col++
		}
	}
	return
}

func (b *naiveBuffer) TextPoint(row, col int) (i int) {
	if row == 0 && col == 0 {
		return 0
	}
	for l := b.Size(); row > 0 && i < l; i++ {
		if b.data[i] == '\n' {
			row--
		}
	}
	if i < b.Size() {
		return i + col
	}
	return i
}
