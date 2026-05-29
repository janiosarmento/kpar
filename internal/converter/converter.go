package converter

import (
	"fmt"
	"os"
	"os/exec"
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
	BestPath        string // path to the best converted file, or "" if no gain
	Saved           int64  // bytes saved (0 if no gain)
	OriginalDeleted bool
}

// Convert processes a single image file according to its type.
// If keepOriginal is false, the original file is deleted when a smaller conversion is found.
// If removeWM is true, the Gemini sparkle watermark is removed before encoding.
// If cropRight is true, 160px are cropped from the right edge before processing.
func Convert(src string, reg encoder.Registry, quality int, keepOriginal, removeWM, cropRight bool) (Result, error) {
	info, err := os.Stat(src)
	if err != nil {
		return Result{}, fmt.Errorf("stat %s: %w", src, err)
	}

	encodeSrc := src

	// Crop 160px from each side
	if cropRight {
		cropped, cleanup, err := cropSides(src, 160)
		if err != nil {
			return Result{}, fmt.Errorf("crop: %w", err)
		}
		defer cleanup()
		encodeSrc = cropped
	}

	// Remove watermark
	if removeWM {
		patched, cleanup, err := removeWatermark(encodeSrc)
		if err != nil {
			return Result{}, fmt.Errorf("watermark removal: %w", err)
		}
		defer cleanup()
		encodeSrc = patched
	}

	result := Result{
		OriginalPath: src,
		OriginalSize: info.Size(),
	}

	ext := strings.ToLower(filepath.Ext(src))
	base := strings.TrimSuffix(src, filepath.Ext(src))

	switch ext {
	case ".webp":
		err = convertWebp(base, encodeSrc, reg, quality, &result)
	case ".jpg", ".jpeg", ".png":
		err = convertJpgPng(base, encodeSrc, reg, quality, &result)
	case ".heic", ".heif":
		err = convertHeic(base, encodeSrc, reg, quality, &result)
	default:
		return result, fmt.Errorf("unsupported format: %s", ext)
	}

	if err != nil {
		return result, err
	}

	pickBest(&result)

	if result.Saved > 0 && !keepOriginal {
		os.Remove(src)
		result.OriginalDeleted = true
	}

	return result, nil
}

func convertWebp(base, src string, reg encoder.Registry, quality int, result *Result) error {
	if reg.AvifEncoder == nil {
		return fmt.Errorf("no AVIF encoder available")
	}

	// avifenc does not support WEBP input directly; decode to a temp PNG first.
	tmpPNG := base + "__tmp.png"
	defer os.Remove(tmpPNG)

	if err := webpToPNG(src, tmpPNG); err != nil {
		return fmt.Errorf("decoding webp to png: %w", err)
	}

	// Resize if wider than 1440px
	encodeSrc, cleanup, err := resizeIfNeeded(tmpPNG)
	if err != nil {
		return fmt.Errorf("resize: %w", err)
	}
	defer cleanup()

	dst := base + ".avif"
	if err := reg.AvifEncoder.Encode(encodeSrc, dst, quality); err != nil {
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

// webpToPNG decodes a WEBP file to PNG using dwebp (preferred) or magick.
func webpToPNG(src, dst string) error {
	if _, err := exec.LookPath("dwebp"); err == nil {
		cmd := exec.Command("dwebp", src, "-o", dst)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("dwebp: %w\n%s", err, out)
		}
		return nil
	}
	if _, err := exec.LookPath("magick"); err == nil {
		cmd := exec.Command("magick", "convert", src, dst)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("magick: %w\n%s", err, out)
		}
		return nil
	}
	return fmt.Errorf("no tool available to decode WEBP (need dwebp or magick)")
}

func convertHeic(base, src string, reg encoder.Registry, quality int, result *Result) error {
	tmpPNG := base + "__tmp.png"
	defer os.Remove(tmpPNG)

	if err := heicToPNG(src, tmpPNG); err != nil {
		return fmt.Errorf("decoding heic to png: %w", err)
	}

	return convertJpgPng(base, tmpPNG, reg, quality, result)
}

// heicToPNG decodes a HEIC/HEIF file to PNG using heif-convert (preferred) or magick.
func heicToPNG(src, dst string) error {
	if _, err := exec.LookPath("heif-convert"); err == nil {
		cmd := exec.Command("heif-convert", src, dst)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("heif-convert: %w\n%s", err, out)
		}
		return nil
	}
	if _, err := exec.LookPath("magick"); err == nil {
		cmd := exec.Command("magick", "convert", src, dst)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("magick: %w\n%s", err, out)
		}
		return nil
	}
	return fmt.Errorf("no tool available to decode HEIC/HEIF (need heif-convert or magick)")
}

func convertJpgPng(base, src string, reg encoder.Registry, quality int, result *Result) error {
	// Resize if wider than 1440px
	encodeSrc, cleanup, err := resizeIfNeeded(src)
	if err != nil {
		return fmt.Errorf("resize: %w", err)
	}
	defer cleanup()

	if reg.WebpEncoder != nil {
		dst := base + ".webp"
		if err := reg.WebpEncoder.Encode(encodeSrc, dst, quality); err != nil {
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
		if err := reg.AvifEncoder.Encode(encodeSrc, dst, quality); err != nil {
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
