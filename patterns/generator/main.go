package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

func generate[T any](context context.Context, fn func() T) chan T {
	output := make(chan T)

	go func() {
		defer close(output)
		for {
			select {
			case <-context.Done():
				return
			case output <- fn():
			}
		}
	}()

	return output
}

func main() {

	fn := func() int {
		return rand.IntN(100_000_000)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ch := generate(ctx, fn)

	for i := range ch {
		fmt.Println(i)
	}
}
