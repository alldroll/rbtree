// Package rbtree implements Red-Black tree data structure (RB-Tree).
package rbtree

// Tree represents Red-Black tree.
type Tree interface {
	// Returns the number of items in the tree.
	Len() int
	// Insert adds the given item to the tree.
	// Returns true if the item was successfully inserted, or returns false if the item was replaced.
	// Returns an error if there was an attempt to add an element out of subtree range.
	Insert(item Item) (bool, error)
	// Remove deletes an item equals to the given item from the tree.
	// Returns true if the item was successfully removes, otherwise returns false.
	// Returns an error if there was an attempt to remove an element out of subtree range.
	Remove(item Item) (bool, error)
	// Returns the item if the given key is in the tree, otherwise return nil.
	Find(item Item) Item
	// Returns the min element in the tree.
	Min() Item
	// Returns the max element in the tree.
	Max() Item
	// Returns an iterator that points at the smallest element in the tree.
	NewIterator() Iterator
	// SubTree returns a view of the portion of this tree whose keys range from
	// fromKey, inclusive, to toKey, exclusive.
	SubTree(fromKey Item, toKey Item) (Tree, error)
}

// Item represents a single object in the tree.
type Item interface {
	// Less tells whether the current element is less than the given argument.
	Less(other Item) bool
}

// Iterator represents an iterator over a tree collection which provides inorder traverse.
type Iterator interface {
	// IsValid returns true if the iterator is valid, otherwise returns false.
	IsValid() bool
	// Next moves the iterator to the next element and returns it.
	Next() Item
	// Get returns the current pointed element. Return nil if the iterator is invalid.
	Get() Item
}
