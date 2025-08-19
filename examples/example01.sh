#!/usr/bin/gosha

func a(x int) func(int) int {
  return func(y int) int {
    return x + y
  }
}

print(a(1)(2))