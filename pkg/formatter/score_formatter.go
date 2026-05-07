package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// ScoreFormatter renders a human-readable similarity score report.
type ScoreFormatter struct {
	w io.Writer
}

// NewScoreFormatter creates a ScoreFormatter that writes to w.
func NewScoreFormatter(w io.Writer) *ScoreFormatter {
	return &ScoreFormatter{w: w}
}

// Format writes the score report to the underlying writer.
func (f *ScoreFormatter) Format(s diff.Score) error {
	bar := buildBar(s.Similarity, 20)
	_, err := fmt.Fprintf(f.w,
		"Similarity Score\n"+
			"  Total changes : %d\n"+
			"  Score         : %.1f%%\n"+
			"  [%s]\n",
		s.Total,
		s.Similarity*100,
		bar,
	)
	return err
}

// buildBar returns an ASCII progress bar of the given width representing ratio (0..1).
func buildBar(ratio float64, width int) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	filled := int(ratio * float64(width))
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}
