#!/usr/bin/gosha

if $(pwd) == $HOME {
  print($HOME)
  return 0
} else {
  print("Error: Current file isn't in the working directory.")
  return 1
}

