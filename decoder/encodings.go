package decoder

import (
	"strings"
)

const (
	EncodedUTF8            = "EncodedUTF8"
	EncodedBase64          = "EncodedBase64"
	EncodedQuotedPrintable = "EncodedQuotedPrintable"
	EncodedUnknown         = "EncodedUnknown"
)

func DetermineEncoding(encoding string) string {
	if strings.Contains(encoding, "utf-8") {
		return EncodedUTF8
	} else if strings.Contains(encoding, "base64") {
		return EncodedBase64
	} else if strings.Contains(encoding, "quoted-printable") {
		return EncodedQuotedPrintable
	} else {
		return EncodedUnknown
	}
}

func DecodeArray(lines []string, encoding string) ([]string, error) {
	switch DetermineEncoding(encoding) {
	case EncodedBase64:
		return decode_base64(lines)
	case EncodedQuotedPrintable:
		return decode_quotedprintable(lines)
	default:
		return lines, nil
	}
}

