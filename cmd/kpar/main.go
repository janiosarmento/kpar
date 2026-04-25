package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/janiosarmento/kpar/internal/converter"
	"github.com/janiosarmento/kpar/internal/encoder"
	"github.com/janiosarmento/kpar/internal/output"
	"github.com/janiosarmento/kpar/internal/picker"
	"github.com/janiosarmento/kpar/internal/scanner"
)

// version is set at build time via:
//
//	go build -ldflags "-X main.version=1.2.3"
var version = "dev"

var (
	quality      int
	keepOriginal bool
	blur         bool
	crop         bool
	original     bool
)

var rootCmd = &cobra.Command{
	Use:     "kpar [file...]",
	Short:   "Image optimizer — converts to WEBP/AVIF and keeps the smallest",
	Long:    "KPAR converts images between modern formats (WEBP, AVIF) and keeps only the smallest result.\nAccepts multiple files: kpar *.png, kpar gato-*, kpar Downloads/*.jpg",
	Version: version,
	RunE:    run,
}

func init() {
	rootCmd.Flags().IntVarP(&quality, "quality", "q", -1, "Encoding quality (0-100, default: encoder default)")
	rootCmd.Flags().BoolVarP(&keepOriginal, "keep", "k", false, "Keep original file after optimization")
	rootCmd.Flags().BoolVarP(&blur, "blur", "b", false, "Remove Gemini sparkle watermark via blur")
	rootCmd.Flags().BoolVarP(&crop, "crop", "c", false, "Force crop 160px from the right edge")
	rootCmd.Flags().BoolVarP(&original, "original", "o", false, "No filters — only resize and convert")
	rootCmd.Flags().BoolP("version", "v", false, "Print version and exit")
	rootCmd.Flags().Lookup("version").Hidden = true // evita duplicata no --help; --version já é gerado pelo Cobra
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if original && (blur || crop) {
		return fmt.Errorf("ei, você pediu para não fazer nada E para fazer algo ao mesmo tempo. Escolhe um lado!")
	}

	registry := encoder.Detect()
	reportEncoders(registry)

	if registry.WebpEncoder == nil && registry.AvifEncoder == nil {
		fmt.Fprintln(os.Stderr, "No encoders found. Install at least one:")
		fmt.Fprintln(os.Stderr, "  brew install webp libavif    (macOS)")
		fmt.Fprintln(os.Stderr, "  apt install webp libavif-bin (Debian/Ubuntu)")
		fmt.Fprintln(os.Stderr, "  Or install ImageMagick 7+: brew install imagemagick")
		return fmt.Errorf("no encoders available")
	}

	var files []string

	if len(args) > 0 {
		for _, arg := range args {
			abs, err := filepath.Abs(arg)
			if err != nil {
				return err
			}
			files = append(files, abs)
		}
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		scanned, err := scanner.Scan(dir)
		if err != nil {
			return err
		}

		if len(scanned) == 0 {
			fmt.Println("No image files found in current directory.")
			return nil
		}

		selected, err := picker.Pick(scanned)
		if err != nil {
			return err
		}
		files = []string{selected}
	}

	var hasError bool
	for _, filePath := range files {
		cropRight := false
		if !original {
			if crop {
				cropRight = true
			} else {
				gemini, err := converter.IsGeminiImage(filePath)
				if err == nil && gemini {
					cropRight = true
				}
			}
		}

		result, err := converter.Convert(filePath, registry, quality, keepOriginal, blur && !original, cropRight)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(filePath), err)
			hasError = true
			continue
		}
		fmt.Print(output.Render(result))
	}

	if hasError {
		return fmt.Errorf("some files failed")
	}
	return nil
}

func reportEncoders(r encoder.Registry) {
	if r.WebpEncoder != nil {
		fmt.Printf("WEBP: %s\n", r.WebpEncoder.Name())
	}
	if r.AvifEncoder != nil {
		fmt.Printf("AVIF: %s\n", r.AvifEncoder.Name())
	}
	fmt.Println()
}
