package main

import (
  "fmt"
)

func main() {
  fmt.Println("return", deferCall())
}

func deferCall() (int) {
  var i int
  defer func() {
    i++
    fmt.Println("defer1", i)
  }()
  defer func() {
    i++
    fmt.Println("defer2", i)
  }()
  return i
}
