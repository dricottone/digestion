package decoder

import (
	"strings"
)

const (
	UTF8            = "UTF8"
	Base64          = "Base64"
	QuotedPrintable = "QuotedPrintable"
	Unknown         = "Unknown"
)

func determine_encoding(encoding string) string {
	if strings.Contains(encoding, "base64") {
		return Base64
	} else if strings.Contains(encoding, "quoted-printable") {
		return QuotedPrintable
	} else if strings.Contains(encoding, "utf-8") {
		return UTF8
	} else {
		return Unknown
	}
}

// Translate an encoding label into a numeric value according to preference
// of use in processing. Preference tiers are:
//  1. UTF-8
//  2. base64, quoted-printable
func EvaluateEncoding(encoding string) int {
	switch determine_encoding(encoding) {
	case UTF8:
		return 0
	case Base64:
		return 1
	case QuotedPrintable:
		return 1
	default:
		return 10
	}
}

func DecodeArray(lines []string, encoding string) ([]string, error) {
	switch determine_encoding(encoding) {
	case Base64:
		return decode_base64(lines)
	case QuotedPrintable:
		return decode_quotedprintable(lines)
	default:
		return lines, nil
	}
}

