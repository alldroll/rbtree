// Package rbtree implements Red-Black tree data structure (RB-Tree).
package rbtree

import (
	"math/rand"
	"sort"
	"testing"
)

// inspired by https://github.com/google/btree
const benchTreeSize = 10000

type IntItem int

func (el IntItem) Less(other Item) bool {
	return el < other.(IntItem)
}

type StringItem string

func (el StringItem) Less(other Item) bool {
	return el < other.(StringItem)
}

var colorNames = map[color]string{
	black: `black`,
	red:   `red`,
}

func TestRotate(t *testing.T) {
	var root, a, b, c, x, y *node

	root = &node{black, nil, nil, tNil, nil}
	a = &node{black, nil, tNil, tNil, nil}
	b = &node{black, nil, tNil, tNil, nil}
	c = &node{black, nil, tNil, tNil, nil}
	x = &node{black, nil, a, tNil, root}
	y = &node{black, nil, b, c, x}

	root.left = x
	x.right = y
	b.parent = y
	c.parent = y
	a.parent = x
	root.parent = tNil

	tree := &rbTree{root, 0}

	tree.leftRotate(x)
	if root.left != y {
		t.Errorf("Expected root.left to be y")
	}

	tree.rightRotate(y)
	if root.left != x {
		t.Errorf("Expected root.left to be x")
	}

	tree.rightRotate(root)
	if tree.root != x {
		t.Errorf("Expected root to be x")
	}
}

// Cormen 13.3.3
func TestInsert(t *testing.T) {
	cases := [][]struct {
		item  IntItem
		color color
	}{
		// nothing
		{
			{41, black},
		},
		// nothing
		{
			{38, red},
			{41, black},
		},
		// case 3
		{
			{31, red},
			{38, black},
			{41, red},
		},
		// case 1, root should be black
		{
			{12, red},
			{31, black},
			{38, black},
			{41, black},
		},
		// case 2 -> case 3
		{
			{12, red},
			{19, black},
			{31, red},
			{38, black},
			{41, black},
		},
		// case 1
		{
			{8, red},
			{12, black},
			{19, red},
			{31, black},
			{38, black},
			{41, black},
		},
	}

	tree := New()
	seq := [...]int{41, 38, 31, 12, 19, 8}
	for i, c := range cases {
		tree.Insert(IntItem(seq[i]))

		nodes := make([]*node, 0)
		iter := tree.NewIterator().(*iterator)
		for {
			val := iter.Next()
			if val == nil {
				break
			}

			nodes = append(nodes, iter.node)
		}

		for j, n := range nodes {
			if n.color != c[j].color || n.item != c[j].item {
				t.Errorf(
					"Expected for %v, item {%d} color to be {%s}, got {%s}",
					seq[:i+1], n.item, colorNames[c[j].color], colorNames[n.color],
				)
			}
		}
	}
}

func TestRemove(t *testing.T) {
	tree := New()
	seq := []int{41, 38, 31, 12, 19, 8}
	for _, item := range seq {
		tree.Insert(IntItem(item))
	}

	cases := [][]struct {
		item  IntItem
		color color
	}{
		{
			{12, black},
			{19, red},
			{31, black},
			{38, black},
			{41, black},
		},
		{
			{19, black},
			{31, red},
			{38, black},
			{41, black},
		},
		{
			{31, black},
			{38, black},
			{41, black},
		},
		{
			{38, black},
			{41, red},
		},
		{
			{41, black},
		},
		{},
	}

	for i, c := range cases {
		item := tree.Min()
		tree.Remove(item)

		nodes := make([]*node, 0)
		iter := tree.NewIterator().(*iterator)
		for {
			item := iter.Next()
			if item == nil {
				break
			}

			nodes = append(nodes, iter.node)
		}

		for j, n := range nodes {
			if n.color != c[j].color || n.item != c[j].item {
				t.Errorf(
					"Expected for case %d, item {%d} color to be {%s}, got {%s}",
					i, n.item, colorNames[c[j].color], colorNames[n.color],
				)
			}
		}
	}

	if tree.Len() != 0 {
		t.Errorf("Expected tree length to be 0, got %d", tree.Len())
	}
}

func TestLength(t *testing.T) {
	tree := New()
	seq := [...]int{41, 38, 31, 12, 19, 8}
	for _, item := range seq {
		tree.Insert(IntItem(item))
	}

	if len(seq) != tree.Len() {
		t.Errorf("Expected tree length to be %d, got %d", len(seq), tree.Len())
	}

	// TODO test removeFixup
}

func TestIterator(t *testing.T) {
	tree := New()
	seq := []int{41, 38, 31, 12, 19, 8, 9, 32, 6, 100, 2, -1, 57, 23, 21, 0, 0, 1}
	for _, item := range seq {
		tree.Insert(IntItem(item))
	}

	expected := []int{-1, 0, 1, 2, 6, 8, 9, 12, 19, 21, 23, 31, 32, 38, 41, 57, 100}
	assertEqualIntDataset(t, tree, expected)
}

