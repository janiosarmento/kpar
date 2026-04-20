# KPAR Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a CLI image optimizer that converts images between WEBP/AVIF and keeps only the smallest result.

**Architecture:** Subprocess-based — Go orchestrates external encoders (`cwebp`, `avifenc`, `magick`) to convert images, compares file sizes, and keeps/discards results. Interactive picker via Bubbletea with numbered fallback.

**Tech Stack:** Go 1.26, Bubbletea/Bubbles/Lipgloss (Charm ecosystem), Cobra (CLI), `os/exec` for subprocesses.

---

## File Structure

```
kpar/
├── cmd/kpar/
│   └── main.go              # Cobra root command, arg parsing, orchestration
├── internal/
│   ├── encoder/
│   │   ├── encoder.go       # Encoder interface, Registry, detection logic
│   │   ├── webp.go          # CwebpEncoder + MagickWebpEncoder
│   │   ├── avif.go          # AvifencEncoder + MagickAvifEncoder
│   │   └── encoder_test.go  # Tests for encoder detection and execution
│   ├── converter/
│   │   ├── converter.go     # Conversion strategies (WEBP input, JPG/PNG input)
│   │   └── converter_test.go
│   ├── picker/
│   │   ├── picker.go        # Bubbletea model + numbered fallback
│   │   └── picker_test.go
│   ├── scanner/
│   │   ├── scanner.go       # Find image files in directory
│   │   └── scanner_test.go
│   └── output/
│       ├── output.go        # Format and print results with Lipgloss
│       └── output_test.go
├── go.mod
└── go.sum
```

---

### Task 1: Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `cmd/kpar/main.go`

- [ ] **Step 1: Initialize Go module**

```bash
cd /Users/janiosarmento/projects/kpar
go mod init github.com/janiosarmento/kpar
```

- [ ] **Step 2: Create minimal main.go**

Create `cmd/kpar/main.go`:

```go
package main

import "fmt"

func main() {
	fmt.Println("kpar")
}
```

- [ ] **Step 3: Verify it compiles and runs**

Run: `go run ./cmd/kpar/`
Expected: prints `kpar`

- [ ] **Step 4: Commit**

```bash
git add go.mod cmd/
git commit -m "Scaffold project with Go module and minimal main"
```

---

### Task 2: Scanner — Find Image Files

**Files:**
- Create: `internal/scanner/scanner.go`
- Create: `internal/scanner/scanner_test.go`

- [ ] **Step 1: Write failing tests for scanner**

Create `internal/scanner/scanner_test.go`:

```go
package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/janiosarmento/kpar/internal/scanner"
)

func TestScanFindsImages(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	files := []string{"photo.jpg", "image.jpeg", "pic.png", "banner.webp", "readme.txt", "data.csv"}
	for _, f := range files {
		os.WriteFile(filepath.Join(dir, f), []byte("fake"), 0644)
	}

	result, err := scanner.Scan(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"banner.webp", "image.jpeg", "photo.jpg", "pic.png"}
	if len(result) != len(expected) {
		t.Fatalf("got %d files, want %d: %v", len(result), len(expected), result)
	}
	for i, name := range expected {
		if filepath.Base(result[i]) != name {
			t.Errorf("result[%d] = %q, want %q", i, filepath.Base(result[i]), name)
		}
	}
}

func TestScanEmptyDir(t *testing.T) {
	dir := t.TempDir()

	result, err := scanner.Scan(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestScanCaseInsensitive(t *testing.T) {
	dir := t.TempDir()

	files := []string{"Photo.JPG", "IMAGE.PNG", "banner.WEBP", "Pic.Jpeg"}
	for _, f := range files {
		os.WriteFile(filepath.Join(dir, f), []byte("fake"), 0644)
	}

	result, err := scanner.Scan(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 4 {
		t.Fatalf("got %d files, want 4: %v", len(result), result)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/janiosarmento/projects/kpar && go test ./internal/scanner/`
Expected: compilation error — package doesn't exist yet

- [ ] **Step 3: Implement scanner**

Create `internal/scanner/scanner.go`:

```go
package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// Scan returns sorted absolute paths of image files in dir.
func Scan(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if imageExtensions[ext] {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}

	sort.Strings(files)
	return files, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/scanner/ -v`
