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
