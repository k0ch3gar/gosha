#!/usr/bin/gosha

func ff() {
  t := 0;
  for t < 10 {
    print(t);
    t = t + 1;
  }
}

func fs() {
  t := 10;
  for t < 20 {
    print(t);
    t = t + 1;
  }
}

go ff()
fs()