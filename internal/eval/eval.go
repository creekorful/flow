package eval

import "strings"

// Evaluate given boolean expression using given arguments
func Evaluate(expr string, args map[string]string) bool {
	finalExpr := expr

	for key, value := range args {
		finalExpr = strings.Replace(finalExpr, key, value, -1)
	}

	// split down expression
	parts := strings.Split(finalExpr, " ")
	if len(parts) != 3 {
		return false
	}

	left := parts[0]
	operator := parts[1]
	right := parts[2]

	switch operator {
	case "==":
		return left == right
	case "!=":
		return left != right
	default:
		return false
	}
}
