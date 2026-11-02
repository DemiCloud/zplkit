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

	got, err := Render(filepath.Join(dir, "basic.zpl"), Data{"Tag": "ABC-000001", "CompanyName": "Example Corp"})
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

	got, err := Render(filepath.Join(dir, "qr.zpl"), Data{"Tag": "ABC-000001", "AssetURL": "https://example.com/assets/"})
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

	if _, err := Render(filepath.Join(dir, "typo.zpl"), Data{"CompanyName": "Example Corp"}); err == nil {
		t.Fatal("expected error for unknown placeholder, got nil")
	}
}

func TestRenderMissingFieldFails(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "needs-zip", "^FD{{.ZIP}}^FS")

	if _, err := Render(filepath.Join(dir, "needs-zip.zpl"), Data{"Name": "Jane Doe"}); err == nil {
		t.Fatal("expected error for a field nobody supplied, got nil")
	}
}

func TestRenderRejectsZPLBreakerCharacters(t *testing.T) {
	tests := []struct {
		name string
		data Data
	}{
		{"caret", Data{"Name": "Evil^XA^FO0,0^XZ"}},
		{"tilde", Data{"Name": "Evil~SD"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeTemplate(t, dir, "basic", "^FD{{.Name}}^FS")

			if _, err := Render(filepath.Join(dir, "basic.zpl"), tt.data); err == nil {
				t.Fatal("expected error for ZPL breaker character, got nil")
			}
		})
	}
}

func TestFitFont(t *testing.T) {
	tests := []struct {
		name             string
		text             string
		maxWidth         int
		max, min         int
		want             int
		wantClampedToMax bool
		wantClampedToMin bool
	}{
		{name: "empty text uses max", text: "", maxWidth: 390, max: 48, min: 20, want: 48},
		{name: "short text uses max", text: "Jam", maxWidth: 390, max: 48, min: 20, wantClampedToMax: true},
		{name: "long text shrinks below max", text: "Smoked Spanish Paprika", maxWidth: 390, max: 48, min: 20},
		{name: "extreme text clamps to min", text: "Extremely Long Spice Jar Label Name Example", maxWidth: 390, max: 48, min: 20, wantClampedToMin: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fitFont(tt.text, tt.maxWidth, tt.max, tt.min)
			if got > tt.max || got < tt.min {
				t.Fatalf("fitFont() = %d, want within [%d, %d]", got, tt.min, tt.max)
			}
			if tt.want != 0 && got != tt.want {
				t.Errorf("fitFont() = %d, want %d", got, tt.want)
			}
			if tt.wantClampedToMax && got != tt.max {
				t.Errorf("fitFont() = %d, want clamped to max %d", got, tt.max)
			}
			if tt.wantClampedToMin && got != tt.min {
				t.Errorf("fitFont() = %d, want clamped to min %d", got, tt.min)
			}
		})
	}
}

func TestRenderFitFontTemplateFunc(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "fit", `{{$sz := fitfont .Name 390 48 20}}^A0N,{{$sz}},{{$sz}}^FD{{.Name}}^FS`)

	got, err := Render(filepath.Join(dir, "fit.zpl"), Data{"Name": "Smoked Spanish Paprika"})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	want := "^A0N,32,32^FDSmoked Spanish Paprika^FS"
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
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
