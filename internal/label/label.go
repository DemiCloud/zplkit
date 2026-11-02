// Package label resolves a short -t name or path to a .zpl template file,
// renders it against a set of fields, and lists the templates available
// under a directory.
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

// zplBreakers are characters that start a new ZPL command (^) or control
// instruction (~); a field value containing one could inject commands into
// a rendered label rather than just appearing as label text.
const zplBreakers = "^~"

// Data is the set of values a template can reference as {{.Key}}. It's
// deliberately a map rather than a fixed struct: a template's field needs
// vary by what it's labeling (an asset tag, a mail-from/to address, a
// generic box/jar label, ...), and a new category or field should never
// require a Go code change. CompanyName/AssetURL (from config) and Tag
// (from the print command's TAG argument or a --field flag) are the keys
// in common use today, but any key a template references via {{.Field}}
// can be supplied the same way.
type Data map[string]string

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
// placeholder that doesn't resolve to a key in data — an unknown name, one
// left over from a typo, or one nobody supplied — is a hard error rather
// than being printed literally or silently dropped (via missingkey=error,
// since a plain map lookup would otherwise render a missing key as empty).
func Render(path string, data Data) (string, error) {
	if err := validateFieldValues(data); err != nil {
		return "", err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read template: %w", err)
	}

	tmpl, err := template.New(filepath.Base(path)).Option("missingkey=error").Parse(string(b))
	if err != nil {
		return "", fmt.Errorf("parsing template %s: %w", path, err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("rendering template %s: %w", path, err)
	}
	return buf.String(), nil
}

// validateFieldValues rejects any Data value containing a ZPL-breaker
// character, regardless of which key it's under or where it came from
// (config, a --field flag, ...). This is the one place that guarantee is
// enforced, rather than relying on every value's source to remember to
// check it itself.
func validateFieldValues(data Data) error {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if strings.ContainsAny(data[k], zplBreakers) {
			return fmt.Errorf("field %q must not contain %q", k, zplBreakers)
		}
	}
	return nil
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
