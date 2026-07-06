package excercise

import "golang.org/x/tour/tree"

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	var walkRec func(t *tree.Tree)
	walkRec = func(t *tree.Tree) {
		if t == nil {
			return
		}
		walkRec(t.Left)
		ch <- t.Value
		walkRec(t.Right)
	}
	walkRec(t)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	if <-ch1 != <-ch2 {
		return false
	}
	return true
}
