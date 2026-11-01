// Package label resolves a short -t name or path to a .zpl template file,
// renders it against a tag and the configured placeholders, and lists the
// templates available under a directory.
package label

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

var nameRE = regexp.MustCompile(`^[A-Za-z0-9_./-]+$`)

// Data is the set of values a template can reference, as {{.Tag}},
// {{.CompanyName}}, and {{.AssetURL}}.
type Data struct {
	Tag         string
	CompanyName string
	AssetURL    string
}

// Resolve returns the .zpl file path for name under dir. A name containing
// "/" is treated as an explicit path under dir; otherwise it's looked up
// in aliases, falling back to a flat name of the same shape so a template
// can be referenced directly without an alias.
func Resolve(dir, name string, aliases map[string]string) (string, error) {
	if !nameRE.MatchString(name) || strings.Contains(name, "..") {
		return "", fmt.Errorf("invalid template name: %s", name)
	}

	rel := name
	if !strings.Contains(name, "/") {
		if target, ok := aliases[name]; ok {
			rel = target
		}
	}

	path := filepath.Join(dir, rel+".zpl")
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("template not found: %s", path)
	}
	return path, nil
}

// Render reads the template at path and executes it against data. Any
// placeholder that doesn't resolve to a Data field — an unknown name, or
// one left over from a typo — is a hard error rather than being printed
// literally or silently dropped.
func Render(path string, data Data) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read template: %w", err)
	}

	tmpl, err := template.New(filepath.Base(path)).Parse(string(b))
	if err != nil {
		return "", fmt.Errorf("parsing template %s: %w", path, err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("rendering template %s: %w", path, err)
	}
	return buf.String(), nil
}

// List returns every template under dir, as paths relative to dir without
// the .zpl extension, sorted.
func List(dir string) ([]string, error) {
	var names []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".zpl") {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		names = append(names, strings.TrimSuffix(filepath.ToSlash(rel), ".zpl"))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}
