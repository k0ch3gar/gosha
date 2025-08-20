#!/usr/bin/gosha

func fn() {
  i := 1
  for i < 10 {
    print($(cat ~/.bash_history | tail -n $i | head -n 1))
    i = i + 1
  }
}


func fns() {
  i := 1
  for i < 5 {
    print($(docker ps | tail -n $i | head -n 1))
    i = i + 1
  }
}

go fn()
fns()