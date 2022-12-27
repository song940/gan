package gan

import (
	"fmt"
	"regexp"
	"strings"
)

type PathToRegexp struct {
	pattern string
	keys    []string
	regexp  *regexp.Regexp
}

func Compile(path string) (re *PathToRegexp) {
	re = &PathToRegexp{}
	arr := strings.Split(path, "/")
	if arr[0] == "" {
		arr = arr[1:]
	}
	var pattern = ""
	var keys = []string{}
	for _, p := range arr {
		l := len(p)
		var c byte
		if l > 0 {
			c = p[0]
		}
		if c == '*' {
			pattern = pattern + "/(.+)"
			keys = append(keys, p[1:])
		} else if c == ':' {
			r := "([^/]+?)"
			o := p[l-1:]

			if o == "?" {
				pattern = pattern + fmt.Sprintf("(?:/%s)?", r)
			} else if o == "*" {
				pattern = pattern + "(?:/)(.*)"
			} else {
				pattern = pattern + "/" + r
			}
			keys = append(keys, p[1:l])
		} else {
			pattern += "/" + p
		}
	}
	if len(keys) > 0 {
		pattern = pattern + "(?:/)?"
	}
	re.keys = keys
	re.pattern = pattern
	re.regexp = regexp.MustCompile(fmt.Sprintf("^%s/?$", pattern))
	// log.Println("=====>", re)
	return
}

func (re *PathToRegexp) Test(path string) bool {
	return re.regexp.Match([]byte(path))
}

func (re *PathToRegexp) Parse(path string) map[string]string {
	var output = map[string]string{}
	result := re.regexp.FindStringSubmatch(path)
	for i, key := range re.keys {
		output[key] = result[1:][i]
	}
	return output
}
