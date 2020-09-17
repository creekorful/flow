package main

import (
	"bytes"
	"fmt"
	"github.com/creekorful/flow/internal/runner"
	"github.com/creekorful/flow/internal/step"
	"github.com/creekorful/flow/internal/workflow"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Level(zerolog.InfoLevel)

	app := &cli.App{
		Usage:   "Run complex workflows using simple definition files.",
		Authors: []*cli.Author{{Name: "Aloïs Micard", Email: "alois@micard.lu"}},
		Version: "0.1.0",
		Before:  before,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "cache-dir",
				Usage:    "Path to the workflows cache directory",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "Default application log level",
				Value: "info",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "run",
				Aliases:   []string{"r"},
				Usage:     "Run a workflow",
				ArgsUsage: "URI",
				Action:    runWorkflow,
			},
			{
				Name:      "validate",
				Aliases:   []string{"v"},
				Usage:     "Validate a workflow",
				ArgsUsage: "URI",
				Action:    validateWorkflow,
			},
			{
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "Create a new workflow",
				Action:  newWorkflow,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err).Msg("error while running app.")
		os.Exit(1)
	}
}

func before(c *cli.Context) error {
	lvl, err := zerolog.ParseLevel(c.String("log-level"))
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	log.Logger.Level(lvl)
	return nil
}

func runWorkflow(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("missing URI")
	}

	uri := c.Args().First()
	w, err := readWorkflow(uri)
	if err != nil {
		return err
	}

	r, err := runner.NewRunner(c.String("cache-dir"))
	if err != nil {
		return err
	}

	return r.Run(&w, runner.NewArgs())
}

func validateWorkflow(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("missing URI")
	}

	uri := c.Args().First()
	_, err := readWorkflow(uri)
	if err != nil { // todo improve validation
		return fmt.Errorf("invalid workflow %s (%s)", uri, err)
	}

	log.Info().Str("uri", uri).Msg("workflow is valid.")

	return nil
}

func newWorkflow(_ *cli.Context) error {
	w := workflow.Workflow{
		ID:          "hello-world",
		Name:        "Hello World",
		Author:      "Aloïs Micard <alois@micard.lu>",
		Description: "Simple hello world.",
		Steps: []step.Step{
			{
				Exec: "echo Hello, world.",
			},
		},
	}

	var buf bytes.Buffer
	if err := workflow.Write(&buf, w); err != nil {
		return err
	}

	fmt.Printf("%s", buf.String())
	return nil
}

func readWorkflow(uri string) (workflow.Workflow, error) {
	// todo manage other sources
	f, err := os.Open(uri)
	if err != nil {
		return workflow.Workflow{}, err
	}
	defer f.Close()

	return workflow.Read(f)
}
