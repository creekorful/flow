package step

// Step represent a part of the workflow
// a step consume arguments and may produce result
type Step struct {
	Name string            `yaml:"name,omitempty"`
	Exec string            `yaml:"exec,omitempty"`
	Uses string            `yaml:"uses,omitempty"`
	Args map[string]string `yaml:"args,omitempty"`
	Cond string            `yaml:"cond,omitempty"`
}

// Runnable determinate if given step may be run directly
func (s *Step) Runnable() bool {
	return s.Exec != ""
}
