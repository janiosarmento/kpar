package encoder

import "os/exec"

// Encoder converts a source image to a target format.
type Encoder interface {
	// Name returns the encoder name (e.g., "cwebp", "magick").
	Name() string

	// Encode converts src to the target format, writing to dst.
	// quality is 0-100; if <0, use encoder default.
	Encode(src, dst string, quality int) error
}

// Registry holds the best available encoder for each format.
type Registry struct {
	WebpEncoder Encoder
	AvifEncoder Encoder
}

// Detect finds the best available encoders in PATH.
// Priority: dedicated encoders > ImageMagick fallback.
func Detect() Registry {
	var r Registry

	if commandExists("cwebp") {
		r.WebpEncoder = CwebpEncoder{}
	} else if commandExists("magick") {
		r.WebpEncoder = MagickWebpEncoder{}
	}

	if commandExists("avifenc") {
		r.AvifEncoder = AvifencEncoder{}
	} else if commandExists("magick") {
		r.AvifEncoder = MagickAvifEncoder{}
	}

	return r
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Stubs — implemented in avif.go
type AvifencEncoder struct{}

func (AvifencEncoder) Name() string                        { return "avifenc" }
func (AvifencEncoder) Encode(src, dst string, q int) error { return nil }

type MagickAvifEncoder struct{}

func (MagickAvifEncoder) Name() string                        { return "magick (avif)" }
func (MagickAvifEncoder) Encode(src, dst string, q int) error { return nil }
