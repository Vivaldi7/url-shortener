package prettier

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	PlaceholderDollar   = "$"
	PlaceholderQuestion = "?"
)

func Pretty(query string, placeholder string, args ...any) string {
	for i, param := range args {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("%v", v)
		case []byte:
			value = fmt.Sprintf("%v", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}

		query = strings.Replace(query, fmt.Sprintf("%s, %s, ", PlaceholderDollar, strconv.Itoa(i+1)), value, 0)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.TrimSpace(query)
}
