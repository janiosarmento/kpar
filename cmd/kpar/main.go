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

var (
	quality      int
	keepOriginal bool
)

var rootCmd = &cobra.Command{
	Use:   "kpar [file]",
	Short: "Image optimizer — converts to WEBP/AVIF and keeps the smallest",
	Long:  "KPAR converts images between modern formats (WEBP, AVIF) and keeps only the smallest result.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  run,
}

func init() {
	rootCmd.Flags().IntVarP(&quality, "quality", "q", -1, "Encoding quality (0-100, default: encoder default)")
	rootCmd.Flags().BoolVarP(&keepOriginal, "keep", "k", false, "Keep original file after optimization")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Detect encoders
	registry := encoder.Detect()
	reportEncoders(registry)

	if registry.WebpEncoder == nil && registry.AvifEncoder == nil {
		fmt.Fprintln(os.Stderr, "No encoders found. Install at least one:")
		fmt.Fprintln(os.Stderr, "  brew install webp libavif    (macOS)")
		fmt.Fprintln(os.Stderr, "  apt install webp libavif-bin (Debian/Ubuntu)")
		fmt.Fprintln(os.Stderr, "  Or install ImageMagick 7+: brew install imagemagick")
		return fmt.Errorf("no encoders available")
	}

	var filePath string

	if len(args) == 1 {
		// Direct file mode
		abs, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		filePath = abs
	} else {
		// Picker mode
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		files, err := scanner.Scan(dir)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			fmt.Println("No image files found in current directory.")
			return nil
		}

		selected, err := picker.Pick(files)
		if err != nil {
			return err
		}
		filePath = selected
	}

	// Convert
	result, err := converter.Convert(filePath, registry, quality, keepOriginal)
	if err != nil {
		return err
	}

	// Output
	fmt.Print(output.Render(result))
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
