package config

import (
	"strings"
)

const (
	stringSliceSeparator = ","
)

type stringSlice []string

func newStringSliceValue(val []string, p *[]string) *stringSlice {
	*p = val
	return (*stringSlice)(p)
}

//goland:noinspection GoMixedReceiverTypes
func (s stringSlice) String() string {
	return strings.Join(s, stringSliceSeparator)
}

//goland:noinspection GoMixedReceiverTypes
func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}
