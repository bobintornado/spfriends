package relationship

import "regexp"

func ExtractEmails(text string) (list []string) {
	re := regexp.MustCompile("[a-z0-9-]{1,30}@[a-z0-9-]{1,65}.[a-z]{1,}")
	list = re.FindAllString(text, -1)

	return list
}
