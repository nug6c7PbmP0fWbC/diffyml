package reporter

import (
	"io"

	"github.com/diffyml/diffyml/pkg/diff"
	"github.com/diffyml/diffyml/pkg/filter"
	"github.com/diffyml/diffyml/pkg/formatter"
)

// Config holds reporter configuration.
type Config struct {
	Format    string
	FilterCfg filter.Config
}

// Reporter orchestrates diffing, filtering, and formatting.
type Reporter struct {
	cfg       Config
	formatter formatter.Formatter
}

// New creates a new Reporter with the given config.
func New(cfg Config) (*Reporter, error) {
	f, err := formatter.New(cfg.Format)
	if err != nil {
		return nil, err
	}
	return &Reporter{cfg: cfg, formatter: f}, nil
}

// Report compares two YAML node maps, applies filters, and writes formatted output.
func (r *Reporter) Report(w io.Writer, a, b map[string]interface{}) error {
	changes := diff.Compare(a, b)
	changes = filter.Apply(changes, r.cfg.FilterCfg)
	return r.formatter.Format(w, changes)
}
