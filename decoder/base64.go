package decoder

import (
	"encoding/base64"
)

func decode_base64(lines []string) ([]string, error) {
	decoded := []string{}
	for _, line := range lines {
		decoded_line, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return decoded, err
		}
		decoded = append(decoded, string(decoded_line))
	}
	return decoded, nil
}

