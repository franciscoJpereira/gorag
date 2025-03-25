package pkg

import "encoding/base64"

func EncodeBase64(name string) string {
	return base64.StdEncoding.EncodeToString([]byte(name))
}

func DecodeBase64(name string) string {
	result, err := base64.StdEncoding.DecodeString(name)
	if err != nil {
		panic(err.Error())
	}
	return string(result)
}

func DecodeBase64Batch(data []string) []string {
	for index, value := range data {
		data[index] = DecodeBase64(value)
	}
	return data
}
