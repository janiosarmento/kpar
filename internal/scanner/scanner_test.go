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
