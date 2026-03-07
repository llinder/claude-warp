package main

import (
	"context"
	"fmt"
	"os"

	"github.com/llinder/claude-warp/internal/hooks"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "claude-warp",
		Usage: "Warp terminal integration for Claude Code",
		Commands: []*cli.Command{
			{
				Name:  "session-start",
				Usage: "Handle session start (Warp detection, project context)",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return hooks.SessionStart()
				},
			},
			{
				Name:  "stop",
				Usage: "Handle session stop (notification with summary)",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return hooks.Stop()
				},
			},
			{
				Name:  "notification",
				Usage: "Handle notification (forward to Warp)",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return hooks.Notification()
				},
			},
			{
				Name:  "save-workflow",
				Usage: "Save a command as a Warp workflow",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Usage:    "workflow name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "command",
						Usage:    "command template with {{args}}",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "description",
						Usage: "workflow description",
					},
					&cli.BoolFlag{
						Name:  "repo",
						Usage: "save to repo .warp/workflows/ instead of personal",
					},
					&cli.StringSliceFlag{
						Name:  "arg",
						Usage: "argument spec: name:description:default (repeatable)",
					},
					&cli.StringSliceFlag{
						Name:  "tag",
						Usage: "tag for the workflow (repeatable)",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return hooks.SaveWorkflow(hooks.SaveWorkflowOpts{
						Name:        cmd.String("name"),
						Command:     cmd.String("command"),
						Description: cmd.String("description"),
						RepoScoped:  cmd.Bool("repo"),
						ArgSpecs:    cmd.StringSlice("arg"),
						Tags:        cmd.StringSlice("tag"),
					})
				},
			},
			{
				Name:  "save-launch",
				Usage: "Save a launch configuration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Usage:    "launch configuration name",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:     "tab",
						Usage:    "tab spec: title:cwd:command (repeatable)",
						Required: true,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return hooks.SaveLaunch(hooks.SaveLaunchOpts{
						Name:     cmd.String("name"),
						TabSpecs: cmd.StringSlice("tab"),
					})
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
