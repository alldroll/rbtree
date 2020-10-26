package rbtree

import "errors"

// ErrorOutOfSubTreeRange tells that there was an attempt to access out of the subtree range.
var ErrorOutOfSubTreeRange error = errors.New("Given key is out of sub tree range")

// subTree is a view of the portion of the tree whose
// keys range from fromKey, inclusive, to toKey, exclusive.
type subTree struct {
	tree    *rbTree
	fromKey Item
	toKey   Item
}

// Returns the number of items in the tree.
func (st *subTree) Len() int {
	iterator := st.NewIterator()
	size := 0

	for iterator.Next() != nil {
		size++
	}

	return size
}

// Insert adds the given item to the tree.
// Returns true if the item was successfully inserted, or returns false if the item was replaced.
// Returns error if there was an attempt to add an element out of subtree range.
func (st *subTree) Insert(item Item) (bool, error) {
	if !st.inRange(item) {
		return false, ErrorOutOfSubTreeRange
	}

	return st.tree.Insert(item)
}

// Removes the given item from the tree
// Returns true if the item was successfuly removes, otherwise returns false
// Returns error if there was an attempt to remove an element out of subtree range.
func (st *subTree) Remove(item Item) (bool, error) {
	if !st.inRange(item) {
		return false, ErrorOutOfSubTreeRange
	}

	return st.tree.Remove(item)
}

// Returns a item if the given key is in the tree, otherwise return nil.
func (st *subTree) Find(item Item) Item {
	if !st.inRange(item) {
		return nil
	}

	return st.tree.Find(item)
}

// Returns the min element in the sub tree
func (st *subTree) Min() Item {
	node := st.tree.root.ceiling(st.fromKey)
	if node == tNil {
		return nil
	}

	return node.item
}

// Returns the max element in the sub tree
func (st *subTree) Max() Item {
	node := st.tree.root.floor(st.toKey)
	if node == tNil {
		return nil
	}

	return node.item
}

// SubTree returns a view of the portion of this tree whose keys range from
// fromKey, inclusive, to toKey, exclusive.
func (st *subTree) NewIterator() Iterator {
	return &subIterator{
		iterator: &iterator{
			node:  st.tree.root.ceiling(st.fromKey),
			state: beforeFirst,
		},
		toKey: st.toKey,
	}
}

// Returns a view of the portion of this map whose keys range from
// fromKey, inclusive, to toKey, exclusive
func (st *subTree) SubTree(fromKey, toKey Item) (Tree, error) {
	if !st.inRange(fromKey) || !st.inRange(toKey) {
		return nil, ErrorOutOfSubTreeRange
	}

	return st.tree.SubTree(fromKey, toKey)
}

// Returns true if the given item in the subTree range, otherwise return false
func (st *subTree) inRange(item Item) bool {
	return !item.Less(st.fromKey) && !st.toKey.Less(item)
}
