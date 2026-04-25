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

	result, err := converter.Convert(src, registry, -1, true, false, false)
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

	result, err := converter.Convert(webpPath, registry, -1, true, false, false)
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
