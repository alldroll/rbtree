package rbtree

type state byte

const (
	deferencable state = iota
	beforeFirst
	pastRear
)

// iterator implements Iterator interface for Tree collection.
type iterator struct {
	node  *node
	state state
}

// IsValid returns true if the iterator is valid, otherwise returns false.
func (it *iterator) IsValid() bool {
	return it.state == deferencable
}

// Next moves the iterator to the next element and returns it.
func (it *iterator) Next() Item {
	if it.state == pastRear || it.node == tNil {
		return nil
	}

	if it.state == beforeFirst {
		it.state = deferencable
		return it.node.item
	}

	if it.node.right != tNil {
		it.node = it.node.right.min()
		return it.node.item
	}

	x := it.node
	y := x.parent
	for y.parent != tNil && y.right == x {
		x, y = y, y.parent
	}

	it.node = y
	if y.right == x {
		it.state = pastRear
		return nil
	}

	return it.node.item
}

// Get returns the current pointed element. Return nil if the iterator is invalid.
func (it *iterator) Get() Item {
	if !it.IsValid() {
		return nil
	}

	return it.node.item
}
