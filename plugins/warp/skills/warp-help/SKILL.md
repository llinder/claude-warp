---
name: warp-help
description: This skill should be used when the user asks about "warp features", "warp help", "warp notifications", "warp workflows", "warp launch configs", "save a workflow", "create a launch config", or wants to know what the Warp plugin can do.
---

# Warp Terminal Integration

This plugin provides Warp terminal integration for Claude Code.

## Features

### Notifications
Task completion and input-needed notifications are sent automatically via Warp native notifications (OSC 777). Notifications fire when:
- Claude finishes a task and returns to the prompt
- Claude needs user input (elicitation dialogs)
- Claude needs tool approval (permission prompts)

The notification title includes the project name for easy identification across multiple sessions.

### Tab Titles
The Warp tab title is automatically set to "Claude: <project>" on session start. If the title reverts, add this to shell config (`~/.zshrc` or `~/.bashrc`):
```bash
export WARP_DISABLE_AUTO_TITLE=true
```

### Warp Workflows
Save useful commands as Warp workflows accessible via `Ctrl+Shift+R`.

Save a personal workflow:
```bash
claude-warp save-workflow --name <name> --command <cmd> [--description <desc>] [--arg <name:description:default>]...
```

Save a repo workflow (shared with team via git):
```bash
claude-warp save-workflow --repo --name <name> --command <cmd> [--description <desc>] [--arg <name:description:default>]...
```

Arguments use `{{arg_name}}` syntax in the command string.

When discovering useful or repeated commands during a session, proactively suggest saving them as Warp workflows.

### Warp Launch Configs
Create multi-tab/pane dev environment layouts:
```bash
claude-warp save-launch --name <name> --tab <title:cwd:command> [--tab <title:cwd:command>]...
```
