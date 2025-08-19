#!/usr/bin/gosha

func readIntAndReturn() int {
  var val int
  read(&val)
  return val
}

var slice []int

slice = append(slice, 1)
print(slice[0])
read(&slice[0])
print(slice[0])

for len(slice) < 5 {
  slice = append(slice, readIntAndReturn())
}

print(slice)