# Gosha: Интерпретатор Go с поддержкой Bash

## Что это?

Gosha - это интерпретатор языка Go, написанный на самом Go, с дополнительной поддержкой выполнения Bash-команд. Он позволяет:

- Выполнять Go-код в интерактивном режиме (REPL)
- Запускать Go-скрипты используя шебанг
- Смешивать Go-код с Bash-командами

## Установка и сборка

```bash
git clone https://github.com/your-username/gosha.git
cd gosha
go build -tags netgo -ldflags '-extldflags "-static"' ./cmd/gosha
mv ./gosha /usr/bin
```

## Использование

1) Использование в интерактивном режиме:
```bash
$ gosha
Hi user!
That's Gosha!
gosha>>
```

2) Использование в скриптах с использованием шебанга:
```bash
(your-script.sh)
#!/usr/bin/gosha

print("Hello world!")
...
```

### Мотивация

В первую очередь, этот проект был создан с целью избавиться от проблемы устаревшего и неудобного синтаксиса Bash-скриптов
