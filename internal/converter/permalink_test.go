package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPermalink(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Foo Bar Baz", "foo-bar-baz"},
		{"Crème Brûlée Café", "creme-brulee-cafe"},
		{"naïve résumé", "naive-resume"},
		{"Straße Æther Øresund", "strasse-aether-oresund"},
		{"my_file.v2", "my-file-v2"},
		{"  many   spaces  ", "many-spaces"},
		{"Hello!!! (copy) #1", "hello-copy-1"},
		{"already-a-slug", "already-a-slug"},
		{"Gemini_Generated_Image_abc123", "gemini-generated-image-abc123"},
		{"???", "image"},
		{"", "image"},
	}

	for _, tt := range tests {
		if got := permalink(tt.name); got != tt.want {
			t.Errorf("permalink(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestUniquePath(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "photo")

	if got := uniquePath(base, "avif"); got != base+".avif" {
		t.Errorf("free path: got %q, want %q", got, base+".avif")
	}

	os.WriteFile(base+".avif", nil, 0o644)
	if got := uniquePath(base, "avif"); got != base+"-2.avif" {
		t.Errorf("taken path: got %q, want %q", got, base+"-2.avif")
	}

	os.WriteFile(base+"-2.avif", nil, 0o644)
	if got := uniquePath(base, "avif"); got != base+"-3.avif" {
		t.Errorf("taken twice: got %q, want %q", got, base+"-3.avif")
	}

	if got := uniquePath(base, "webp"); got != base+".webp" {
		t.Errorf("other extension should be free: got %q, want %q", got, base+".webp")
	}
}
