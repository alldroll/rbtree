// Package rbtree implements Red-Black tree data structure (RB-Tree).
package rbtree

import "errors"

// ErrorFromGreaterThanToKey informs that the fromKey should be less or equal to toKey
var ErrorFromGreaterThanToKey error = errors.New("fromKey should be >= toKey")

// rBTree is an implementation of red-black tree.
type rbTree struct {
	root   *node
	length int
}

// New returns a new instance of Tree.
func New() Tree {
	return &rbTree{tNil, 0}
}

// Returns the number of items in the tree.
func (rb *rbTree) Len() int {
	return rb.length
}

// Insert adds the given item to the tree.
// Returns true if the item was successfully inserted, or returns false if the item was replaced.
// Returns an error if there was an attempt to add an element out of subtree range.
func (rb *rbTree) Insert(item Item) (bool, error) {
	z := &node{red, item, tNil, tNil, tNil}
	res := rb.insert(z)
	result := false

	if res == z { // if we insert z
		rb.length++
		result = true
	}

	return result, nil
}

// Remove deletes an item equals to the given item from the tree.
// Returns true if the item was successfully removes, otherwise returns false.
// Returns an error if there was an attempt to remove an element out of subtree range.
func (rb *rbTree) Remove(item Item) (bool, error) {
	z, _ := rb.find(item)
	if z == tNil {
		return false, nil
	}

	rb.remove(z)
	rb.length--
	return true, nil
}

// Returns a item if the given key is in the tree, otherwise return nil.
func (rb *rbTree) Find(item Item) Item {
	x, _ := rb.find(item)
	if x == tNil {
		return nil
	}

	return x.item
}

// Returns the min element in the tree
func (rb *rbTree) Min() Item {
	return rb.root.min().item
}

// Returns the max element in the tree
func (rb *rbTree) Max() Item {
	return rb.root.max().item
}

// Returns an iterator that points at the smallest element in the tree.
func (rb *rbTree) NewIterator() Iterator {
	if rb.Len() == 0 {
		return &iterator{tNil, beforeFirst}
	}

	return &iterator{rb.root.min(), beforeFirst}
}

// SubTree returns a view of the portion of this tree whose keys range from
// fromKey, inclusive, to toKey, exclusive.
func (rb *rbTree) SubTree(fromKey, toKey Item) (Tree, error) {
	if toKey.Less(fromKey) {
		return nil, ErrorFromGreaterThanToKey
	}

	return &subTree{
		tree:    rb,
		fromKey: fromKey,
		toKey:   toKey,
	}, nil
}

// insert adds the given node in the tree.
func (rb *rbTree) insert(z *node) *node {
	x, y := rb.find(z.item)
	if x != tNil {
		x.item = z.item
		return x
	}

	z.parent = y
	if y == tNil {
		rb.root = z
	} else if z.item.Less(y.item) {
		y.left = z
	} else {
		y.right = z
	}

	z.color = red
	z.left = tNil
	z.right = tNil

	rb.insertFixup(z)
	return z
}

// remove deletes the given node from the tree.
func (rb *rbTree) remove(z *node) {
	x, y := tNil, z
	yColor := y.color

	if z.left == tNil {
		x = z.right
		rb.transplant(z, z.right)
	} else if z.right == tNil {
		x = z.left
		rb.transplant(z, z.left)
	} else {
		y = z.right.min()
		yColor = y.color
		x = y.right
		if y.parent == z {
			x.parent = y
		} else {
			rb.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}

		rb.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}

	if yColor == black {
		rb.removeFixup(x)
	}
}

// find searches the node if the given key is in the tree, otherwise return nil.
func (rb *rbTree) find(item Item) (*node, *node) {
	x := rb.root
	y := tNil

	for x != tNil {
		if item.Less(x.item) {
			y, x = x, x.left
		} else if x.item.Less(item) {
			y, x = x, x.right
		} else {
			break
		}
	}

	return x, y
}

// Performs fixup with insertion
func (rb *rbTree) insertFixup(z *node) {
	for z.parent.color == red {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y.color == red { // case 1, uncle "y" is red
				// restore rule 4
				z.parent.color = black
				y.color = black // uncle "y" should be black
				// uphold rule 5
				z.parent.parent.color = red
				// z - grandparent
				z = z.parent.parent
			} else {
				if z == z.parent.right { // case 2 -> case 3
					z = z.parent
					rb.leftRotate(z)
				}

				// case 3, uncle "y" is black and z - left child
				z.parent.color = black
				z.parent.parent.color = red
				rb.rightRotate(z.parent.parent)
			}
		} else {
			y := z.parent.parent.left
			if y.color == red {
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					rb.rightRotate(z)
				}

				z.parent.color = black
				z.parent.parent.color = red
				rb.leftRotate(z.parent.parent)
			}
		}
	}

	rb.root.color = black
}

// leftRotate performs the left rotation for given node.
func (rb *rbTree) leftRotate(x *node) {
	y := x.right
	x.right = y.left
	if y.left != tNil {
		y.left.parent = x
	}

	y.parent = x.parent
	if x.parent == tNil {
		rb.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.left = x
	x.parent = y
}

// rightRotate performs the right rotation for given node.
func (rb *rbTree) rightRotate(y *node) {
	x := y.left
	y.left = x.right
	if x.right != tNil {
		x.right.parent = y
	}

	x.parent = y.parent
	if y.parent == tNil {
		rb.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	x.right = y
	y.parent = x
}

// removeFixup deletes the given node and performs fixup of the tree.
func (rb *rbTree) removeFixup(x *node) {
	for x != rb.root && x.color == black {
		if x == x.parent.left {
			w := x.parent.right //right brother
			if w.color == red {
				// case 1
				w.color = black
				x.parent.color = red
				rb.leftRotate(x.parent)
				w = x.parent.right
			}

			if w.left.color == black && w.right.color == black {
				// case 2
				w.color = red
				x = x.parent
			} else {
				if w.right.color == black {
					// case 3
					w.left.color = black
					w.color = red
					rb.rightRotate(w)
					x = x.parent
				}
				// case 4
				w.color = x.parent.color
				x.parent.color = black
				w.right.color = black
				rb.leftRotate(x.parent)
				x = rb.root
			}
		} else {
			w := x.parent.left //left brother
			if w.color == red {
				// case 1
				w.color = black
				x.parent.color = red
				rb.rightRotate(x.parent)
				w = x.parent.left
			}

			if w.right.color == black && w.left.color == black {
				// case 2
				w.color = red
				x = x.parent
			} else {
				if w.left.color == black {
					// case 3
					w.right.color = black
					w.color = red
					rb.leftRotate(w)
					x = x.parent
				}
				// case 4
				w.color = x.parent.color
				x.parent.color = black
				w.left.color = black
				rb.rightRotate(x.parent)
				x = rb.root
			}
		}
	}

	x.color = black
}

// transplant performs the transplant operation.
func (rb *rbTree) transplant(u, v *node) {
	if u.parent == tNil {
		rb.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}

	v.parent = u.parent
}
