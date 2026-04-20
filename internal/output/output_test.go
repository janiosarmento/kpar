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
