# KPAR — Design Spec

## Overview

KPAR (corruptela de "capar" — cortar, encurtar) is a CLI tool for optimizing image file sizes by converting between modern formats (WEBP, AVIF) and keeping only the smallest result.

## Usage

```
kpar [arquivo]          # process a specific file
kpar                    # scan current directory, show interactive picker
kpar --quality 80       # set encoding quality (applies to both modes)
```

## Supported Input Formats

- JPEG (`.jpg`, `.jpeg`)
- PNG (`.png`)
- WEBP (`.webp`)

## Conversion Strategies

### WEBP input

1. Convert to AVIF
2. Compare sizes: AVIF vs original WEBP
3. If AVIF >= WEBP → delete AVIF, report "no gain"
4. If AVIF < WEBP → keep AVIF as `basename.avif`, report savings

### JPG/PNG input

1. Convert to WEBP
2. Convert to AVIF
3. Compare all three: original, WEBP, AVIF
4. Keep the original always (current behavior — see Future Work)
5. Between WEBP and AVIF, keep only the smaller one; discard the other
6. If both conversions are larger than the original → discard both, report "no gain"

## Encoder Priority

For each target format, KPAR uses the best available encoder:

### WEBP encoding
1. `cwebp` (dedicated encoder, best quality/compression) — preferred
2. `magick` (ImageMagick 7+) — fallback

### AVIF encoding
1. `avifenc` (dedicated encoder, best quality/compression) — preferred
2. `magick` (ImageMagick 7+) — fallback

### Startup behavior
- KPAR checks which encoders are available in PATH
- For each format, selects the best available encoder
- Reports which encoder is being used (e.g., "using cwebp for WEBP, magick for AVIF")
- Errors only if no encoder is available for a required format
- Provides clear install instructions when missing: `brew install webp libavif` / `apt install webp libavif-bin`

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--quality`, `-q` | encoder default | Encoding quality (0-100) |

## Interactive Picker

When no file argument is provided:

1. Scanner finds all JPG/PNG/WEBP files in the current directory
2. Picker displays an interactive menu:
   - **Primary:** fzf-style fuzzy search (via Bubbletea)
   - **Fallback:** numbered list if terminal doesn't support interactive mode
3. User selects a file, processing continues

## Output Format

After processing, KPAR displays results with colors (via Lipgloss):

```
foto.jpg (2.4 MB)
  ├─ webp: 1.1 MB (54% do original) ✗ descartado
  ├─ avif:  820 KB (33% do original) ✓ salvo
  └─ ganho: 1.6 MB economizados
```

When there's no gain:

```
foto.webp (45 KB)
  ├─ avif: 52 KB (115% do original) ✗ descartado
  └─ sem ganho
```

Color scheme:
- Green: saved/kept file
- Red/gray: discarded file
- Highlighted: total savings

## Project Structure

```
kpar/
├── cmd/kpar/
│   └── main.go           # entry point, arg parsing
├── internal/
│   ├── converter/
│   │   ├── converter.go  # interface + keep/discard logic
│   │   ├── webp.go       # cwebp / magick webp encoding
│   │   └── avif.go       # avifenc / magick avif encoding
│   ├── picker/
│   │   └── picker.go     # interactive menu (Bubbletea + numbered fallback)
│   └── scanner/
│       └── scanner.go    # find JPG/PNG/WEBP in directory
├── go.mod
└── go.sum
```

## Dependencies

### Go libraries
- [Bubbletea](https://github.com/charmbracelet/bubbletea) — interactive TUI (picker)
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling/colors
- CLI parsing library (Cobra or similar)

### External tools (at least one required)
- `cwebp` — WEBP encoder (`brew install webp` / `apt install webp`)
- `avifenc` — AVIF encoder (`brew install libavif` / `apt install libavif-bin`)
- `magick` — ImageMagick 7+ (fallback for both formats)

## Future Work

Items deliberately deferred:

1. **Pipeline support** — currently standalone only; design for composability later
2. **Original file handling** — currently keeps originals; later evaluate deleting or moving to `originals/`
3. **Batch processing** — currently one file at a time; later add "process all" option
