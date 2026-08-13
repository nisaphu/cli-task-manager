package architecture

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ADR 0001 requires task persistence to use local JSON files.
// This fitness function fails if an alternative database/persistence
// technology is introduced into Go source code.
func TestADR0001PersistenceGuardrail(t *testing.T) {
	root := ".."

	forbiddenImports := []string{
		"database/sql",
		"github.com/mattn/go-sqlite3",
		"modernc.org/sqlite",
		"gorm.io/",
		"entgo.io/",
		"go.etcd.io/bbolt",
		"github.com/dgraph-io/badger",
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == ".specify" {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, imp := range file.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)

			for _, forbidden := range forbiddenImports {
				if importPath == forbidden || strings.HasPrefix(importPath, forbidden) {
					t.Errorf(
						"ADR 0001 violation: forbidden persistence import %q found in %s; tasks must use local JSON files",
						importPath,
						path,
					)
				}
			}
		}

		return nil
	})
	if err != nil {
		t.Fatalf("architecture scan failed: %v", err)
	}

	// Positive assertion: the JSON persistence implementation
	// must continue to use encoding/json.
	jsonStorePath := filepath.Join(root, "internal", "store", "json.go")

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, jsonStorePath, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("cannot inspect JSON store %s: %v", jsonStorePath, err)
	}

	usesEncodingJSON := false
	for _, imp := range file.Imports {
		if strings.Trim(imp.Path.Value, `"`) == "encoding/json" {
			usesEncodingJSON = true
			break
		}
	}

	if !usesEncodingJSON {
		t.Errorf(
			"ADR 0001 violation: %s must use encoding/json for local JSON persistence",
			jsonStorePath,
		)
	}
}
