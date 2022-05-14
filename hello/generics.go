package main

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}

	return y
}

type Tree[T interface{}] struct {
	left, right *Tree[T]
	value T
}

func (t *Tree[T]) LookUp(x interface{}) *Tree[T] {
	println(x)
	return t.left
}

func main() {
	min := Min[float64](14.1, 2.0)
	println(min)

	var tree Tree[float64]
	tree.LookUp(1.0)
}