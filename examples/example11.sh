#!/usr/bin/gosha

var count = 0

for true {
  var num = 0
  read(&num)
  if num % 2 == 0 {
    return count
  }

  count = count + 1
}
