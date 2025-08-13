#!/bin/gosha

var all = ""
var cur = ""
for cur != "q" {
  read(cur)
  if cur == "q" {
    return all
  }

  all = all + cur
}
