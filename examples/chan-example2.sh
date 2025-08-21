#!/usr/bin/gosha

var first = make(chan int, 1)
var second = make(chan string, 1)

func workers() {
  for true {
    val := <- second
    print(val)
    if val == "docker ps" {
      break
    }
  }
}

func workerf() {
  for true {
    i := <- first
    val := $(cat ~/.bash_history | tail -n $i | head -n 1)
    second<- val
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