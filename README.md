# 🧭 Quest

An interactive CLI tool for learning Go through hands-on coding quests with AI-powered hints.

## Quick Start

```bash
# Start a new quest
./quest begin

# Get your next task
./quest next

# Check if your code passes
./quest check

# Need help? Get AI hints
./quest hint

# Mark current task as complete
./quest complete

# View your progress
./quest summary
```

## Commands

### `quest begin`
Start a new coding quest. Choose from:
- **Pick a Legendary Path** - Select from curated templates
- **Forge Your Own Quest** - Customize with AI generation
- **Seek a Mystery Quest** - Get a random AI-generated adventure

### `quest next`
Show your current task with objectives and steps.

### `quest check`
Validate your code against the task requirements. Shows what passed and what failed.

### `quest hint`
Get AI-powered hints for your current task. Uses GitHub Copilot CLI to provide contextual guidance.

### `quest annotate`
Get AI-generated code annotations explaining what needs to be done.

### `quest complete`
Mark the current task as complete and move to the next one.

### `quest jump [task-number]`
Jump to a specific task number.

Options:
- `-l, --last-complete` - Jump to the last completed task

### `quest summary`
View your quest progress across all chapters and tasks.

### `quest help`
Show help for any command.

## How It Works

1. **Begin** - Initialize a quest in the `.quest/` folder
2. **Next** - See what you need to build
3. **Code** - Write your implementation
4. **Check** - Validate your code automatically
5. **Hint** - Get AI help if stuck
6. **Complete** - Move to the next challenge

## Templates

- **go-web-api** - Build a REST API (Quick)
- **go-cli-tool** - Create a CLI tool (Quick)
- **go-concurrency** - Master goroutines (Quick)
- **go-todo-api** - Full REST API with database (Standard)
- **go-auth-system** - Complete auth system (Extended)
- **go-fairy-garden** - Fairy worker service (Standard)
- **go-isekai-server** - Distributed world manager (Master)

## Requirements

- Go 1.18 or higher
- GitHub Copilot CLI (for AI hints)

## Example Session


---

Built with ❤️ for learning Go through interactive quests.
