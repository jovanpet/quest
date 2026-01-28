package format

import (
	"fmt"
	"strings"
	"time"
)

// ANSI color codes - Extended palette
const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
	
	// Primary colors
	ColorGreen  = "\033[38;5;42m"   // Bright green
	ColorCyan   = "\033[38;5;51m"   // Bright cyan
	ColorYellow = "\033[38;5;226m"  // Bright yellow
	ColorRed    = "\033[38;5;196m"  // Bright red
	ColorPink   = "\033[38;5;213m"  // Bright pink
	ColorMagenta = "\033[38;5;198m" // Hot pink/magenta
	
	// Accent colors
	ColorPurple = "\033[38;5;141m"  // Soft purple
	ColorBlue   = "\033[38;5;75m"   // Sky blue
	ColorOrange = "\033[38;5;214m"  // Orange
	
	// Backgrounds
	BgGreen  = "\033[48;5;22m"      // Dark green background
	BgCyan   = "\033[48;5;24m"      // Dark cyan background
	BgYellow = "\033[48;5;58m"      // Dark yellow background
	BgPink   = "\033[48;5;126m"     // Pink background
)

// Box drawing characters
const (
	BoxTopLeft     = "╭"
	BoxTopRight    = "╮"
	BoxBottomLeft  = "╰"
	BoxBottomRight = "╯"
	BoxHorizontal  = "─"
	BoxVertical    = "│"
	BoxDot         = "•"
	BoxArrow       = "→"
	BoxCheck       = "✓"
	BoxCross       = "✗"
	BoxStar        = "✨"
)

// Header prints a styled header with icon
func Header(title string) {
	fmt.Println()
	fmt.Printf("  %s%s%s %s%s\n", ColorPink, ColorBold, "🧭", title, ColorReset)
	fmt.Printf("  %s%s%s\n\n", ColorDim, strings.Repeat("─", len(title)+2), ColorReset)
}

// Success prints a success message
func Success(msg string) {
	fmt.Printf("  %s%s %s%s%s\n", BgGreen, ColorBold, BoxCheck, ColorReset, fmt.Sprintf(" %s%s%s", ColorGreen, msg, ColorReset))
	fmt.Println()
}

// Warning prints a warning message
func Warning(msg string) {
	fmt.Printf("  %s⚠  %s%s\n", ColorYellow, msg, ColorReset)
	fmt.Println()
}

// Error prints an error message
func Error(msg string, err error) {
	if err != nil {
		fmt.Printf("  %s%s %s: %v%s\n", ColorRed, BoxCross, msg, err, ColorReset)
	} else {
		fmt.Printf("  %s%s %s%s\n", ColorRed, BoxCross, msg, ColorReset)
	}
	fmt.Println()
}

// ErrorWithTip prints an error with a helpful tip
func ErrorWithTip(msg string, err error, tip string) {
	Error(msg, err)
	if tip != "" {
		fmt.Printf("  %s💡 %s%s\n\n", ColorPink, tip, ColorReset)
	}
}

// Info prints an info message
func Info(msg string) {
	fmt.Printf("  %s%s %s%s\n\n", ColorPink, BoxArrow, msg, ColorReset)
}

// Line prints a plain line
func Line(msg string) {
	fmt.Printf("  %s\n", msg)
}

// Dim prints dimmed text
func Dim(msg string) {
	fmt.Printf("  %s%s%s\n", ColorDim, msg, ColorReset)
}

// List prints a bulleted list
func List(items []string) {
	for _, item := range items {
		fmt.Printf("  %s%s%s %s\n", ColorCyan, BoxDot, ColorReset, item)
	}
	fmt.Println()
}

// Bold prints bold text
func Bold(msg string) {
	fmt.Printf("  %s%s%s\n", ColorBold, msg, ColorReset)
}

// Print prints a plain message
func Print(msg string) {
	fmt.Println(msg)
}

// Printf prints a formatted message
func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Box is deprecated - use Header and Separator instead
// Keeping for backwards compatibility but simplified
func Box(title string, lines []string) {
	Header(title)
	for _, line := range lines {
		fmt.Println("  " + line)
	}
	fmt.Println()
}

// Divider prints a horizontal line separator
func Divider() {
	fmt.Println(ColorDim + "  " + strings.Repeat("─", 50) + ColorReset)
}

// SectionHeader prints a bold section title with divider
func SectionHeader(title string) {
	fmt.Println()
	fmt.Println("  " + ColorPink + ColorBold + title + ColorReset)
	Divider()
}

// Separator prints a decorative separator
func Separator() {
	fmt.Printf("  %s%s%s\n\n", ColorDim, strings.Repeat("─", 60), ColorReset)
}

// Prompt prints a prompt with arrow
func Prompt(msg string) {
	fmt.Printf("  %s→%s %s: ", ColorPink, ColorReset, msg)
}