func TestSubTreeIterator(t *testing.T) {
	tree := New()
	seq := []int{41, 38, 31, 12, 19, 8, 9, 32, 6, 100, 2, -1, 57, 23, 21, 0, 0, 1}
	for _, item := range seq {
		tree.Insert(IntItem(item))
	}

	subTree, err := tree.SubTree(IntItem(5), IntItem(31))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	expected := []int{6, 8, 9, 12, 19, 21, 23, 31}
	assertEqualIntDataset(t, subTree, expected)
}

func TestMinMax(t *testing.T) {
	tree := New()
	seq := []int{41, 38, 31, 12, 19, 8, 9, 32, 6, 100, 2, -1, 57, 23, 21, 0, 0, 1}
	for _, item := range seq {
		tree.Insert(IntItem(item))
	}

	assertEqualItems(t, IntItem(-1), tree.Min())
	assertEqualItems(t, IntItem(100), tree.Max())

	subTree, err := tree.SubTree(IntItem(11), IntItem(1000))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	assertEqualItems(t, IntItem(12), subTree.Min())
	assertEqualItems(t, IntItem(100), subTree.Max())
}

func TestSubTree(t *testing.T) {
	tree := New()
	seq := [...]int{41, 38, 31, 12, 19, 8, 9, 32, 6, 100, 2, -1, 57, 23, 21, 0, 0, 1}
	for _, item := range seq {
		tree.Insert(IntItem(item))
	}

	subTree, err := tree.SubTree(IntItem(11), IntItem(1000))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	_, err = subTree.SubTree(IntItem(10), IntItem(90))
	if err != ErrorOutOfSubTreeRange {
		t.Errorf("Expected error %v, got %v", ErrorOutOfSubTreeRange, err)
	}

	subTree, err = subTree.SubTree(IntItem(31), IntItem(999))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	expected := []int{31, 32, 38, 41, 57, 100}
	assertEqualIntDataset(t, subTree, expected)
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]byte, n)
	m := len(letters)

	for i := range b {
		b[i] = letters[rand.Intn(m)]
	}

	return string(b)
}

func assertEqualItems(t *testing.T, a, b Item) {
	if (a != b) || a.Less(b) || b.Less(a) {
		t.Errorf("Expected %d, got %d", a, b)
	}
}

func assertEqualIntDataset(t *testing.T, tree Tree, dataset []int) {
	i := 0

	iter := tree.NewIterator()
	for {
		val := iter.Next()
		if val == nil {
			break
		}

		if IntItem(dataset[i]) != val {
			t.Errorf("Expected at {%d} to be %d, got %d", i, dataset[i], val)
		}

		i++
	}

	if i != len(dataset) {
		t.Errorf("Expected to iterate {%d}, got %d", len(dataset), i)
	}
}

func BenchmarkExistsString(b *testing.B) {
	b.StopTimer()

	n := 10
	vals := make([]StringItem, 0, benchTreeSize)

	for i := 0; i < benchTreeSize; i++ {
		vals = append(vals, StringItem(randString(n)))
	}

	b.StartTimer()

	for i := 0; i < b.N; {
		b.StopTimer()

		tree := New()
		for _, val := range vals {
			tree.Insert(val)
		}

		rand.Shuffle(benchTreeSize, func(i, j int) {
			vals[i], vals[j] = vals[j], vals[i]
		})

		b.StartTimer()
		for _, val := range vals {
			tree.Find(val)
			i++

			if i >= b.N {
				break
			}
		}
	}
}

func BenchmarkExistsInt(b *testing.B) {
	b.StopTimer()
	toInsert := perm(benchTreeSize)
	toFind := perm(benchTreeSize)
	b.StartTimer()

	for i := 0; i < b.N; {
		b.StopTimer()

		tree := New()
		for _, val := range toInsert {
			tree.Insert(val)
		}

		b.StartTimer()
		for _, val := range toFind {
			tree.Find(val)
			i++

			if i >= b.N {
				break
			}
		}
	}
}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	vals := perm(benchTreeSize)
	b.StartTimer()

	for i := 0; i < b.N; {
		tree := New()

		for _, v := range vals {
			tree.Insert(v)
			i++

			if i >= b.N {
				break
			}
		}
	}
}

func BenchmarkRemove(b *testing.B) {
	b.StopTimer()
	toInsert := perm(benchTreeSize)
	toDelete := perm(benchTreeSize)
	b.StartTimer()

	for i := 0; i < b.N; {
		b.StopTimer()

		tree := New()
		for _, val := range toInsert {
			tree.Insert(val)
		}

		b.StartTimer()

		for _, val := range toDelete {
			tree.Remove(val)
			i++

			if i >= b.N {
				break
			}
		}
	}

}

func BenchmarkRange(b *testing.B) {
	b.StopTimer()

	vals := perm(benchTreeSize)
	tree := New()

	for _, val := range vals {
		tree.Insert(val)
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Less(vals[j])
	})

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		iter := tree.NewIterator()

		for j := 0; j < benchTreeSize; j++ {
			val := iter.Next()

			if val != vals[j] {
				b.Fatalf("expected: %v, got: %v", vals[j], val)
			}
		}
	}
}

func perm(size int) []IntItem {
	vals := make([]IntItem, 0, size)

	for v := range rand.Perm(size) {
		vals = append(vals, IntItem(v))
	}

	return vals
}
