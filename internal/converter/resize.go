package converter

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const maxWidth = 1440

// IsGeminiImage reports whether the file contains a Google C2PA provenance
// certificate, which Gemini embeds in all AI-generated images.
func IsGeminiImage(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return false, err
	}
	return bytes.Contains(data, []byte("Google C2PA")), nil
}

// imageWidth returns the width of a JPEG or PNG image by reading only the header.
func imageWidth(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, fmt.Errorf("decode config %s: %w", path, err)
	}
	return cfg.Width, nil
}

// cropSides removes pixels from both the left and right edges of an image.
// Returns the path to the cropped file and a cleanup function.
func cropSides(src string, pixels int) (string, func(), error) {
	noop := func() {}

	w, err := imageWidth(src)
	if err != nil {
		return "", noop, err
	}

	if w <= pixels*2 {
		return src, noop, nil
	}

	dir := filepath.Dir(src)
	base := filepath.Base(src)
	tmp := filepath.Join(dir, "__cropped_"+base)

	newWidth := w - pixels*2
	crop := fmt.Sprintf("%dx+%d+0", newWidth, pixels)
	cmd := exec.Command("magick", src, "-crop", crop, "+repage", tmp)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", noop, fmt.Errorf("magick crop: %w\n%s", err, out)
	}

	cleanup := func() { os.Remove(tmp) }
	return tmp, cleanup, nil
}

// resizeIfNeeded checks if the image is wider than maxWidth and, if so,
// resizes it to a temporary file. Returns the path to use for encoding
// (original or resized) and a cleanup function.
func resizeIfNeeded(src string) (string, func(), error) {
	noop := func() {}

	w, err := imageWidth(src)
	if err != nil {
		return "", noop, err
	}

	if w <= maxWidth {
		return src, noop, nil
	}

	ext := filepath.Ext(src)
	dir := filepath.Dir(src)
	base := filepath.Base(src)
	tmp := filepath.Join(dir, "__resized_"+base)
	// Preserve extension for format detection
	if filepath.Ext(tmp) != ext {
		tmp = tmp + ext
	}

	// Use magick to resize, maintaining aspect ratio
	resize := fmt.Sprintf("%dx", maxWidth)
	cmd := exec.Command("magick", "convert", src, "-resize", resize+">", tmp)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", noop, fmt.Errorf("magick resize: %w\n%s", err, out)
	}

	cleanup := func() { os.Remove(tmp) }
	return tmp, cleanup, nil
}
