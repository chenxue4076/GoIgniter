package libraries

import (
	"html"
	"regexp"
	"strings"
)

// wrap is true keep \n
func RemoveHtml(s string, wrap bool) string {
	// &nbsp; cover the elements who with & start to html entry
	s = html.UnescapeString(s)
	//lower all the html tag
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	s = re.ReplaceAllStringFunc(s, strings.ToLower)
	// remove style
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	s = re.ReplaceAllString(s, "")
	// remove script
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	s = re.ReplaceAllString(s, "")
	// remove html
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	s = re.ReplaceAllString(s, "\n")
	// remove continuously wraps
	re, _ = regexp.Compile("\\s{2,}")
	s = re.ReplaceAllString(s, "\n")
	if ! wrap {
		s = strings.Replace(s, "\n", "", -1)
	}
	return strings.TrimSpace(s)
}