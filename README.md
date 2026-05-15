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
| `--no-letters` | Отключить латиницу |
| `-m`, `--memo` | Запоминаемый пароль |
| `--clip` | Скопировать пароль в буфер обмена вместо вывода в stdout |
| `--clip-ttl` | TTL для очистки буфера при `--clip` (по умолчанию `30s`) |

`-d`/`--no-digits` и `-s`/`--no-special` — взаимоисключающие пары.

`-m`/`--memo` нельзя комбинировать с класс-флагами (`-d`, `-s`, `--no-digits`, `--no-special`, `--no-letters`): memo-режим всегда содержит цифры и `@` по фиксированному шаблону.

При `--clip` пароль не печатается в stdout — он попадает прямо в системный буфер обмена (`pbcopy` / `wl-copy` / `xclip` / `xsel`) и автоматически вытирается через `--clip-ttl`. Процесс блокируется до истечения TTL.

Источник энтропии — `crypto/rand`.

## Верификация релизов

Релизы подписываются [cosign](https://github.com/sigstore/cosign) в keyless-режиме (Sigstore OIDC), для каждого архива публикуется SBOM.

```bash
# Скачать релиз (пример):
gh release download vX.Y.Z --pattern 'jtgpwgen_*'
gh release download vX.Y.Z --pattern 'checksums.txt*'

# Проверить подпись checksums.txt:
cosign verify-blob \
  --certificate checksums.txt.pem \
  --signature checksums.txt.sig \
  --certificate-identity-regexp 'https://github\.com/jtprogru/jtgpwgen/.+' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  checksums.txt

# Проверить, что архив соответствует подписанной сумме:
sha256sum -c checksums.txt --ignore-missing
```

SBOM (`*.sbom.json`, формат SPDX-JSON) публикуется рядом с архивами.

## Энтропия memorable-режима

Memo-пароль строится из слогов CVC, разделённых дефисом, и завершается `NN@` (две цифры + `@`).
Алфавит: 18 согласных (без `q`, `x`) и 6 гласных. Сила одного слога — `log2(18·6·18) ≈ 10.92` бит, суффикса — `log2(100) ≈ 6.64` бит. Например:

| `--length` | Слогов | Энтропия |
|------------|--------|----------|
| 24 | 6 | ~72 бит |
| 32 | 8 | ~94 бит |
| 40 | 10 | ~116 бит |

Минимально допустимая энтропия memo-режима — **64 бита** (~`--length 23`). При меньшем значении генератор вернёт ошибку `ErrMemoEntropyTooLow`.

Для критичных аккаунтов используйте либо `--length 40+` в memo-режиме, либо обычный режим (24 символа из 73-символьного алфавита ≈ 149 бит).
