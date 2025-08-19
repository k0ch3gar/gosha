#!/bin/gosha

print("1 - nano")
print("2 - vi")
print("3 - links")
print("4 - exit")

var a = 0
for true {
  read(&a)
  if a == 1 {
    $(nano)
  }

  if a == 2 {
    $(vi)
  }

  if a == 3 {
    $(links)
  }

  if a == 4 {
    return 0;
  }
}
