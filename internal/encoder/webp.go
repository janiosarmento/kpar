package encoder

import (
	"fmt"
	"os/exec"
)

type CwebpEncoder struct{}

func (CwebpEncoder) Name() string { return "cwebp" }

func (CwebpEncoder) Encode(src, dst string, quality int) error {
	args := []string{"-metadata", "none", src, "-o", dst}
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
	args := []string{"convert", src, "-strip"}
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
