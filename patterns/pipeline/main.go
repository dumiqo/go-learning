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

func take[T any](context context.Context, input chan T, count int) chan T {
	output := make(chan T)

	go func() {
		defer close(output)
		for i := 0; i < count; i++ {
			select {
			case <-context.Done():
				return
			case output <- <-input:
			}
		}
	}()

	return output
}

func prime(context context.Context, input chan int) chan int {
	output := make(chan int)

	go func() {
		defer close(output)
		for {
			select {
			case <-context.Done():
				return
			case i := <-input:
				if isPrime(i) {
					output <- i
				}
			}
		}
	}()

	return output
}

func isPrime(n int) bool {
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	fn := func() int {
		return rand.IntN(1_000_000_000)
	}

	gen := generate(ctx, fn)
	prim := prime(ctx, gen)
	take := take(ctx, prim, 1000)
	for i := range take {
		fmt.Println(i)
	}
}
