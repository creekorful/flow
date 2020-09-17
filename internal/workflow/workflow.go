package workflow

import (
	"github.com/creekorful/flow/internal/step"
	"gopkg.in/yaml.v2"
	"io"
)

// Workflow represent a sequence of instructions (step.Step)
type Workflow struct {
	ID          string      `yaml:"id,omitempty"`
	Name        string      `yaml:"name,omitempty"`
	Author      string      `yaml:"author,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Steps       []step.Step `yaml:"steps,omitempty"`
}

// Write given Workflow into given io.Writer
func Write(writer io.Writer, workflow Workflow) error {
	b, err := yaml.Marshal(&workflow)
	if err != nil {
		return err
	}

	_, err = writer.Write(b)
	return err
}

// Read Workflow from given io.Reader
func Read(reader io.Reader) (Workflow, error) {
	var workflow Workflow
	if err := yaml.NewDecoder(reader).Decode(&workflow); err != nil {
		return Workflow{}, err
	}

	return workflow, nil
}
