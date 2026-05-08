package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// maskFormatter renders masked changes, highlighting which paths were masked.
type maskFormatter struct {
	w io.Writer
}

// NewMaskFormatter returns a Formatter that shows masked changes
// and indicates which values have been hidden.
func NewMaskFormatter(w io.Writer) Formatter {
	return &maskFormatter{w: w}
}

func (f *maskFormatter) Format(changes []diff.Change) error {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(f.w, "no changes")
		return err
	}
	for _, c := range changes {
		masked := isMasked(c)
		line := buildMaskLine(c, masked)
		if _, err := fmt.Fprintln(f.w, line); err != nil {
			return err
		}
	}
	return nil
}

func isMasked(c diff.Change) bool {
	const placeholder = "***"
	ov, _ := c.OldValue.(string)
	nv, _ := c.NewValue.(string)
	return ov == placeholder || nv == placeholder
}

func buildMaskLine(c diff.Change, masked bool) string {
	tag := strings.ToUpper(string(c.Type))
	suffix := ""
	if masked {
		suffix = " [masked]"
	}
	switch c.Type {
	case diff.ChangeAdded:
		return fmt.Sprintf("+ [%s] %s: %v%s", tag, c.Path, c.NewValue, suffix)
	case diff.ChangeRemoved:
		return fmt.Sprintf("- [%s] %s: %v%s", tag, c.Path, c.OldValue, suffix)
	default:
		return fmt.Sprintf("~ [%s] %s: %v -> %v%s", tag, c.Path, c.OldValue, c.NewValue, suffix)
	}
}
