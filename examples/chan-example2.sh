#!/usr/bin/gosha

var first = make(chan int, 1)
var second = make(chan string, 1)

func workers() {
  for true {
    val := make(int, <- second)
    print(val)
    if val == 0 {
      break
    }
  }
}

func workerf() {
  for true {
    second<- make(string, <-first)
  }
}

func generator() {
  for true {
    var val int
    read(&val)
    first <- val
  }
}

go workerf()
go generator()
workers()