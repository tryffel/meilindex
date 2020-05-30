package indexer

import (
	"regexp"
	"strings"
)

var attachmentRegex = regexp.MustCompile(`^(.+); name=\"?(.+)\"?$`)

func ParseAttachments(contentType string) string {
	names := attachmentRegex.FindStringSubmatch(contentType)
	if len(names) == 3 {
		name := names[2]
		if strings.HasSuffix(name, `"`) {
			return strings.TrimSuffix(name, `"`)
		} else {
			return name
		}
	}
	return contentType
}
