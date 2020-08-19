package renderer

import (
	"strings"
)

const (
	FormattedTextPlain = "FormattedTextPlain"
	FormattedTextHTML  = "FormattedTextHTML"
	FormattedUnknown   = "FormattedUnknown"
)

func DetermineFormatting(formatting string) string {
	if strings.Contains(formatting, "text/plain") {
		return FormattedTextPlain
	} else if strings.Contains(formatting, "text/html") {
		return FormattedTextHTML
	} else {
		return FormattedUnknown
	}
}

