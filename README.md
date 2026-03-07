# claude-warp

A Claude Code plugin that integrates with [Warp terminal](https://warp.dev), providing native notifications, workflow generation, and launch configurations.

## Features

- **Native notifications** — Task completion and input-needed alerts are forwarded to Warp's notification system via OSC 777 escape sequences
- **Workflow generation** — Save useful commands as [Warp workflows](https://docs.warp.dev/features/workflows) accessible via `Ctrl+Shift+R`
- **Launch configurations** — Create multi-tab dev environment layouts that can be opened in Warp
- **Session context** — Automatically detects Warp terminal and injects integration instructions into Claude's system prompt, including discovery of existing workflows

## Installation

### From the Claude Code marketplace

Add the marketplace and install the plugin:

```sh
claude plugin marketplace add github:llinder/claude-warp
claude plugin install warp@claude-warp
```

To install for a specific scope:

```sh
claude plugin install warp@claude-warp --scope user      # all projects (default)
claude plugin install warp@claude-warp --scope project   # current project only
```

### From source

Requires Go 1.25+.

```sh
git clone https://github.com/llinder/claude-warp.git
cd claude-warp
make build
```

This compiles the binary into `plugins/warp/scripts/claude-warp`.

### Manual install

Copy or symlink the `plugins/warp` directory into your Claude Code plugins location:

```sh
# User-scoped (all projects)
cp -r plugins/warp ~/.claude/plugins/warp

# Project-scoped (per-repo)
cp -r plugins/warp .claude/plugins/warp
```

Claude Code will automatically detect the plugin on the next session start.

## How it works

The plugin registers three Claude Code hooks:

| Hook | Trigger | Behavior |
|------|---------|----------|
| `SessionStart` | Claude Code session begins | Detects Warp, discovers existing workflows, injects system prompt with available commands |
| `Stop` | Session ends | Parses the transcript and sends a summary notification (prompt, files changed, commands run) |
| `Notification` | Claude needs user input | Forwards the notification to Warp's native notification system |

## Usage

Once installed, the plugin works automatically. During a session, Claude can also use the CLI to save workflows and launch configs:

### Save a workflow

```sh
claude-warp save-workflow \
  --name "Deploy staging" \
  --command "kubectl rollout restart deployment/{{service}} -n staging" \
  --description "Restart a staging deployment" \
  --arg "service:Service name:web"
```

- `--repo` flag saves to `.warp/workflows/` in the repo root (shared with the team) instead of your personal `~/.warp/workflows/`
- `--tag` adds searchable tags (repeatable)
- `--arg` format is `name:description:default` (repeatable, default is optional)
- Arguments use `{{arg_name}}` syntax in the command string

### Save a launch configuration

```sh
claude-warp save-launch \
  --name "Dev environment" \
  --tab "Editor:~/project:nvim" \
  --tab "Server:~/project:make run" \
  --tab "Logs:~/project:tail -f logs/dev.log"
```

- `--tab` format is `title:cwd:command` (repeatable)
- Saved to `~/.warp/launch_configurations/`

## Development

```sh
make build    # Build the binary
make test     # Run tests
make lint     # Run go vet
make dist     # Cross-compile for darwin/linux (amd64/arm64)
```

## Plugin structure

```
plugins/warp/
├── .claude-plugin/
│   └── plugin.json      # Plugin metadata
├── hooks/
│   └── hooks.json       # Hook definitions
└── scripts/
    └── claude-warp      # Compiled binary (gitignored)
```

## License

MIT
