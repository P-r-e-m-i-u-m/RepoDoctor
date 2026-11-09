package doctor

import (
	"regexp"
	"strings"
)

var blockedPatterns = []*regexp.Regexp{
	regexp.MustCompile("(?i)\\brm\\s+-rf\\b"),
	regexp.MustCompile("(?i)\\bgit\\s+reset\\s+--hard\\b"),
	regexp.MustCompile("(?i)\\bgit\\s+push\\b[^\\n]*(?:\\s--force(?:-with-lease)?\\b|\\s-[a-z]*f[a-z]*(?:\\s|$))"),
	regexp.MustCompile("(?i)\\b(drop|truncate)\\s+(database|table)\\b"),
	regexp.MustCompile("(?i)\\bdel\\s+/[fq]\\b"),
	regexp.MustCompile("(?i)\\brmdir\\s+/s\\b"),
	regexp.MustCompile("(?i)\\bgit\\s+clean\\b[^\\n]*(?:-[a-z]*f[a-z]*(?:\\s|$)|--force\\b)"),
	regexp.MustCompile("(?i)\\bdocker(?:-compose|\\s+compose)\\s+down\\b[^\\n]*(?:--volumes\\b|-v\\b)"),
	regexp.MustCompile("(?i)\\bremove-item\\b[^\\n]*-recurse\\b[^\\n]*-force\\b"),
	regexp.MustCompile("(?i)\\bremove-item\\b[^\\n]*-force\\b[^\\n]*-recurse\\b"),
}

func IsDestructive(command string) bool {
	command = strings.TrimSpace(command)
	for _, pattern := range blockedPatterns {
		if pattern.MatchString(command) {
			return true
		}
	}
	return false
}
