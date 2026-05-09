package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// ThresholdFormatter renders threshold violation results as human-readable text.
type ThresholdFormatter struct{}

// NewThresholdFormatter returns a new ThresholdFormatter.
func NewThresholdFormatter() *ThresholdFormatter {
	return &ThresholdFormatter{}
}

// Format writes threshold violations to w. If there are none, it writes an
// "all thresholds satisfied" message instead.
func (f *ThresholdFormatter) Format(w io.Writer, violations []diff.ThresholdViolation) error {
	if len(violations) == 0 {
		_, err := fmt.Fprintln(w, "threshold: all limits satisfied")
		return err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("threshold: %d violation(s) detected\n", len(violations)))
	for _, v := range violations {
		sb.WriteString(fmt.Sprintf("  [%s] %s\n", strings.ToUpper(v.Kind), v.Message))
	}

	_, err := fmt.Fprint(w, sb.String())
	return err
}
