package iplists

import "strings"

type ListType string

const (
	White ListType = "white"
	Black ListType = "black"
)

var allowedListTypes = map[string]ListType{
	"white": White,
	"black": Black,
}

func ParseType(str string) (ListType, bool) {
	t, ok := allowedListTypes[strings.ToLower(str)]
	return t, ok
}
