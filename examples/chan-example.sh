#!/usr/bin/gosha

var a chan string

func worker() {
  var val = ""
  for val != "q" {
    val = <-a
    print(val)
  }
}

func generator() {
  var val = ""
  for val != "q" {
    read(&val)
    a <- val
  }
}

go worker()
generator()