package copilot_helper

// getIndentation extracts leading whitespace from a line
func GetIndentation(line string) string {
	for i, ch := range line {
		if ch != ' ' && ch != '\t' {
			return line[:i]
		}
	}
	return line // Return entire line if it's all whitespace
}
