package excercise

import (
	"testing"

	"golang.org/x/tour/tree"
)

func TestWalk(t *testing.T) {
	inputTree := tree.New(1)

	ch := make(chan int)

	go Walk(inputTree, ch)

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, want := range expected {
		v, ok := <-ch
		if !ok {
			t.Fatalf("канал был закрыт преждевременно, ожидалось %d", want)
		}
		if v != want {
			t.Errorf("получено %d, ожидалось %d", v, want)
		}
	}
	_, ok := <-ch
	if ok {
		t.Error("канал не был закрыт после всех значений")
	}
}

func TestSameEqualTree(t *testing.T) {
	t1 := tree.New(1)
	t2 := tree.New(1)
	ok := Same(t1, t2)
	if !ok {
		t.Fatalf("Деревья должны быть равны")
	}
}

func TestSameDifferentTree(t *testing.T) {
	t1 := tree.New(1)
	t2 := tree.New(2)
	ok := Same(t1, t2)
	if ok {
		t.Fatalf("Деревья должны быть разные")
	}
}
