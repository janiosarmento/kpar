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
