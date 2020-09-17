package runner

import (
	"fmt"
	"github.com/creekorful/flow/internal/eval"
	"github.com/creekorful/flow/internal/step"
	"github.com/creekorful/flow/internal/workflow"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Runner is the running context of a given workflow.Workflow
type Runner struct {
	workflows map[string]workflow.Workflow
}

// NewRunner create a Runner and load workflow cache by searching in given directory (cacheDir)
func NewRunner(cacheDir string) (*Runner, error) {
	r := &Runner{
		workflows: map[string]workflow.Workflow{},
	}

	if err := filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			f, err := os.Open(path)
			if err != nil {
				log.Warn().Str("path", path).Msg("unable to open file.")
				return nil
			}

			w, err := workflow.Read(f)
			if err != nil {
				log.Warn().Str("path", path).Msg("invalid workflow file.")
				return nil
			}

			if w.ID == "" {
				log.Warn().Str("path", path).Msg("skipping workflow without id.")
				return nil
			}

			// prevent conflicts
			if _, exist := r.workflows[w.ID]; exist {
				return fmt.Errorf("duplicate workflow %s found", w.ID)
			}

			r.workflows[w.ID] = w
			log.Trace().Str("id", w.ID).Msg("workflow loaded.")
		}

		return nil
	}); err != nil {
		return nil, err
	}
	log.Debug().Int("workflows", len(r.workflows)).Msg("initialized runner.")

	return r, nil
}

// Run perform the workflow run, using given arguments
func (r *Runner) Run(workflow *workflow.Workflow, args Arguments) error {
	log.Debug().Str("workflow", workflow.Name).Msg("running workflow.")

	for _, s := range workflow.Steps {
		// make sure step should be executed
		if s.Cond != "" {
			if !eval.Evaluate(s.Cond, args.Values()) {
				log.Debug().Str("step", s.Name).Str("cond", s.Cond).Msg("skipping step.")
				continue
			}
		}

		// if step is directly runnable, well... run it
		if s.Runnable() {
			res, err := r.runStep(s, args)
			if err != nil {
				return err
			}

			// propagate step result
			args = args.Update("in", "result", res)

			continue
		}

		// not runnable, fetch linked workflow and run it
		w, err := r.findWorkflow(s.Uses)
		if err != nil {
			return err
		}
		if err := r.Run(w, args.Merge(s.Args)); err != nil {
			return err
		}
	}

	return nil
}

func (r *Runner) runStep(step step.Step, args Arguments) (string, error) {
	log.Debug().Str("step", step.Name).Msg("running step.")

	if step.Exec == "" {
		return "", fmt.Errorf("step is not runnable")
	}

	cmdLine, err := interpolateCmd(step.Exec, args.Values())
	if err != nil {
		return "", err
	}

	parts := strings.Split(cmdLine, " ")

	cmd := exec.Command(parts[0], parts[1:]...)
	log.Trace().Str("cmd", cmd.String()).Msg("running command.")

	b, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (r *Runner) findWorkflow(id string) (*workflow.Workflow, error) {
	log.Debug().Str("workflow", id).Msg("finding workflow.")

	if w, ok := r.workflows[id]; ok {
		return &w, nil
	}

	return nil, fmt.Errorf("no workflow `%s` found", id)
}

func interpolateCmd(cmd string, args map[string]string) (string, error) {
	res := cmd

	for key, value := range args {
		res = strings.Replace(res, fmt.Sprintf("{%s}", key), value, -1)
	}

	// make sure everything is extrapolate (todo improve)
	if strings.Contains(res, "{") && strings.Contains(res, "}") {
		return "", fmt.Errorf("missing variables")
	}

	return res, nil
}
