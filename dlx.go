// Copyright 2014 Zanicar. All rights reserved.

// Utilizes a BSD-3-Clause license. Refer to the included LICENSE file for details.

// Package dlx implements Dancing Links (Algorithm X).
// The algorithm is described in the "Dancing Links" paper by Donald Knuth
// published in "Millennial Perspectives in Computer Science. P159. Volume 187"
// (2000).
package dlx

// Matrix represents a sparse matrix.
// The zero value of a Matrix is an empty matrix ready to use.
type Matrix struct {
	h         Element
	o         []*Element
	solutions [][]string
}

// Init initializes the matrix, empty and ready to use.
func (m *Matrix) Init() *Matrix {
	m.h.up = &m.h
	m.h.down = &m.h
	m.h.left = &m.h
	m.h.right = &m.h
	m.h.column = &m.h
	m.o = nil
	m.solutions = nil
	return m
}

// New returns a pointer to a newly created and initialzed matrix.
func New() *Matrix { return new(Matrix).Init() }

// Head returns the first Head element from the matrix, or nil if empty.
func (m *Matrix) Head() *Element {
	if m.h.right == &m.h {
		return nil
	}
	return m.h.right
}

// Lazy initialization
func (m *Matrix) lazyInit() {
	if m.h.right == nil {
		m.Init()
	}
}

// Helper function to insert a Head element into the matrix and returns a
// pointer to the element.
func (m *Matrix) insertHead(e, at *Element) *Element {
	// Positional pointers
	n := at.right
	at.right = e
	e.left = at
	e.right = n
	n.left = e
	e.up = e
	e.down = e
	// Structural pointers
	e.matrix = m
	e.column = e
	return e
}

// Helper function to insert a given value into the header of the matrix at the
// given head element.
func (m *Matrix) insertHeadValue(v interface{}, at *Element) *Element {
	return m.insertHead(&Element{Value: v}, at)
}

// PushHead pushes a Head element onto the matrix with the given name and
// returns a pointer to the element.
func (m *Matrix) PushHead(name string) *Element {
	m.lazyInit()
	head := Head{name, 0}
	return m.insertHeadValue(head, m.h.left)
}

// Inserts an element at the given row and column and returns a pointer to the
// element.
func (m *Matrix) insertItem(e, atR *Element, atC *Element) *Element {
	if atR == nil {
		e.left = e
		e.right = e
	} else {
		n := atR.right
		atR.right = e
		e.left = atR
		e.right = n
		n.left = e
	}
	ch := atC.down
	atC.down = e
	e.up = atC
	e.down = ch
	ch.up = e
	// Structural pointers
	e.matrix = m
	e.column = ch
	// Update Column Header
	// Utilizes workaround for Go issue 3117
	ch.Value = Head{ch.Value.(Head).name, ch.Value.(Head).size + 1}
	return e
}

// PushItem pushes the given row onto the matrix under the given column head
// element and returns a pointer to the row element.
func (m *Matrix) PushItem(row, colHead *Element) *Element {
	return m.insertItem(&Element{Value: true}, row, colHead.up)
}

// Finds any solutions within the matrix at the given level.
func (m *Matrix) search(k int) {
	if m.Head() == nil {
		solStr := make([]string, len(m.o))
		for i := range m.o {
			j := 0
			rowStr := m.o[i].column.Value.(Head).name
			for e := m.o[i].Right(); e != m.o[i]; e = e.Right() {
				j++
				rowStr += (" " + e.column.Value.(Head).name)
			}
			solStr[i] = rowStr
		}
		m.solutions = append(m.solutions, solStr)
		return
	}
	c := m.getColumn()
	m.cover(c)
	for r := c.Down(); r != c; r = r.Down() {
		m.o = append(m.o, r)
		for j := r.Right(); j != r; j = j.Right() {
			m.cover(j.column)
		}
		m.search(k + 1)
		r = m.o[k]

		m.o[k] = nil
		m.o = m.o[0 : len(m.o)-1]

		c = r.column
		for j := r.Left(); j != r; j = j.Left() {
			m.uncover(j.column)
		}
	}
	m.uncover(c)
}

// Solve invokes a search for solutions from the root (level 0) and returns
// a slice of all found solutions as a slice of strings denoting valid
// constraint options that exactly covers the problem space.
func (m *Matrix) Solve() [][]string {
	m.search(0)
	return m.solutions
}

// Returns a pointer to the head element of the column with the smallest size.
func (m *Matrix) getColumn() *Element {
	var c *Element
	s := uint64(18446744073709551615)
	for ce := m.Head(); ce != nil; ce = ce.Right() {
		ces := ce.Value.(Head).size
		if ces < s {
			c = ce
			s = ces
		}
	}
	return c
}

// The cover operation of algorithm X.
func (m *Matrix) cover(c *Element) {
	c.right.left = c.left
	c.left.right = c.right
	for i := c.Down(); i != c; i = i.Down() {
		for j := i.Right(); j != i; j = j.Right() {
			j.down.up = j.up
			j.up.down = j.down
			j.column.Value = Head{j.column.Value.(Head).name, j.column.Value.(Head).size - 1}
		}
	}
}

// The uncover operation of algorithm X.
func (m *Matrix) uncover(c *Element) {
	for i := c.Up(); i != c; i = i.Up() {
		for j := i.Left(); j != i; j = j.Left() {
			j.column.Value = Head{j.column.Value.(Head).name, j.column.Value.(Head).size + 1}
			j.down.up = j
			j.up.down = j
		}
	}
	c.right.left = c
	c.left.right = c
}

// Element is an element of a matrix. Contains a Value interface{}.
type Element struct {
	// Pointers in the matrix of elements.
	// Column points to the column head.
	up, down, left, right, column *Element

	// The matrix to which the element belongs.
	matrix *Matrix

	Value interface{}
}

// Up returns the above matrix element or nil.
func (e *Element) Up() *Element {
	if p := e.up; e.matrix != nil && p != &e.matrix.h {
		return p
	}
	return nil
}

// Down returns the below matrix element or nil.
func (e *Element) Down() *Element {
	if p := e.down; e.matrix != nil && p != &e.matrix.h {
		return p
	}
	return nil
}

// Left returns the left matrix element or nil.
func (e *Element) Left() *Element {
	if p := e.left; e.matrix != nil && p != &e.matrix.h {
		return p
	}
	return nil
}

// Right returns the right matrix element or nil.
func (e *Element) Right() *Element {
	if p := e.right; e.matrix != nil && p != &e.matrix.h {
		return p
	}
	return nil
}

// Head represents a header element for the matrix.
type Head struct {
	name string
	size uint64
}

// Name returns the column name.
func (h Head) Name() string {
	return h.name
}

// Size returns the column size.
func (h Head) Size() uint64 {
	return h.size
}
