// Utilities for scoring
package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
)

func Initialize(ctx context.Context, db *sql.DB) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	pattern := filepath.Join(root, "sql", "*.sql")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	sort.Slice(paths, func(i, j int) bool { return paths[i] < paths[j] })
	for _, path := range paths {
		log.Printf("Load sql file: %s\n", path)
		f, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to load sql: %s", path))
		}

		if _, err = db.ExecContext(ctx, string(f)); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to exec sql: %s", path))
		}
	}

	return nil
}
