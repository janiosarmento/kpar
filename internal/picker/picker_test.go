package picker_test

import (
	"strings"
	"testing"

	"github.com/janiosarmento/kpar/internal/picker"
)

func TestFormatChoices(t *testing.T) {
	files := []string{
		"/home/user/photos/beach.jpg",
		"/home/user/photos/sunset.png",
		"/home/user/photos/logo.webp",
	}

	formatted := picker.FormatChoices(files)
	if len(formatted) != 3 {
		t.Fatalf("got %d choices, want 3", len(formatted))
	}

	if !strings.Contains(formatted[0], "beach.jpg") {
		t.Errorf("formatted[0] = %q, want to contain 'beach.jpg'", formatted[0])
	}
}
