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

# Need help? Get AI explanations
./quest explain

# Mark current task as complete
./quest complete

# View your progress
./quest summary

#Check the health of the quest
./quest health
```

## Commands

### `quest begin`
Start a new coding quest. Choose from curated templates or generate custom quests with AI.

### `quest next`
Move to the next task in your quest and display what you need to work on.

### `quest check`
Validate that you have completed the requirements for the current task. Shows what passed and what failed.

Options:
- `-a, --annotate` - Add inline comments to code showing check results

### `quest explain`
Get AI-powered explanations and hints for the current task. The AI analyzes your code and provides contextual guidance, with increasing detail based on how many times you've requested help.

### `quest complete`
Mark the current task as complete and move to the next one.

### `quest jumpTo [task-number]`
Jump to a specific task index or the last completed task in your quest.

Options:
- `-l, --last-complete` - Jump to the last completed task

### `quest summary`
View your quest progress across all chapters and tasks.

### `quest health`
Checks the health of the quest.

## How It Works

1. **Begin** - Initialize a quest in the `.quest/` folder
2. **Next** - See what you need to build
3. **Code** - Write your implementation
4. **Check** - Validate your code automatically
5. **Explain** - Get AI help if stuck
6. **Complete** - Move to the next challenge

## Available Templates

### Quick (3 tasks each)
- **go-web-api** - REST API with Go - Learn to build HTTP servers, handle JSON, and write tests
- **go-cli-tool** - CLI Tool with Cobra - Create command-line tools with subcommands and flags
- **go-concurrency** - Go Concurrency Patterns - Master goroutines, channels, and worker pools

### Normal (10 tasks)
- **go-fairy-garden** - Fairy Worker Service - Build a whimsical worker service with Go

### Advanced
- **go-todo-api** - Complete Todo REST API (15 tasks) - Build a full-featured REST API with CRUD, testing, and middleware
- **go-auth-system** - User Authentication System (19 tasks) - Build complete auth with database, sessions, JWT, and security
- **go-isekai-server** - Distributed World Manager (20 tasks) - Manage a distributed virtual world with Go

## Requirements

- Go 1.18 or higher
- GitHub Copilot CLI (for AI hints via `quest explain`)

## Installation

### Option 1: Install via Go (Recommended)

```bash
go install github.com/jovanpet/quest@latest
```

This will install the `quest` binary to your `$GOPATH/bin` directory. Make sure `$GOPATH/bin` is in your PATH.

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/jovanpet/quest.git
cd quest

# Build the binary
go build -o quest

# Run it
./quest begin
```

### Option 3: Download Pre-built Binaries

Download the latest release from the [releases page](https://github.com/jovanpet/quest/releases).

---

Built with ❤️ for learning Go through interactive quests.
