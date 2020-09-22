package utils

import "strings"

func ChangeSpan(str string) []string {
	str = strings.Replace(str, `<span class="name">`, " ", -1)
	str = strings.Replace(str, `<span class="divide">`, " ", -1)
	str = strings.Replace(str, `</span>`, " ", -1)
	str = strings.Replace(str, `<small>`, " ", -1)
	str = strings.Replace(str, `</small>`, " ", -1)
	return strings.Split(str, " ")
}

func DelSpace(str string) string {
	str = strings.Replace(str, "\n", " ", -1)
	str = strings.Replace(str, "\n\r", " ", -1)
	str = strings.Replace(str, "\r", " ", -1)
	str = strings.Replace(str, "	", " ", -1)
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
