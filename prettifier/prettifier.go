package prettifier

import (
	"strings"
)

func Prettify(s string) string {
	var b strings.Builder
	fields := strings.Fields(s)
	var parenthesisStack uint8
	var afterOrderBy bool

	for _, field := range fields {
		switch {

		case strings.HasPrefix(field, "(") && parenthesisStack == 0 && !afterOrderBy:
			b.WriteString("\n\t")
			b.WriteString("(")
			b.WriteString("\n\t\t")
			b.WriteString(field[1:])
			b.WriteString(" ")
			parenthesisStack = parenthesisStack + 1

		case strings.HasSuffix(field, ")") && parenthesisStack == 1:
			b.WriteString(field[:len(field)-1])
			b.WriteString("\n\t) ")
			parenthesisStack = parenthesisStack - 1

		case field == "ORDER" || field == "ENGINE" || field == "PARTITION" || field == "SETTINGS" || field == "FROM" || field == "SELECT" || field == "GROUP":
			b.WriteString("\n")
			b.WriteString(field)
			b.WriteString(" ")
			if field == "ORDER" {
				afterOrderBy = true
			}

		case strings.HasSuffix(field, ",") && !afterOrderBy:
			b.WriteString(field)
			b.WriteString("\n\t\t")

		default:
			b.WriteString(field)
			b.WriteString(" ")
		}
	}

	return b.String()
}
