package label

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemplate(t *testing.T, dir, rel, contents string) {
	t.Helper()
	path := filepath.Join(dir, rel+".zpl")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func TestResolveByAlias(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "2x1/001", "^XA^XZ")
	aliases := map[string]string{"default": "2x1/001"}

	path, err := Resolve(dir, "default", aliases)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	want := filepath.Join(dir, "2x1/001.zpl")
	if path != want {
		t.Errorf("Resolve() = %q, want %q", path, want)
	}
}

func TestResolveByExplicitPath(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "2x1/001", "^XA^XZ")

	path, err := Resolve(dir, "2x1/001", nil)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	want := filepath.Join(dir, "2x1/001.zpl")
	if path != want {
		t.Errorf("Resolve() = %q, want %q", path, want)
	}
}

func TestResolveByExplicitPathWithDottedDir(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "2.25x1.25/001", "^XA^XZ")

	path, err := Resolve(dir, "2.25x1.25/001", nil)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	want := filepath.Join(dir, "2.25x1.25/001.zpl")
	if path != want {
		t.Errorf("Resolve() = %q, want %q", path, want)
	}
}

func TestResolveFlatFallback(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "solo", "^XA^XZ")

	path, err := Resolve(dir, "solo", nil)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	want := filepath.Join(dir, "solo.zpl")
	if path != want {
		t.Errorf("Resolve() = %q, want %q", path, want)
	}
}

func TestResolveUnknownNameFails(t *testing.T) {
	dir := t.TempDir()
	if _, err := Resolve(dir, "bogus", nil); err == nil {
		t.Fatal("expected error for unknown template, got nil")
	}
}

func TestResolveRejectsPathTraversal(t *testing.T) {
	dir := t.TempDir()
	if _, err := Resolve(dir, "../../etc/passwd", nil); err == nil {
		t.Fatal("expected error for path traversal, got nil")
	}
}

func TestResolveRejectsInvalidCharacters(t *testing.T) {
	dir := t.TempDir()
	if _, err := Resolve(dir, "owner;rm -rf", nil); err == nil {
		t.Fatal("expected error for invalid characters, got nil")
	}
}

func TestRenderSubstitutesFields(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "basic", "^FD{{.CompanyName}}^FS^FD{{.Tag}}^FS")

	got, err := Render(filepath.Join(dir, "basic.zpl"), Data{Tag: "ABC-000001", CompanyName: "Example Corp"})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	want := "^FDExample Corp^FS^FDABC-000001^FS"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestRenderConcatenatesAssetURLAndTag(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "qr", "^FDQA,{{.AssetURL}}{{.Tag}}^FS")

	got, err := Render(filepath.Join(dir, "qr.zpl"), Data{Tag: "ABC-000001", AssetURL: "https://example.com/assets/"})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	want := "^FDQA,https://example.com/assets/ABC-000001^FS"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestRenderUnknownPlaceholderFails(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "typo", "^FD{{.CompnayName}}^FS")

	if _, err := Render(filepath.Join(dir, "typo.zpl"), Data{CompanyName: "Example Corp"}); err == nil {
		t.Fatal("expected error for unknown placeholder, got nil")
	}
}

func TestListSortedRelativePaths(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "2x1/001", "^XA^XZ")
	writeTemplate(t, dir, "2x1/003", "^XA^XZ")
	writeTemplate(t, dir, "2.25x1.25/001", "^XA^XZ")

	got, err := List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	want := []string{
		"2.25x1.25/001",
		"2x1/001",
		"2x1/003",
	}
	if len(got) != len(want) {
		t.Fatalf("List() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("List()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
