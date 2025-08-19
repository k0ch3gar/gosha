#!/usr/bin/gosha

var slice []int

slice = append(slice, 1)
print(slice[0])
read(&slice[0])
print(slice[0])

for len(slice) < 5 {
  var x int
  read(&x)
  slice = append(slice, x)
}

print(slice)