package rbtree

// subIterator implements Iterator interface for the sub tree collection.
type subIterator struct {
	iterator *iterator
	toKey    Item
}

// IsValid returns true if the iterator is valid, otherwise returns false.
func (it *subIterator) IsValid() bool {
	node := it.iterator.node

	return it.iterator.IsValid() && node != nil && node.item != nil && !it.toKey.Less(node.item)
}

// Next moves the iterator to the next element and returns it.
func (it *subIterator) Next() Item {
	if it.iterator.state == pastRear {
		return nil
	}

	item := it.iterator.Next()

	if item != nil && it.toKey.Less(item) {
		it.iterator.state = pastRear
		return nil
	}

	return item
}

// Get returns the current pointed element. Return nil if the iterator is invalid.
func (it *subIterator) Get() Item {
	if !it.IsValid() {
		return nil
	}

	return it.iterator.node.item
}
