#!/usr/bin/gosha

func readIntAndReturn() int {
  var val int
  read(&val)
  return val
}

var slice []string

slice = append(slice, "1")
print(slice[0])
read(&slice[0])
print(slice[0])

for len(slice) < 2 {
  val := readIntAndReturn()
  slice = append(slice, $(cat ~/.bash_history | tail -n $val))
}

print(slice)