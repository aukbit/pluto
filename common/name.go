package common

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func SafeName(name, sufix string) string {
	// support only alphanumeric and underscore characters
	reg, err := regexp.Compile("[^A-Za-z0-9_]+")
	if err != nil {
		log.Fatal(err)
	}
	safe := reg.ReplaceAllString(name, "_")
	return fmt.Sprintf("%s_%s", strings.ToLower(safe), sufix)
}