Expected: all 3 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/scanner/
git commit -m "Add scanner to find image files in directory"
```

---

### Task 3: Encoder — Interface and Detection

**Files:**
- Create: `internal/encoder/encoder.go`
- Create: `internal/encoder/encoder_test.go`

- [ ] **Step 1: Write failing tests for encoder detection**

Create `internal/encoder/encoder_test.go`:

```go
package encoder_test

import (
	"testing"

	"github.com/janiosarmento/kpar/internal/encoder"
)

func TestDetectEncoders(t *testing.T) {
	registry := encoder.Detect()

	// At least one encoder should be found on a dev machine
	// (we have cwebp, avifenc, and magick installed)
	if registry.WebpEncoder == nil && registry.AvifEncoder == nil {
		t.Fatal("no encoders detected — need at least cwebp/avifenc/magick in PATH")
	}
}

func TestRegistryNames(t *testing.T) {
	registry := encoder.Detect()

	if registry.WebpEncoder != nil {
		name := registry.WebpEncoder.Name()
		if name == "" {
			t.Error("WebpEncoder.Name() returned empty string")
		}
	}

	if registry.AvifEncoder != nil {
		name := registry.AvifEncoder.Name()
		if name == "" {
			t.Error("AvifEncoder.Name() returned empty string")
		}
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/encoder/`
Expected: compilation error — package doesn't exist

- [ ] **Step 3: Implement encoder interface and detection**

Create `internal/encoder/encoder.go`:

```go
package encoder

import "os/exec"

// Encoder converts a source image to a target format.
type Encoder interface {
	// Name returns the encoder name (e.g., "cwebp", "magick").
	Name() string

	// Encode converts src to the target format, writing to dst.
	// quality is 0-100; if <0, use encoder default.
	Encode(src, dst string, quality int) error
}

// Registry holds the best available encoder for each format.
type Registry struct {
	WebpEncoder Encoder
	AvifEncoder Encoder
}

// Detect finds the best available encoders in PATH.
// Priority: dedicated encoders > ImageMagick fallback.
func Detect() Registry {
	var r Registry

	if commandExists("cwebp") {
		r.WebpEncoder = CwebpEncoder{}
	} else if commandExists("magick") {
		r.WebpEncoder = MagickWebpEncoder{}
	}

	if commandExists("avifenc") {
		r.AvifEncoder = AvifencEncoder{}
	} else if commandExists("magick") {
		r.AvifEncoder = MagickAvifEncoder{}
	}

	return r
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/encoder/ -v`
Expected: FAIL — CwebpEncoder etc. not defined yet. That's expected; we implement them in the next tasks.

Note: Tests will pass once Task 4 and Task 5 are complete. For now, create stub types so tests compile:

Add to the bottom of `encoder.go`:

```go
// Stubs — implemented in webp.go and avif.go
type CwebpEncoder struct{}
func (CwebpEncoder) Name() string                         { return "cwebp" }
func (CwebpEncoder) Encode(src, dst string, q int) error  { return nil }

type MagickWebpEncoder struct{}
func (MagickWebpEncoder) Name() string                         { return "magick (webp)" }
func (MagickWebpEncoder) Encode(src, dst string, q int) error  { return nil }

type AvifencEncoder struct{}
func (AvifencEncoder) Name() string                         { return "avifenc" }
func (AvifencEncoder) Encode(src, dst string, q int) error  { return nil }

type MagickAvifEncoder struct{}
func (MagickAvifEncoder) Name() string                         { return "magick (avif)" }
func (MagickAvifEncoder) Encode(src, dst string, q int) error  { return nil }
```

Run: `go test ./internal/encoder/ -v`
Expected: all tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/encoder/
git commit -m "Add encoder interface, registry, and detection logic"
```

---

### Task 4: WEBP Encoder Implementations

**Files:**
- Create: `internal/encoder/webp.go`
- Modify: `internal/encoder/encoder.go` (remove CwebpEncoder and MagickWebpEncoder stubs)
- Modify: `internal/encoder/encoder_test.go` (add encoding tests)

- [ ] **Step 1: Add encoding test**

Append to `internal/encoder/encoder_test.go`:

```go
func TestCwebpEncode(t *testing.T) {
	if _, err := exec.LookPath("cwebp"); err != nil {
		t.Skip("cwebp not installed")
	}

	src := createTestPNG(t)
	dst := filepath.Join(t.TempDir(), "out.webp")

	enc := encoder.CwebpEncoder{}
	if err := enc.Encode(src, dst, -1); err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

func TestMagickWebpEncode(t *testing.T) {
	if _, err := exec.LookPath("magick"); err != nil {
		t.Skip("magick not installed")
	}

	src := createTestPNG(t)
	dst := filepath.Join(t.TempDir(), "out.webp")

	enc := encoder.MagickWebpEncoder{}
	if err := enc.Encode(src, dst, -1); err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

// createTestPNG creates a minimal valid PNG file and returns its path.
func createTestPNG(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.png")

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for i := range img.Pix {
		img.Pix[i] = 128
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	return path
}
```

Add these imports to the test file's import block:

```go
"image"
"image/png"
"os/exec"
"path/filepath"
```

- [ ] **Step 2: Run tests to verify new tests fail**

Run: `go test ./internal/encoder/ -v -run TestCwebp`
Expected: FAIL — stub Encode does nothing useful (but actually stubs return nil, so we need the real impl to validate output). The test should fail because the stub doesn't actually create an output file.

- [ ] **Step 3: Implement real WEBP encoders**

Create `internal/encoder/webp.go`:

```go
package encoder

import (
	"fmt"
	"os/exec"
)

type CwebpEncoder struct{}

func (CwebpEncoder) Name() string { return "cwebp" }

func (CwebpEncoder) Encode(src, dst string, quality int) error {
	args := []string{src, "-o", dst}
	if quality >= 0 {
		args = append([]string{"-q", fmt.Sprintf("%d", quality)}, args...)
	}
	cmd := exec.Command("cwebp", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("cwebp: %w\n%s", err, out)
	}
	return nil
}

type MagickWebpEncoder struct{}

func (MagickWebpEncoder) Name() string { return "magick (webp)" }

func (MagickWebpEncoder) Encode(src, dst string, quality int) error {
	args := []string{"convert", src}
	if quality >= 0 {
		args = append(args, "-quality", fmt.Sprintf("%d", quality))
	}
	args = append(args, dst)
	cmd := exec.Command("magick", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("magick: %w\n%s", err, out)
	}
	return nil
}
```

Remove the `CwebpEncoder` and `MagickWebpEncoder` stubs from `encoder.go`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/encoder/ -v`
Expected: PASS for TestCwebpEncode and TestMagickWebpEncode

- [ ] **Step 5: Commit**

```bash
git add internal/encoder/
git commit -m "Implement cwebp and magick WEBP encoders"
```

---

### Task 5: AVIF Encoder Implementations

**Files:**
- Create: `internal/encoder/avif.go`
- Modify: `internal/encoder/encoder.go` (remove AvifencEncoder and MagickAvifEncoder stubs)
- Modify: `internal/encoder/encoder_test.go` (add AVIF tests)

- [ ] **Step 1: Add encoding tests**

Append to `internal/encoder/encoder_test.go`:

```go
func TestAvifencEncode(t *testing.T) {
	if _, err := exec.LookPath("avifenc"); err != nil {
		t.Skip("avifenc not installed")
	}

	src := createTestPNG(t)
	dst := filepath.Join(t.TempDir(), "out.avif")

	enc := encoder.AvifencEncoder{}
	if err := enc.Encode(src, dst, -1); err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

func TestMagickAvifEncode(t *testing.T) {
	if _, err := exec.LookPath("magick"); err != nil {
		t.Skip("magick not installed")
	}

	src := createTestPNG(t)
	dst := filepath.Join(t.TempDir(), "out.avif")

	enc := encoder.MagickAvifEncoder{}
	if err := enc.Encode(src, dst, -1); err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}
```

- [ ] **Step 2: Run tests to verify new tests fail**

Run: `go test ./internal/encoder/ -v -run TestAvifenc`
Expected: FAIL — stub doesn't produce output file

- [ ] **Step 3: Implement real AVIF encoders**

Create `internal/encoder/avif.go`:

```go
package encoder

import (
	"fmt"
	"os/exec"
)

type AvifencEncoder struct{}

func (AvifencEncoder) Name() string { return "avifenc" }

func (AvifencEncoder) Encode(src, dst string, quality int) error {
	args := []string{}
	if quality >= 0 {
		// avifenc uses --min/--max for quality range; for simplicity map quality to both
		args = append(args, "--min", fmt.Sprintf("%d", quality), "--max", fmt.Sprintf("%d", quality))
	}
	args = append(args, src, dst)
	cmd := exec.Command("avifenc", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("avifenc: %w\n%s", err, out)
	}
	return nil
}

type MagickAvifEncoder struct{}

func (MagickAvifEncoder) Name() string { return "magick (avif)" }

func (MagickAvifEncoder) Encode(src, dst string, quality int) error {
	args := []string{"convert", src}
	if quality >= 0 {
		args = append(args, "-quality", fmt.Sprintf("%d", quality))
	}
	args = append(args, dst)
	cmd := exec.Command("magick", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("magick: %w\n%s", err, out)
	}
	return nil
}
```

Remove the `AvifencEncoder` and `MagickAvifEncoder` stubs from `encoder.go`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/encoder/ -v`
Expected: all encoder tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/encoder/
git commit -m "Implement avifenc and magick AVIF encoders"
```

---

### Task 6: Converter — Conversion Strategies

**Files:**
- Create: `internal/converter/converter.go`
- Create: `internal/converter/converter_test.go`

- [ ] **Step 1: Write failing tests**

Create `internal/converter/converter_test.go`:

```go
package converter_test

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/janiosarmento/kpar/internal/converter"
	"github.com/janiosarmento/kpar/internal/encoder"
)

func createTestImage(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)

	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for i := range img.Pix {
		img.Pix[i] = 128
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestConvertPNG(t *testing.T) {
	registry := encoder.Detect()
	if registry.WebpEncoder == nil || registry.AvifEncoder == nil {
		t.Skip("need both webp and avif encoders")
	}

	dir := t.TempDir()
	src := createTestImage(t, dir, "test.png")

	result, err := converter.Convert(src, registry, -1)
	if err != nil {
		t.Fatalf("convert failed: %v", err)
	}

	if result.OriginalPath != src {
		t.Errorf("OriginalPath = %q, want %q", result.OriginalPath, src)
	}
	if result.OriginalSize <= 0 {
		t.Error("OriginalSize should be > 0")
	}
	// At least one conversion should have been attempted
	if len(result.Conversions) == 0 {
		t.Error("expected at least one conversion attempt")
	}
}

func TestConvertWEBP(t *testing.T) {
	registry := encoder.Detect()
	if registry.WebpEncoder == nil || registry.AvifEncoder == nil {
		t.Skip("need both webp and avif encoders")
	}

	// First create a WEBP from a PNG
	dir := t.TempDir()
	pngPath := createTestImage(t, dir, "test.png")
	webpPath := filepath.Join(dir, "test.webp")
	if err := registry.WebpEncoder.Encode(pngPath, webpPath, -1); err != nil {
		t.Fatalf("setup: creating webp: %v", err)
	}

	result, err := converter.Convert(webpPath, registry, -1)
	if err != nil {
		t.Fatalf("convert failed: %v", err)
	}

	// WEBP input should only attempt AVIF
	for _, c := range result.Conversions {
		if c.Format == "webp" {
			t.Error("WEBP input should not convert to WEBP")
		}
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/converter/`
Expected: compilation error — package doesn't exist

- [ ] **Step 3: Implement converter**

Create `internal/converter/converter.go`:

```go
package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/janiosarmento/kpar/internal/encoder"
)

// Conversion holds the result of a single format conversion.
type Conversion struct {
	Format  string // "webp" or "avif"
	Path    string // path to converted file
	Size    int64
	Kept    bool   // true if this was the winner
	Encoder string // which encoder was used
}

// Result holds the full result of processing one image.
type Result struct {
	OriginalPath string
	OriginalSize int64
	Conversions  []Conversion
	BestPath     string // path to the best converted file, or "" if no gain
	Saved        int64  // bytes saved (0 if no gain)
}

// Convert processes a single image file according to its type.
func Convert(src string, reg encoder.Registry, quality int) (Result, error) {
	info, err := os.Stat(src)
	if err != nil {
		return Result{}, fmt.Errorf("stat %s: %w", src, err)
	}

	result := Result{
		OriginalPath: src,
		OriginalSize: info.Size(),
	}

	ext := strings.ToLower(filepath.Ext(src))
	base := strings.TrimSuffix(src, filepath.Ext(src))

	switch ext {
	case ".webp":
		err = convertWebp(base, src, reg, quality, &result)
	case ".jpg", ".jpeg", ".png":
		err = convertJpgPng(base, src, reg, quality, &result)
	default:
		return result, fmt.Errorf("unsupported format: %s", ext)
	}

	if err != nil {
		return result, err
	}

	pickBest(&result)
	return result, nil
}

func convertWebp(base, src string, reg encoder.Registry, quality int, result *Result) error {
	if reg.AvifEncoder == nil {
		return fmt.Errorf("no AVIF encoder available")
	}

	dst := base + ".avif"
	if err := reg.AvifEncoder.Encode(src, dst, quality); err != nil {
		return err
	}

	info, err := os.Stat(dst)
	if err != nil {
		return err
	}

	result.Conversions = append(result.Conversions, Conversion{
		Format:  "avif",
		Path:    dst,
		Size:    info.Size(),
		Encoder: reg.AvifEncoder.Name(),
	})
	return nil
}

func convertJpgPng(base, src string, reg encoder.Registry, quality int, result *Result) error {
	if reg.WebpEncoder != nil {
		dst := base + ".webp"
		if err := reg.WebpEncoder.Encode(src, dst, quality); err != nil {
			return err
		}
		info, err := os.Stat(dst)
		if err != nil {
			return err
		}
		result.Conversions = append(result.Conversions, Conversion{
			Format:  "webp",
			Path:    dst,
			Size:    info.Size(),
			Encoder: reg.WebpEncoder.Name(),
		})
	}

	if reg.AvifEncoder != nil {
		dst := base + ".avif"
		if err := reg.AvifEncoder.Encode(src, dst, quality); err != nil {
			return err
		}
		info, err := os.Stat(dst)
		if err != nil {
			return err
		}
		result.Conversions = append(result.Conversions, Conversion{
			Format:  "avif",
			Path:    dst,
			Size:    info.Size(),
			Encoder: reg.AvifEncoder.Name(),
		})
	}

	if len(result.Conversions) == 0 {
		return fmt.Errorf("no encoders available")
	}
	return nil
}

// pickBest selects the smallest conversion (if smaller than original),
// marks it as kept, and removes the rest.
func pickBest(result *Result) {
	var bestIdx = -1
	var bestSize int64 = result.OriginalSize

	for i, c := range result.Conversions {
		if c.Size < bestSize {
			bestSize = c.Size
			bestIdx = i
		}
	}

	for i := range result.Conversions {
		if i == bestIdx {
			result.Conversions[i].Kept = true
			result.BestPath = result.Conversions[i].Path
			result.Saved = result.OriginalSize - result.Conversions[i].Size
		} else {
			// Remove the discarded file
			os.Remove(result.Conversions[i].Path)
		}
	}
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/converter/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/converter/
git commit -m "Add converter with WEBP and JPG/PNG conversion strategies"
```

---

### Task 7: Output — Formatted Results

**Files:**
- Create: `internal/output/output.go`
- Create: `internal/output/output_test.go`

- [ ] **Step 1: Write failing tests**

Create `internal/output/output_test.go`:

```go
package output_test

import (
	"testing"

	"github.com/janiosarmento/kpar/internal/converter"
	"github.com/janiosarmento/kpar/internal/output"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{2621440, "2.5 MB"},
	}

	for _, tt := range tests {
		got := output.FormatSize(tt.bytes)
		if got != tt.want {
			t.Errorf("FormatSize(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

func TestRenderResult(t *testing.T) {
	result := converter.Result{
		OriginalPath: "/tmp/foto.jpg",
		OriginalSize: 2457600,
		Conversions: []converter.Conversion{
			{Format: "webp", Path: "/tmp/foto.webp", Size: 1126400, Kept: false, Encoder: "cwebp"},
			{Format: "avif", Path: "/tmp/foto.avif", Size: 839680, Kept: true, Encoder: "avifenc"},
		},
		BestPath: "/tmp/foto.avif",
		Saved:    1617920,
	}

	rendered := output.Render(result)
	if rendered == "" {
		t.Error("Render returned empty string")
	}
}

func TestRenderNoGain(t *testing.T) {
	result := converter.Result{
		OriginalPath: "/tmp/foto.webp",
		OriginalSize: 45000,
		Conversions: []converter.Conversion{
			{Format: "avif", Path: "/tmp/foto.avif", Size: 52000, Kept: false, Encoder: "avifenc"},
		},
		Saved: 0,
	}

	rendered := output.Render(result)
	if rendered == "" {
		t.Error("Render returned empty string")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/output/`
Expected: compilation error — package doesn't exist

- [ ] **Step 3: Install Lipgloss dependency**

```bash
cd /Users/janiosarmento/projects/kpar && go get github.com/charmbracelet/lipgloss
```

- [ ] **Step 4: Implement output formatter**

Create `internal/output/output.go`:

```go
package output

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/janiosarmento/kpar/internal/converter"
)

var (
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	red    = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	gray   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	bold   = lipgloss.NewStyle().Bold(true)
)

// FormatSize formats bytes into a human-readable string.
func FormatSize(bytes int64) string {
	switch {
	case bytes >= 1048576:
		return fmt.Sprintf("%.1f MB", float64(bytes)/1048576)
	case bytes >= 1024:
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// Render formats a conversion result for terminal display.
func Render(r converter.Result) string {
	var b strings.Builder

	name := filepath.Base(r.OriginalPath)
	fmt.Fprintf(&b, "%s (%s)\n", bold.Render(name), FormatSize(r.OriginalSize))

	last := len(r.Conversions) - 1
	for i, c := range r.Conversions {
		connector := "├─"
		if i == last && r.Saved == 0 {
			connector = "└─"
		}

		pct := float64(c.Size) / float64(r.OriginalSize) * 100
		sizeStr := fmt.Sprintf("%s: %s (%.*f%% do original)",
			c.Format, FormatSize(c.Size), 0, pct)

		if c.Kept {
			fmt.Fprintf(&b, "  %s %s %s\n", connector, green.Render(sizeStr), green.Render("✓ salvo"))
		} else {
			fmt.Fprintf(&b, "  %s %s %s\n", connector, gray.Render(sizeStr), red.Render("✗ descartado"))
		}
	}

	if r.Saved > 0 {
		savings := fmt.Sprintf("ganho: %s economizados", FormatSize(r.Saved))
		fmt.Fprintf(&b, "  └─ %s\n", green.Render(savings))
	} else {
		fmt.Fprintf(&b, "  └─ %s\n", gray.Render("sem ganho"))
	}

	return b.String()
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/output/ -v`
Expected: all tests PASS

- [ ] **Step 6: Commit**

```bash
git add internal/output/ go.mod go.sum
git commit -m "Add output formatter with colored results display"
```

---

### Task 8: Picker — Interactive File Selection

**Files:**
- Create: `internal/picker/picker.go`
- Create: `internal/picker/picker_test.go`

- [ ] **Step 1: Write failing test for fallback picker**

Create `internal/picker/picker_test.go`:

```go
package picker_test

import (
	"strings"
	"testing"

	"github.com/janiosarmento/kpar/internal/picker"
)

func TestFormatChoices(t *testing.T) {
	files := []string{
		"/home/user/photos/beach.jpg",
		"/home/user/photos/sunset.png",
		"/home/user/photos/logo.webp",
	}

	formatted := picker.FormatChoices(files)
	if len(formatted) != 3 {
		t.Fatalf("got %d choices, want 3", len(formatted))
	}

	if !strings.Contains(formatted[0], "beach.jpg") {
		t.Errorf("formatted[0] = %q, want to contain 'beach.jpg'", formatted[0])
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/picker/`
Expected: compilation error — package doesn't exist

- [ ] **Step 3: Install Bubbletea dependencies**

```bash
cd /Users/janiosarmento/projects/kpar
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
```

- [ ] **Step 4: Implement picker**

Create `internal/picker/picker.go`:

```go
package picker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FormatChoices returns display strings for file paths (basename only).
func FormatChoices(files []string) []string {
	choices := make([]string, len(files))
	for i, f := range files {
		choices[i] = filepath.Base(f)
	}
	return choices
}

// item implements list.Item for Bubbletea list.
type item struct {
	path string
	name string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.name }

// Pick shows an interactive picker and returns the selected file path.
// Falls back to numbered list if terminal is not interactive.
func Pick(files []string) (string, error) {
	if len(files) == 0 {
		return "", fmt.Errorf("no image files found")
	}

	if !isInteractive() {
		return pickNumbered(files)
	}

	return pickBubbletea(files)
}

func isInteractive() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func pickNumbered(files []string) (string, error) {
	names := FormatChoices(files)
	fmt.Println("Select an image to optimize:")
	for i, name := range names {
		fmt.Printf("  %d) %s\n", i+1, name)
	}

	fmt.Print("\n> ")
	var choice int
	if _, err := fmt.Scan(&choice); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	if choice < 1 || choice > len(files) {
		return "", fmt.Errorf("choice %d out of range (1-%d)", choice, len(files))
	}

	return files[choice-1], nil
}

type model struct {
	list     list.Model
	files    []string
	selected string
	quit     bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = i.path
			}
			return m, tea.Quit
		case "q", "ctrl+c", "esc":
			m.quit = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quit {
		return ""
	}
	return "\n" + m.list.View()
}

func pickBubbletea(files []string) (string, error) {
	items := make([]list.Item, len(files))
	for i, f := range files {
		items[i] = item{path: f, name: filepath.Base(f)}
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("2")).
		BorderLeftForeground(lipgloss.Color("2"))

	l := list.New(items, delegate, 40, 14)
	l.Title = "Select an image to optimize"
	l.SetShowStatusBar(false)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Padding(0, 1)

	// Suppress the filter prompt prefix being too wide
	filterPrompt := "Filter: "
	_ = filterPrompt
	_ = strings.TrimSpace

	m := model{list: l, files: files}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("picker: %w", err)
	}

	final := finalModel.(model)
	if final.quit || final.selected == "" {
		return "", fmt.Errorf("cancelled")
	}

	return final.selected, nil
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/picker/ -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/picker/ go.mod go.sum
git commit -m "Add interactive picker with Bubbletea and numbered fallback"
```

---

### Task 9: CLI Wiring — Cobra + Main

**Files:**
- Modify: `cmd/kpar/main.go`

- [ ] **Step 1: Install Cobra**

```bash
cd /Users/janiosarmento/projects/kpar && go get github.com/spf13/cobra
```

- [ ] **Step 2: Implement full CLI**

Replace `cmd/kpar/main.go` with:

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/janiosarmento/kpar/internal/converter"
	"github.com/janiosarmento/kpar/internal/encoder"
	"github.com/janiosarmento/kpar/internal/output"
	"github.com/janiosarmento/kpar/internal/picker"
	"github.com/janiosarmento/kpar/internal/scanner"
)

var quality int

var rootCmd = &cobra.Command{
	Use:   "kpar [file]",
	Short: "Image optimizer — converts to WEBP/AVIF and keeps the smallest",
	Long:  "KPAR converts images between modern formats (WEBP, AVIF) and keeps only the smallest result.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  run,
}

func init() {
	rootCmd.Flags().IntVarP(&quality, "quality", "q", -1, "Encoding quality (0-100, default: encoder default)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Detect encoders
	registry := encoder.Detect()
	reportEncoders(registry)

	if registry.WebpEncoder == nil && registry.AvifEncoder == nil {
		fmt.Fprintln(os.Stderr, "No encoders found. Install at least one:")
		fmt.Fprintln(os.Stderr, "  brew install webp libavif    (macOS)")
		fmt.Fprintln(os.Stderr, "  apt install webp libavif-bin (Debian/Ubuntu)")
		fmt.Fprintln(os.Stderr, "  Or install ImageMagick 7+: brew install imagemagick")
		return fmt.Errorf("no encoders available")
	}

	var filePath string

	if len(args) == 1 {
		// Direct file mode
		abs, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		filePath = abs
	} else {
		// Picker mode
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		files, err := scanner.Scan(dir)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			fmt.Println("No image files found in current directory.")
			return nil
		}

		selected, err := picker.Pick(files)
		if err != nil {
			return err
		}
		filePath = selected
	}

	// Convert
	result, err := converter.Convert(filePath, registry, quality)
	if err != nil {
		return err
	}

	// Output
	fmt.Print(output.Render(result))
	return nil
}

func reportEncoders(r encoder.Registry) {
	if r.WebpEncoder != nil {
		fmt.Printf("WEBP: %s\n", r.WebpEncoder.Name())
	}
	if r.AvifEncoder != nil {
		fmt.Printf("AVIF: %s\n", r.AvifEncoder.Name())
	}
	fmt.Println()
}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./cmd/kpar/`
Expected: no errors, `kpar` binary created

- [ ] **Step 4: Run all tests**

Run: `go test ./...`
Expected: all tests PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/kpar/main.go go.mod go.sum
git commit -m "Wire up CLI with Cobra, connect all components"
```

---

### Task 10: Integration Test — End to End

**Files:**
- Create: `internal/converter/integration_test.go`

- [ ] **Step 1: Write integration test**

Create `internal/converter/integration_test.go`:

```go
package converter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/janiosarmento/kpar/internal/converter"
	"github.com/janiosarmento/kpar/internal/encoder"
)

func TestIntegrationPNGProducesSmaller(t *testing.T) {
	registry := encoder.Detect()
	if registry.WebpEncoder == nil || registry.AvifEncoder == nil {
		t.Skip("need both webp and avif encoders")
	}

	dir := t.TempDir()
	src := createTestImage(t, dir, "large.png")

	result, err := converter.Convert(src, registry, -1)
	if err != nil {
		t.Fatalf("convert failed: %v", err)
	}

	// Original must still exist
	if _, err := os.Stat(src); err != nil {
		t.Error("original file was deleted")
	}

	// At most one converted file should remain
	kept := 0
	for _, c := range result.Conversions {
		if c.Kept {
			kept++
			if _, err := os.Stat(c.Path); err != nil {
				t.Errorf("kept file %s doesn't exist", c.Path)
			}
		} else {
			if _, err := os.Stat(c.Path); err == nil {
				t.Errorf("discarded file %s still exists", c.Path)
			}
		}
	}
	if kept > 1 {
		t.Errorf("expected at most 1 kept file, got %d", kept)
	}
}

func TestIntegrationWEBPToAVIF(t *testing.T) {
	registry := encoder.Detect()
	if registry.WebpEncoder == nil || registry.AvifEncoder == nil {
		t.Skip("need both webp and avif encoders")
	}

	dir := t.TempDir()
	pngSrc := createTestImage(t, dir, "source.png")

	// Create a real WEBP first
	webpPath := filepath.Join(dir, "source.webp")
	if err := registry.WebpEncoder.Encode(pngSrc, webpPath, -1); err != nil {
		t.Fatalf("setup webp: %v", err)
	}
	os.Remove(pngSrc) // clean up PNG

	result, err := converter.Convert(webpPath, registry, -1)
	if err != nil {
		t.Fatalf("convert failed: %v", err)
	}

	// Should only have AVIF conversion, no WEBP
	for _, c := range result.Conversions {
		if c.Format == "webp" {
			t.Error("WEBP input should not produce WEBP conversion")
		}
	}

	// Original WEBP must still exist
	if _, err := os.Stat(webpPath); err != nil {
		t.Error("original WEBP was deleted")
	}
}
```

- [ ] **Step 2: Run integration tests**

Run: `go test ./internal/converter/ -v -run TestIntegration`
Expected: PASS

- [ ] **Step 3: Run full test suite**

Run: `go test ./... -v`
Expected: all tests PASS

- [ ] **Step 4: Test the binary manually**

```bash
cd /Users/janiosarmento/projects/kpar
go build -o kpar ./cmd/kpar/
# Test with a real image if available, or just verify help
./kpar --help
```

- [ ] **Step 5: Clean up and commit**

```bash
rm -f kpar
git add internal/converter/integration_test.go
git commit -m "Add integration tests for end-to-end conversion flow"
```

---

### Task 11: Polish — .gitignore and Go Tidy

**Files:**
- Create: `.gitignore`

- [ ] **Step 1: Create .gitignore**

Create `.gitignore`:

```
kpar
*.exe
dist/
```

- [ ] **Step 2: Tidy Go modules**

```bash
cd /Users/janiosarmento/projects/kpar && go mod tidy
```

- [ ] **Step 3: Verify everything passes**

Run: `go test ./... && go build ./cmd/kpar/`
Expected: all tests pass, build succeeds

- [ ] **Step 4: Commit**

```bash
git add .gitignore go.mod go.sum
git commit -m "Add .gitignore and tidy Go modules"
```
