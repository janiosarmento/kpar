# kpar

Image optimizer CLI — converts images between modern formats (WEBP, AVIF) and keeps only the smallest result.

The name is a corruption of the Portuguese word "capar" (to cut, to trim).

## How it works

- **WEBP input** → converts to AVIF; keeps it only if smaller than the original
- **JPG/PNG input** → converts to both WEBP and AVIF; keeps only the smallest

When a smaller conversion is found, the original file is **deleted** by default. Use `--keep` to preserve it.

If no conversion produces a smaller file, nothing changes.

## Usage

```bash
# Process a specific file
kpar photo.jpg

# Interactive picker — scans current directory for images
kpar

# Set encoding quality (0-100)
kpar --quality 80 photo.png

# Keep original file after optimization
kpar --keep photo.jpg
```

## Requirements

At least one of the following must be installed:

| Tool | Install (macOS) | Install (Debian/Ubuntu) |
|------|----------------|------------------------|
| `cwebp` / `dwebp` | `brew install webp` | `apt install webp` |
| `avifenc` | `brew install libavif` | `apt install libavif-bin` |
| `magick` (ImageMagick 7+) | `brew install imagemagick` | `apt install imagemagick` |

Dedicated encoders (`cwebp`, `avifenc`) are preferred when available. ImageMagick is used as a fallback.

## Install

### From source

```bash
git clone https://github.com/janiosarmento/kpar.git
cd kpar
./install.sh
```

This builds the binary and symlinks it to `~/.local/bin/kpar`. Make sure `~/.local/bin` is in your `PATH`.

### Go install

```bash
go install github.com/janiosarmento/kpar/cmd/kpar@latest
```

## License

MIT