// TaskHeader prints a task header with divider
func TaskHeader(taskNum int, title string) {
	fmt.Println()
	fmt.Printf("  %s📋 Task %d: %s%s%s\n", ColorPink, taskNum, ColorBold, title, ColorReset)
	fmt.Printf("  %s%s%s\n\n", ColorDim, strings.Repeat("─", 50), ColorReset)
}

// CheckHeader prints check header
func CheckHeader(taskNum int, title string) {
	fmt.Println()
	fmt.Printf("  %s🔍 Checking Task %d: %s%s\n\n",
		ColorPink, taskNum, title, ColorReset)
}

// CheckPass prints a passing check
func CheckPass(name string, reason string) {
	fmt.Printf("  %s✓%s %s%s%s%s\n",
		ColorGreen, ColorReset,
		ColorBold, name, ColorReset,
		fmt.Sprintf(" %s%s%s", ColorDim, reason, ColorReset))
}

// CheckFail prints a failing check
func CheckFail(name string, reason string) {
	fmt.Printf("  %s✗%s %s%s%s%s\n",
		ColorRed, ColorReset,
		ColorBold, name, ColorReset,
		fmt.Sprintf(" %s%s%s", ColorDim, reason, ColorReset))
}

// CheckSummaryPass prints passing summary
func CheckSummaryPass(count int) {
	fmt.Println()
	fmt.Printf("  %s🎉 All %d checks passed!%s\n\n",
		ColorGreen, count, ColorReset)
}

// CheckSummaryFail prints failing summary
func CheckSummaryFail(passed, failed int) {
	fmt.Println()
	fmt.Printf("  %s⚠️  Results: %d passed, %d failed%s\n\n",
		ColorYellow, passed, failed, ColorReset)
}

// CommandHint prints a command suggestion
func CommandHint(label, command string) {
	fmt.Printf("  %s→ %s:%s %s%s%s\n\n",
		ColorPink, label, ColorReset,
		ColorCyan, command, ColorReset)
}

// Step prints a numbered step
func Step(num int, text string) {
	fmt.Printf("    %s%d.%s %s\n", ColorYellow, num, ColorReset, text)
}

// File prints a file item
func File(path string) {
	fmt.Printf("    %s📄 %s%s\n", ColorGreen, path, ColorReset)
}

// Newline prints a blank line
func Newline() {
	fmt.Println()
}

// ExplainHeader prints explain command header with attempt number
func ExplainHeader(attemptNum int, taskTitle string) {
	fmt.Println()
	if attemptNum == 1 {
		fmt.Printf("  %s💡 Getting hints for: %s%s\n\n", 
			ColorPink, taskTitle, ColorReset)
	} else if attemptNum == 2 {
		fmt.Printf("  %s💡 Getting more specific hints (attempt #%d)%s\n", 
			ColorPink, attemptNum, ColorReset)
		fmt.Printf("  %s   Task: %s%s\n\n", 
			ColorDim, taskTitle, ColorReset)
	} else {
		fmt.Printf("  %s💡 Getting direct guidance (attempt #%d)%s\n", 
			ColorPink, attemptNum, ColorReset)
		fmt.Printf("  %s   Task: %s%s\n", 
			ColorDim, taskTitle, ColorReset)
		fmt.Printf("  %s   Don't worry, we'll be more direct this time!%s\n\n", 
			ColorDim, ColorReset)
	}
}

// AnnotationSummary prints annotation results
func AnnotationSummary(count int) {
	if count > 0 {
		fmt.Printf("  %s✓ Added %d inline comment(s) to your code%s\n\n",
			ColorGreen, count, ColorReset)
	} else {
		fmt.Printf("  %s✓ Your code looks good, no additional comments needed%s\n\n",
			ColorGreen, ColorReset)
	}
}

// HintSummary prints hint results
func HintSummary(count int) {
	fmt.Printf("  %s✓ Added %d hint(s) to your code%s\n\n", 
		ColorGreen, count, ColorReset)
}

// FileLocation prints file and line number
func FileLocation(file string, line int) {
	fmt.Printf("    %s%s:%d%s\n", ColorDim, file, line, ColorReset)
}

// Spinner represents a loading spinner
type Spinner struct {
	frames  []string
	current int
	msg     string
	stop    chan bool
}

// NewSpinner creates a new spinner with a message
func NewSpinner(msg string) *Spinner {
	return &Spinner{
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		current: 0,
		msg:     msg,
		stop:    make(chan bool),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	go func() {
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-s.stop:
				return
			case <-ticker.C:
				fmt.Printf("\r  %s%s %s%s", ColorCyan, s.frames[s.current], s.msg, ColorReset)
				s.current = (s.current + 1) % len(s.frames)
			}
		}
	}()
}

// Stop stops the spinner and clears the line
func (s *Spinner) Stop() {
	s.stop <- true
	fmt.Print("\r\033[K")
}
