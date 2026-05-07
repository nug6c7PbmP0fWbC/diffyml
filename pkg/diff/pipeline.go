package diff

// Pipeline represents a chain of transformation steps applied to a slice of Changes.
// Each step is a function that receives the current changes and returns a (possibly modified)
// slice of changes along with an error.
type Pipeline struct {
	steps []PipelineStep
}

// PipelineStep is a single transformation applied inside a Pipeline.
type PipelineStep func(changes []Change) ([]Change, error)

// NewPipeline creates an empty Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// Use appends one or more steps to the pipeline.
func (p *Pipeline) Use(steps ...PipelineStep) *Pipeline {
	p.steps = append(p.steps, steps...)
	return p
}

// Run executes all steps in order, threading the changes through each one.
// Execution stops and the error is returned as soon as any step fails.
func (p *Pipeline) Run(changes []Change) ([]Change, error) {
	var err error
	for _, step := range p.steps {
		changes, err = step(changes)
		if err != nil {
			return nil, err
		}
	}
	return changes, nil
}

// IgnoreStep returns a PipelineStep that removes changes matching the given ignore rules.
func IgnoreStep(rules []IgnoreRule) PipelineStep {
	return func(changes []Change) ([]Change, error) {
		return ApplyIgnore(changes, rules), nil
	}
}

// RedactStep returns a PipelineStep that redacts sensitive paths.
func RedactStep(rules []RedactRule) PipelineStep {
	return func(changes []Change) ([]Change, error) {
		return Redact(changes, rules), nil
	}
}

// AnnotateStep returns a PipelineStep that annotates changes with metadata.
func AnnotateStep(rules []AnnotateRule) PipelineStep {
	return func(changes []Change) ([]Change, error) {
		return Annotate(changes, rules), nil
	}
}
