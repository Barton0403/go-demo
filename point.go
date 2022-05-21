package main

import "fmt"

type Tree[T any] struct {
	cmp func(T, T) int
	root *leaf[T]
}

type leaf[T any] struct {
	val T
	left, right *leaf[T]
}

func (bt *Tree[T]) find(val T) {
	pl := &(bt.root)
	fmt.Println(*bt, bt.root)
	println(bt, pl, *pl)
} 

func main() {
	var p *int
	var v int

	v = 1
	p = &v // 指向v的地址
	p2 := &p // 指向p的地址，所以p2也是个指针

	println(v, p, p2, *p2)

	var tree Tree[int]
	tree.find(1)
}