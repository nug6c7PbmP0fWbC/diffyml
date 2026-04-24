package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// AnnotatedFormatter wraps another Formatter and appends inline annotations
// produced by diff.Annotate to each changed line.
type AnnotatedFormatter struct {
	inner   Formatter
	rules   []diff.AnnotationRule
	writer  io.Writer
}

// NewAnnotatedFormatter creates an AnnotatedFormatter that decorates inner
// with annotation rules. Output is written to w.
func NewAnnotatedFormatter(inner Formatter, rules []diff.AnnotationRule, w io.Writer) *AnnotatedFormatter {
	return &AnnotatedFormatter{inner: inner, rules: rules, writer: w}
}

// Format writes the inner format output and then appends an annotation
// summary section when annotations are present.
func (a *AnnotatedFormatter) Format(changes []diff.Change) error {
	if err := a.inner.Format(changes); err != nil {
		return err
	}

	annotations := diff.Annotate(changes, a.rules)
	if len(annotations) == 0 {
		return nil
	}

	_, err := fmt.Fprintln(a.writer, "\n--- Annotations ---")
	if err != nil {
		return err
	}

	for path, anns := range annotations {
		for _, ann := range anns {
			line := fmt.Sprintf("[%s] %s: %s",
				strings.ToUpper(ann.Severity), path, ann.Message)
			if _, err := fmt.Fprintln(a.writer, line); err != nil {
				return err
			}
		}
	}
	return nil
}
