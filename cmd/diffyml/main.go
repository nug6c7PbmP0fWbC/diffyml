// Command diffyml compares two YAML files and reports the differences.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/diffyml/diffyml/pkg/filter"
	"github.com/diffyml/diffyml/pkg/loader"
	"github.com/diffyml/diffyml/pkg/reporter"
)

const version = "0.1.0"

func main() {
	var (
		format      = flag.String("format", "text", "Output format: text, json, yaml, markdown, table")
		filterPath  = flag.String("filter-path", "", "Filter changes by path prefix (e.g. 'server.host')")
		filterType  = flag.String("filter-type", "", "Filter changes by type: added, removed, modified")
		showVersion = flag.Bool("version", false, "Print version and exit")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: diffyml [options] <file1.yml> <file2.yml>\n\n")
		fmt.Fprintf(os.Stderr, "Compare two YAML files and report differences.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("diffyml version %s\n", version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	file1, file2 := args[0], args[1]

	a, err := loader.LoadFile(file1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", file1, err)
		os.Exit(1)
	}

	b, err := loader.LoadFile(file2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", file2, err)
		os.Exit(1)
	}

	var f *filter.Filter
	if *filterPath != "" || *filterType != "" {
		f = &filter.Filter{
			Path: *filterPath,
			Type: strings.ToLower(*filterType),
		}
	}

	r, err := reporter.New(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	output, err := r.Report(a, b, f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating report: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(output)
}
