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

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2

		// Если один канал закрыт, а другой нет — деревья разной длины
		if ok1 != ok2 {
			return false
		}
		// Если оба закрыты — все значения сравнены успешно
		if !ok1 {
			return true
		}
		// Сравниваем значения
		if v1 != v2 {
			return false
		}
	}
}
