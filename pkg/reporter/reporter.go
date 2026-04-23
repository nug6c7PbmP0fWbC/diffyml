package reporter

import (
	"fmt"
	"io"

	"github.com/szhekpisov/diffyml/pkg/diff"
	"github.com/szhekpisov/diffyml/pkg/filter"
	"github.com/szhekpisov/diffyml/pkg/formatter"
)

// Config holds the options for a Reporter.
type Config struct {
	Format     string
	FilterPath string
	FilterType string
}

// Reporter compares two YAML documents and writes a formatted diff report.
type Reporter struct {
	cfg Config
}

// New creates a Reporter with the given Config.
func New(cfg Config) *Reporter {
	return &Reporter{cfg: cfg}
}

// Run compares old and new YAML maps, applies optional filters, formats the
// result using the configured formatter, and writes the output to w.
func (r *Reporter) Run(old, new map[string]interface{}, w io.Writer) error {
	changes, err := diff.Compare(old, new)
	if err != nil {
		return fmt.Errorf("reporter: compare: %w", err)
	}

	if r.cfg.FilterPath != "" || r.cfg.FilterType != "" {
		changes = filter.Apply(changes, r.cfg.FilterPath, r.cfg.FilterType)
	}

	fmt, err := formatter.New(r.cfg.Format)
	if err != nil {
		return fmt.Errorf("reporter: formatter: %w", err)
	}

	output, err := fmt.Format(changes)
	if err != nil {
		return fmt.Errorf("reporter: format: %w", err)
	}

	_, err = w.Write([]byte(output))
	return err
}
