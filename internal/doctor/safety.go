package doctor

import (
 "regexp"
 "strings"
)

var blockedPatterns=[]*regexp.Regexp{
 regexp.MustCompile("(?i)\\brm\\s+-rf\\b"),
 regexp.MustCompile("(?i)\\bgit\\s+reset\\s+--hard\\b"),
 regexp.MustCompile("(?i)\\bgit\\s+push\\b.*\\s--force(?:-with-lease)?\\b"),
 regexp.MustCompile("(?i)\\b(drop|truncate)\\s+(database|table)\\b"),
 regexp.MustCompile("(?i)\\bdel\\s+/[fq]\\b"),
 regexp.MustCompile("(?i)\\brmdir\\s+/s\\b"),
}

func IsDestructive(command string) bool {
 command=strings.TrimSpace(command)
 for _,pattern:=range blockedPatterns { if pattern.MatchString(command) { return true } }
 return false
}
