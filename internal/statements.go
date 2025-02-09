package internal

import (
	"regexp"
	"strings"

	"github.com/gofrs/uuid"
)

var (
	stripCommentsRe   = regexp.MustCompile(`\s+--.*`)
	splitStatementsRe = regexp.MustCompilePOSIX(`;$`)
)

// builtins implements some client-side replacements:
//
// - uuid() replaced with a literal uuid
func builtins(s string) string {
	r := regexp.MustCompile(`(?i)uuid\(\)`)
	s = r.ReplaceAllStringFunc(s, func(_ string) string {
		val := uuid.Must(uuid.NewV4())
		return `'` + val.String() + `'`
	})
	return s
}

func Statements(contents []byte) []string {
	result := []string{}

	// remove sql comments from anywhere ([whitespace]--*\n)
	contents = stripCommentsRe.ReplaceAll(contents, nil)

	// split statements by trailing ; at the end of the line
	stmts := splitStatementsRe.Split(string(contents), -1)
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, builtins(stmt))
		}
	}

	return result
}
