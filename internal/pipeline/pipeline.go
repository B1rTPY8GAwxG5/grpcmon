// Package pipeline chains multiple entry processors in sequence.
package pipeline

import "github.com/grpcmon/internal/capture"

// Processor transforms or filters a slice of entries.
type Processor func([]capture.Entry) []capture.Entry

// Pipeline applies a sequence of Processors to entries.
type Pipeline struct {
	steps []Processor
}

// New returns a Pipeline with the given processors.
func New(steps ...Processor) *Pipeline {
	return &Pipeline{steps: steps}
}

// Add appends a processor to the pipeline.
func (p *Pipeline) Add(proc Processor) {
	p.steps = append(p.steps, proc)
}

// Run passes entries through each processor in order and returns the result.
func (p *Pipeline) Run(entries []capture.Entry) []capture.Entry {
	out := entries
	for _, step := range p.steps {
		out = step(out)
	}
	return out
}
