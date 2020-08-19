package decoder

import (
	"io/ioutil"
	"mime/quotedprintable"
	"strings"
)

func decode_quotedprintable(lines []string) ([]string, error) {
	decoded := []string{}
	for _, line := range lines {
		decoded_line, err := ioutil.ReadAll(quotedprintable.NewReader(strings.NewReader(line)))
		if err != nil {
			return nil, err
		}
		decoded = append(decoded, string(decoded_line))
	}
	return decoded, nil
}

