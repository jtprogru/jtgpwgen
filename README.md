# jtgpwgen

Генератор паролей с настраиваемыми классами символов.

Извлечён из [gch](https://github.com/jtprogru/gch) в рамках распиливания монолита.

## Установка

### Homebrew (cask)

```bash
brew tap jtprogru/tap
brew install --cask jtgpwgen
```

### go install

```bash
go install github.com/jtprogru/jtgpwgen@latest
```

## Использование

```bash
# дефолт: 24 символа, латиница + цифры + '@'
jtgpwgen

# нужная длина
jtgpwgen -l 32

# дополнительные спецсимволы (добавляются к '@')
jtgpwgen -s '#$%^&*'

# без спецсимволов
jtgpwgen --no-special

# без цифр
jtgpwgen --no-digits

# запоминаемый пароль (слоги CVC через дефис + 2 цифры + '@')
jtgpwgen -m -l 24
```

## Флаги

| Флаг | Описание |
|------|----------|
| `-l`, `--length` | Длина пароля (по умолчанию 24) |
| `-s`, `--special` | Дополнительные спецсимволы поверх дефолтного `@` |
| `-d`, `--digits` | Явно включить цифры |
| `--no-special` | Отключить спецсимволы |
| `--no-digits` | Отключить цифры |
| `-m`, `--memo` | Запоминаемый пароль |

`-d`/`--no-digits` и `-s`/`--no-special` — взаимоисключающие пары.

Источник энтропии — `crypto/rand`.
