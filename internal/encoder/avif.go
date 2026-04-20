package encoder

import (
	"fmt"
	"os/exec"
)

type AvifencEncoder struct{}

func (AvifencEncoder) Name() string { return "avifenc" }

func (AvifencEncoder) Encode(src, dst string, quality int) error {
	args := []string{}
	if quality >= 0 {
		args = append(args, "--min", fmt.Sprintf("%d", quality), "--max", fmt.Sprintf("%d", quality))
	}
	args = append(args, src, dst)
	cmd := exec.Command("avifenc", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("avifenc: %w\n%s", err, out)
	}
	return nil
}

type MagickAvifEncoder struct{}

func (MagickAvifEncoder) Name() string { return "magick (avif)" }

func (MagickAvifEncoder) Encode(src, dst string, quality int) error {
	args := []string{"convert", src}
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
