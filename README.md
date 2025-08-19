# Gosha: A Go Interpreter with Bash Support

## Overview

Gosha is a Go language interpreter implemented in Go, featuring seamless Bash command integration. It provides:

- Interactive Go code execution (REPL environment)
- Go script execution via shebang
- Hybrid Go/Bash command execution

## Installation & Build

```bash
git clone https://github.com/k0ch3gar/gosha.git
cd gosha
go build -tags netgo -ldflags '-extldflags "-static"' ./cmd/gosha
mv ./gosha /usr/bin
```

## Usage

1) Interactive execution:
```bash
$ gosha
Hi user!
That's Gosha!
gosha>>
```

2) Script execution with Shebang:
```bash
(your-script.sh)
#!/usr/bin/gosha

print("Hello world!")
...
```

### Motivation

This project was born from the frustration with Bash's complex syntax.
