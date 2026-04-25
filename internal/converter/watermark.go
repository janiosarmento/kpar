package converter

import (
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// removeWatermark removes the Gemini sparkle watermark from the bottom-right
// corner. It paints over the sparkle with a sampled background color, then
// applies a feathered blur to blend the patch into the surrounding image.
// Returns the path to use for encoding and a cleanup function.
func removeWatermark(src string) (string, func(), error) {
	noop := func() {}

	w, err := imageWidth(src)
	if err != nil {
		return "", noop, err
	}
	h, err := imageHeight(src)
	if err != nil {
		return "", noop, err
	}

	dir := filepath.Dir(src)
	base := filepath.Base(src)
	tmp := filepath.Join(dir, "__nowm_"+base)
	tmpFill := filepath.Join(dir, "__nowm_fill_"+base)

	// Sample the background color near the sparkle (above and left of it).
	bgColor, err := sampleColor(src, w, h)
	if err != nil {
		return "", noop, err
	}

	// Step 1: paint a solid rectangle over the sparkle area.
	sx1 := w - 85
	sy1 := h - 65
	rect := fmt.Sprintf("rectangle %d,%d %d,%d", sx1, sy1, w, h)
	cmd := exec.Command("magick", src,
		"-fill", bgColor, "-draw", rect,
		tmpFill,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", noop, fmt.Errorf("magick fill: %w\n%s", err, out)
	}

	// Step 2: feathered blur over the painted area to blend it naturally.
	bx := w - 110
	by := h - 85
	maskRect := fmt.Sprintf("rectangle %d,%d %d,%d", bx, by, w, h)
	ws := strconv.Itoa(w)
	hs := strconv.Itoa(h)
	cmd = exec.Command("magick", tmpFill,
		"(", "+clone", "-blur", "0x15",
		"(", "-size", ws+"x"+hs, "xc:black",
		"-fill", "white", "-draw", maskRect,
		"-blur", "0x12", ")",
		"-alpha", "off", "-compose", "CopyOpacity", "-composite", ")",
		"-compose", "over", "-composite",
		tmp,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpFill)
		return "", noop, fmt.Errorf("magick blend: %w\n%s", err, out)
	}
	os.Remove(tmpFill)

	cleanup := func() { os.Remove(tmp) }
	return tmp, cleanup, nil
}

// sampleColor returns the average color of a small patch above the sparkle area.
func sampleColor(src string, w, h int) (string, error) {
	// Sample a 30x15 patch above-left of the sparkle position.
	crop := fmt.Sprintf("30x15+%d+%d", w-90, h-80)
	cmd := exec.Command("magick", src,
		"-crop", crop, "-resize", "1x1!",
		"-format", "%[pixel:u]", "info:",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("magick sample: %w\n%s", err, out)
	}
	return string(out), nil
}

// imageHeight returns the height of an image by reading only the header.
func imageHeight(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, fmt.Errorf("decode config %s: %w", path, err)
	}
	return cfg.Height, nil
}
