# Guessr

**Guessr** is a small Go tool that reads the first few bytes of a file (its "magic bytes") and attempts to guess the file format with a confidence percentage. It’s useful for quickly identifying file types without relying on file extensions.

---

## Features

- Reads a file’s magic bytes
- Compares them with known file signatures
- Provides a percentage-based confidence for each format
- Simple CLI usage

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/username/guessr.git
cd guessr
```

2. Build the Go binary:

```bash
go build -o guessr
```

or just run
```
go run main.go <filename>
```


Usage:

```bash
./guessr <file_path>
```

Example:

```bash
./guessr sample.png
```

Example output:

```bash
Read bytes: 89 50 4E 47 0D 0A 1A 0A

Format likelihoods:
  JPEG   : 0% match
  PNG    : 100% match
  GIF    : 0% match
  PDF    : 0% match
  ZIP    : 0% match

Most likely format: PNG (100%)
```
## Planned features

- [ ] External signature list
- [ ] Colored table-like output
- [ ] Verbose mode (display matching and different bytes)
- [ ] Option to display only the top N guesses
- [ ] Guess by extension vs magic bytes