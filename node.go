package rbtree

type color byte

const (
	red color = iota
	black
)

type node struct {
	color               color
	item                Item
	left, right, parent *node
}

var tNil = &node{black, nil, nil, nil, nil}

// Returns the min element for this node.
func (nd *node) min() *node {
	n := nd
	for n.left != tNil {
		n = n.left
	}

	return n
}

// Returns the max element for this node.
func (nd *node) max() *node {
	n := nd
	for n.right != tNil {
		n = n.right
	}

	return n
}

// Inspired by java.util.TreeMap#getCeilingEntry
// Gets the node corresponding to the specified item; if no such node
// exists, returns the node for the least item greater than the specified
// item; otherwise returns tNil
func (nd *node) ceiling(item Item) *node {
	p := nd
	for p != tNil {
		if item.Less(p.item) {
			if p.left != tNil {
				p = p.left
			} else {
				return p
			}
		} else if p.item.Less(item) {
			if p.right != tNil {
				p = p.right
			} else {
				parent := p.parent
				ch := p
				for parent != tNil && ch == parent.right {
					ch = parent
					parent = parent.parent
				}

				return parent
			}
		} else {
			return p
		}
	}

	return tNil
}

// Inspired by java.util.TreeMap#getFloorEntry
// Gets the node corresponding to the specified item; if no such node
// exists, returns the node for the greatest item less than the specified
// item; otherwise returns tNil
func (nd *node) floor(item Item) *node {
	p := nd
	for p != tNil {
		if p.item.Less(item) {
			if p.right != tNil {
				p = p.right
			} else {
				return p
			}
		} else if item.Less(p.item) {
			if p.left != tNil {
				p = p.left
			} else {
				parent := p.parent
				ch := p
				for parent != tNil && ch == parent.left {
					ch = parent
					parent = parent.parent
				}

				return parent
			}
		} else {
			return p
		}
	}

	return tNil
}
