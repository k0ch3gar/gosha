#!/usr/bin/gosha

if $(pwd) == $(echo $HOME) {
  print($(echo $HOME))
  return 0
} else {
  print("Error: Current file isn't in the working directory.")
  return 1
}

