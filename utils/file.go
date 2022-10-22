package utils

import "strings"

func FID(name string) string{
	return strings.Split(name, ".")[0]
}