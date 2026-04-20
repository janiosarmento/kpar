package output

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/janiosarmento/kpar/internal/converter"
)

var (
	green = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	red   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	gray  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	bold  = lipgloss.NewStyle().Bold(true)
)

// FormatSize formats bytes into a human-readable string.
func FormatSize(bytes int64) string {
	switch {
	case bytes >= 1048576:
		return fmt.Sprintf("%.1f MB", float64(bytes)/1048576)
	case bytes >= 1024:
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// Render formats a conversion result for terminal display.
func Render(r converter.Result) string {
	var b strings.Builder

	name := filepath.Base(r.OriginalPath)
	fmt.Fprintf(&b, "%s (%s)\n", bold.Render(name), FormatSize(r.OriginalSize))

	last := len(r.Conversions) - 1
	for i, c := range r.Conversions {
		connector := "├─"
		if i == last && r.Saved == 0 {
			connector = "└─"
		}

		pct := float64(c.Size) / float64(r.OriginalSize) * 100
		sizeStr := fmt.Sprintf("%s: %s (%.*f%% do original)",
			c.Format, FormatSize(c.Size), 0, pct)

		if c.Kept {
			fmt.Fprintf(&b, "  %s %s %s\n", connector, green.Render(sizeStr), green.Render("✓ salvo"))
		} else {
			fmt.Fprintf(&b, "  %s %s %s\n", connector, gray.Render(sizeStr), red.Render("✗ descartado"))
		}
	}

	if r.Saved > 0 {
		savings := fmt.Sprintf("ganho: %s economizados", FormatSize(r.Saved))
		fmt.Fprintf(&b, "  └─ %s\n", green.Render(savings))
	} else {
		fmt.Fprintf(&b, "  └─ %s\n", gray.Render("sem ganho"))
	}

	return b.String()
}
