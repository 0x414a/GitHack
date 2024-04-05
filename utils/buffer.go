package utils

import "bytes"

func ReadStringWithoutDelimiter(buffer *bytes.Buffer, delim byte) (string, error) {
	str, err := buffer.ReadString(delim)
	if err != nil {
		return "", err
	}

	if len(str) > 0 {
		str = str[:len(str)-1]
	}
	return str, nil
}
