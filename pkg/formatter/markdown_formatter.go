package formatter

import (
	"fmt"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

type markdownFormatter struct{}

// NewMarkdownFormatter returns a Formatter that renders changes as a Markdown table.
func NewMarkdownFormatter() Formatter {
	return &markdownFormatter{}
}

func (f *markdownFormatter) Format(changes []diff.Change) (string, error) {
	if len(changes) == 0 {
		return "_No changes detected._\n", nil
	}

	var sb strings.Builder
	sb.WriteString("| Path | Type | Old Value | New Value |\n")
	sb.WriteString("|------|------|-----------|-----------|\n")

	for _, c := range changes {
		oldVal := formatMDValue(c.OldValue)
		newVal := formatMDValue(c.NewValue)
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
			c.Path, c.Type, oldVal, newVal))
	}

	return sb.String(), nil
}

func formatMDValue(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
