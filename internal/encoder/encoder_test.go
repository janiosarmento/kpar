package encoder_test

import (
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
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
