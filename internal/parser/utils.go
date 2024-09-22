package parser

import "strings"

func removePrefixes(value string, prefixesToRemove []string) string {

	for _, prefix := range prefixesToRemove {
		//value = strings.TrimPrefix(value, prefix)
		value = strings.Replace(value, prefix, "", -1)
	}
	return value
}
