package runner

import (
	"fmt"
	"runtime"
)

// Arguments represent the values consumed by step.Step during the execution
type Arguments struct {
	values map[string]string
}

// NewArgs create a new Arguments with default values populated
func NewArgs() Arguments {
	return Arguments{
		values: map[string]string{
			"os.name": runtime.GOOS,
			"os.arch": runtime.GOARCH,
		},
	}
}

// Values return the inner arguments values
func (a *Arguments) Values() map[string]string {
	return a.copyValues()
}

// Update given scope value and return new instance with value set
func (a *Arguments) Update(scope, key, value string) Arguments {
	values := a.copyValues()
	values[fmt.Sprintf("%s.%s", scope, key)] = value
	return Arguments{
		values: values,
	}
}

// Merge merge given arguments into current Arguments and return corresponding copy
func (a *Arguments) Merge(values map[string]string) Arguments {
	v := a.copyValues()

	for key, value := range values {
		v[key] = value
	}

	return Arguments{
		values: v,
	}
}

func (a *Arguments) copyValues() map[string]string {
	values := map[string]string{}

	for key, value := range a.values {
		values[key] = value
	}

	return values
}
