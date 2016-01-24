package dashboard

import "strings"

func cleanUp(str string) string {
	s := strings.Fields(strings.TrimSpace(str))
	return strings.Join(s, " ")
}
