package ts2go

import (
	"regexp"
	"strings"
)

// FindWordsWithPattern return sub text that match give pattern
func FindWordsWithPattern(str string, pattern string) []string {
	r := regexp.MustCompile(pattern)
	matches := r.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 {
		return []string{}
	}
	return matches[0]
}

// ExactlyMathPattern like its name
func ExactlyMathPattern(str string, pattern string) bool {
	r := regexp.MustCompile(pattern)
	return r.MatchString(str)
}

// SplitAndTrim like its name
func SplitAndTrim(str, delim string, isTrim bool) []string {
	ss := strings.Split(str, delim)
	if isTrim {
		for i, s := range ss {
			ss[i] = strings.TrimSpace(s)
		}
	}
	return ss
}

// UpperFirstChar like it name
func UpperFirstChar(s string) string {
	if s == "" {
		return s
	}
	fc := s[:1]
	fc = strings.ToUpper(fc)
	return fc + s[1:]
}

// LowerFirstChar like it name
func LowerFirstChar(s string) string {
	if s == "" {
		return s
	}
	fc := s[:1]
	fc = strings.ToLower(fc)
	return fc + s[1:]
}
