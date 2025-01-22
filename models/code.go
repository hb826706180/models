package models

import "encoding/base64"

type Code_ struct {
}

func Base64Encode(encrypt interface{}) string {
	s := base64.StdEncoding.EncodeToString(encrypt.([]byte))
	return s
}
func Base64Decode(encrypt string) ([]byte, error) {
	by, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return nil, err
	}
	return by, nil
}
