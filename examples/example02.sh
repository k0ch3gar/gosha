#!/bin/gosha

var all = ""
var cur = ""
for cur != "q" {
  readln(cur)
  if cur == "q" {
    return all
  }

  all = all + cur
}